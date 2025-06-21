package main

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"io"
	"log"
	"maps"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"sunnyvaleserv.org/portal/maillist"
	"sunnyvaleserv.org/portal/util/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	aconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"k8s.io/apimachinery/pkg/util/sets"
	"zombiezen.com/go/sqlite"
)

var (
	toHandle  sets.Set[string]
	dbconn    *sqlite.Conn
	sesClient *ses.Client
)

func main() {
	var (
		conf aws.Config
		err  error
	)
	// Move to maillist directory.
	if err = os.Chdir("/home/snyserv/sunnyvaleserv.org/data"); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: chdir data: %s\n", err)
		os.Exit(1)
	}
	// Lock the lockfile to ensure only one instance running.
	if lockf, err := os.OpenFile("maillist/LOCK", os.O_CREATE|os.O_WRONLY, 0666); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	} else if err = syscall.Flock(int(lockf.Fd()), syscall.LOCK_EX); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: LOCK: %s\n", err)
		os.Exit(1)
	}
	// Open the logfile and prepare to log to it.
	logname := "maillist/log/" + time.Now().Format("2006-01")
	if logf, err := os.OpenFile(logname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	} else {
		log.SetOutput(logf)
	}
	// Open the database.
	if dbconn, err = sqlite.OpenConn(config.Get("databaseFilename"), sqlite.OpenReadOnly|sqlite.OpenNoMutex); err != nil {
		log.Fatalf("ERROR: open DB: %s", err)
	}
	defer dbconn.Close()
	// Set the journal mode to truncate.
	if stmt, _, err := dbconn.PrepareTransient("PRAGMA journal_mode = TRUNCATE"); err != nil {
		log.Fatalf("ERROR: prepare stmt: %s", err)
	} else if _, err = stmt.Step(); err != nil {
		log.Fatalf("ERROR: set journal mode: %s", err)
	} else if err = stmt.Finalize(); err != nil {
		log.Fatalf("ERROR: finalize stmt: %s", err)
	}
	if conf, err = aconfig.LoadDefaultConfig(context.Background(), aconfig.WithCredentialsProvider(aws.CredentialsProviderFunc(func(_ context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     config.Get("sendmailAccessKey"),
			SecretAccessKey: config.Get("sendmailSecretKey"),
		}, nil
	}))); err != nil {
		log.Fatalf("ERROR: load AWS config: %s", err)
	}
	sesClient = ses.NewFromConfig(conf)
	// Get the list of mails to be handled.
	toHandle = getMailsToHandle()
	// Handle each of them.
	for {
		if messageID, ok := toHandle.PopAny(); ok {
			handleMail(messageID)
		} else {
			break
		}
	}
}

// getMailsToHandle returns the initial set of mails that need to be handled.
// (Additional ones could be added later during processing.)  This is
// determined by looking for QUEUE/*.data files that have their other-read mode
// bit set.
//
// As a side effect, it also removes outdated message files:  those without the
// other-read mode bit that are over one month old.
func getMailsToHandle() (toHandle sets.Set[string]) {
	var (
		ents        []os.DirEntry
		err         error
		oneMonthAgo = time.Now().AddDate(0, -1, 0)
	)
	toHandle = sets.New[string]()
	if ents, err = os.ReadDir("maillist/QUEUE"); err != nil {
		log.Fatalf("ERROR: QUEUE: %s", err)
	}
	for _, ent := range ents {
		if !ent.Type().IsRegular() || !strings.HasSuffix(ent.Name(), ".data") {
			continue
		}
		if stat, err := ent.Info(); err == nil {
			messageID := strings.TrimSuffix(ent.Name(), ".data")
			if stat.Mode()&04 == 04 {
				toHandle.Insert(messageID)
			} else if stat.ModTime().Before(oneMonthAgo) {
				log.Printf("Removing outdated QUEUE/%s", messageID)
				os.Remove("maillist/QUEUE/" + ent.Name())
				os.Remove("maillist/QUEUE/" + messageID)
			}
		}
	}
	return toHandle
}

func handleMail(messageID string) {
	var (
		tf       *os.File
		mail     *mailMetadata
		hasError bool
		err      error
	)
	log.Printf("Handling %s:", messageID)
	// Open the data file, lock it, and read it.
	if tf, err = os.OpenFile(fmt.Sprintf("maillist/QUEUE/%s.data", messageID), os.O_RDWR, 0666); err != nil {
		log.Printf("ERROR: %s", err)
		return
	}
	defer tf.Close()
	if mail, err = readTracking(tf); err != nil {
		log.Printf("ERROR: QUEUE/%s.data: %s", messageID, err)
		return
	}
	for list := range maps.Keys(mail.lists) {
		if err = handleMailToList(tf, messageID, mail, list); err != nil {
			log.Printf("ERROR: %s: %s: %s", messageID, list, err)
			hasError = true
		}
	}
	if !hasError {
		if err = tf.Chmod(0640); err != nil {
			log.Printf("ERROR: QUEUE/%s.data: chmod: %s", messageID, err)
		}
		log.Printf("  Marked as handled.")
	}
}

func handleMailToList(tf *os.File, messageID string, mail *mailMetadata, listaddr string) (err error) {
	listname, _, _ := strings.Cut(strings.ToLower(listaddr), "@")
	if mail.sentToList.Has(listname) {
		log.Printf("  Already handled destination %s.", listname)
		return nil
	}
	if strings.HasSuffix(listname, ".mod") {
		return handleModerationResponse(messageID, strings.TrimSuffix(listname, ".mod"))
	}
	list := maillist.GetList(dbconn, listname)
	if list == nil {
		return handleUnknownList(tf, messageID, listname)
	}
	if !mail.approved.Has(listname) {
		if mail.moderating.Has(listname) {
			log.Printf("  Moderation already requested for %s.", listname)
			return nil
		} else if problems := messageNeedsModeration(messageID, list, mail); len(problems) != 0 {
			return requestModeration(tf, messageID, list, problems)
		}
	}
	return sendListEmail(tf, messageID, list, mail)
}

func handleUnknownList(tf *os.File, messageID, listname string) (err error) {
	var (
		fname   string
		raw     []byte
		msg     *mail.Message
		body    []byte
		subject string
		comment string
	)
	log.Printf("  Unknown recipient %s, forwarding to admin", listname)
	fname = filepath.Join("maillist/QUEUE", messageID)
	if raw, err = os.ReadFile(fname); err != nil {
		return fmt.Errorf("read message: %w", err)
	}
	if msg, err = mail.ReadMessage(bytes.NewReader(raw)); err != nil {
		return fmt.Errorf("parse message: %w", err)
	}
	if body, err = io.ReadAll(msg.Body); err != nil {
		return fmt.Errorf("read message body: %w", err)
	}
	if s := msg.Header.Get("Subject"); s != "" {
		subject = "ADMIN: " + s
	} else {
		subject = "ADMIN: (no subject)"
	}
	comment = fmt.Sprintf("ERROR: No such mailing list %q.", html.EscapeString(listname))
	if err = forwardMessage(msg.Header, body, raw, config.Get("adminFrom"), "", strings.Split(config.Get("adminTo"), ","), subject, comment); err != nil {
		return err
	}
	tstamp := time.Now().Format(time.RFC3339)
	fmt.Fprintf(tf, "L %s %s\n", tstamp, listname)
	return nil
}
