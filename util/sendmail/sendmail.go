package sendmail

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var mailSenderPath string

func init() {
	var err error

	if mailSenderPath, err = exec.LookPath("mail-sender"); err != nil {
		home := os.Getenv("HOME")
		if home == "" {
			home = "/home/snyserv"
		}
		mailSenderPath = filepath.Join(home, "go", "bin", "mail-sender")
	}
}

// A Mailer is a handler for sending email.  While the current email sending
// method is stateless, other possible methods aren't (e.g. directly connecting
// to an SMTP server), so this API allows for state to be preserved across
// multiple messages.
type Mailer struct{}

// OpenMailer creates a handler for sending email.
func OpenMailer() (m *Mailer, err error) {
	return new(Mailer), nil
}

// SendMessage sends a single message through the Mailer.  If it returns an
// error, the Mailer is no longer usable.
func (m *Mailer) SendMessage(from string, to []string, body []byte) (err error) {
	return SendMessage(from, to, body)
}

// Close closes the Mailer.  The Mailer may not be used after this is called.
func (m *Mailer) Close() {}

// SendMessage sends a single email message.
func SendMessage(from string, to []string, body []byte) (err error) {
	var (
		args []string
		cmd  *exec.Cmd
		out  []byte
	)
	args = make([]string, 0, len(to)+1)
	args = append(args, from)
	args = append(args, to...)
	cmd = exec.Command(mailSenderPath, args...)
	cmd.Stdin = bytes.NewReader(body)
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(out))
	}
	return nil
}
