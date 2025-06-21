package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cgi"
	"net/mail"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"sunnyvaleserv.org/portal/maillist/private"
)

func main() {
	if fh, err := os.OpenFile("/tmp/serv-mailrecv.err", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666); err == nil {
		os.Stderr = fh
	}
	if err := os.Chdir("/home/snyserv/sunnyvaleserv.org/data/maillist"); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	cgi.Serve(http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	var (
		sesID   string
		sesHash string
		mf      *os.File
		msg     *mail.Message
		mid     string
		qfn     string
		tfn     string
		tfh     *os.File
		existed bool
		err     error
	)
	// Verify that we have a valid request:  SSL over POST with a correctly
	// hashed SES ID.
	if r.TLS == nil || !r.TLS.HandshakeComplete {
		sendError(w, "403 TLS Required", http.StatusForbidden)
		return
	}
	if r.Method != http.MethodPost {
		sendError(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if sesID = r.FormValue("sesID"); sesID == "" {
		sendError(w, "400 sesID Required", http.StatusBadRequest)
		return
	}
	if sesHash = r.FormValue("sesHash"); sesHash == "" {
		sendError(w, "400 sesHash Required", http.StatusBadRequest)
		return
	}
	if sesHash != private.ComputeHash(sesID) {
		sendError(w, "403 Incorrect Hash", http.StatusForbidden)
		return
	}
	// Open a temporary file and read the message into it.
	if mf, err = os.CreateTemp(".", ""); err != nil {
		sendError(w, "500 CreateTemp "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(mf.Name())
	if _, err = io.WriteString(mf, r.FormValue("message")); err != nil {
		sendError(w, "500 copy message "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Get the message ID from the message.
	if _, err = mf.Seek(0, 0); err != nil {
		sendError(w, "500 rewind message "+err.Error(), http.StatusInternalServerError)
		return
	}
	if msg, err = mail.ReadMessage(mf); err != nil {
		sendError(w, "400 invalid message "+err.Error(), http.StatusBadRequest)
		return

	}
	mf.Close()
	if mid = msg.Header.Get("Message-ID"); mid == "" {
		mid = sesID
	} else if addr, err := mail.ParseAddress(mid); err != nil {
		mid = sesID
	} else {
		mid = addr.Address
	}
	// Hash the message ID and base64-encode the hash.  This ensures that
	// our message key has nothing in it that mailers will consider to be
	// a URL, and it's short enough that mailers won't word-wrap it.
	// This is important when message keys are embedded in moderation
	// emails.
	midhash := md5.Sum([]byte(mid))
	mid = base64.RawURLEncoding.EncodeToString(midhash[:])
	// Add message information to the tracking file, creating it if needed.
	qfn = filepath.Join("QUEUE", mid)
	tfn = qfn + ".data"
	if tfh, err = os.OpenFile(tfn, os.O_CREATE|os.O_EXCL|os.O_APPEND|os.O_WRONLY, 0666); os.IsExist(err) {
		existed = true
		tfh, err = os.OpenFile(tfn, os.O_APPEND|os.O_WRONLY, 0666)
	}
	if err != nil {
		sendError(w, "500 open tracking file "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err = syscall.Flock(int(tfh.Fd()), syscall.LOCK_EX); err != nil {
		sendError(w, "500 lock tracking file "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(tfh, "R %s %s\n", time.Now().Format(time.RFC3339), strings.Join(r.Form["recipient"], " "))
	fmt.Fprintf(tfh, "V %s %s %s %s %s\n", r.FormValue("dkim"), r.FormValue("dmarc"), r.FormValue("spf"), r.FormValue("spam"), r.FormValue("virus"))
	// If the tracking file existed, set its other-read bit.
	// Otherwise, save the message file before unlocking it.
	if existed {
		st, _ := tfh.Stat()
		if err = tfh.Chmod(st.Mode() | 04); err != nil {
			sendError(w, "500 chmod tracking file "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if err = os.Rename(mf.Name(), qfn); err != nil {
			sendError(w, "500 move message "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if err = tfh.Close(); err != nil {
		sendError(w, "500 close tracking file "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Start the routemail program in the background to handle it.
	cmd := exec.Command("/home/snyserv/bin/routemail")
	if err = cmd.Start(); err != nil {
		sendError(w, "500 start routemail "+err.Error(), http.StatusInternalServerError)
		return
	}
	// All is well.
	w.WriteHeader(http.StatusNoContent)
}

func sendError(w http.ResponseWriter, estr string, code int) {
	http.Error(w, estr, code)
	log.Printf("ERROR: %s", estr)
}
