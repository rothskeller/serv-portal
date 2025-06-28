package sendmail

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	aconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"sunnyvaleserv.org/portal/util/config"
)

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
func (m *Mailer) SendMessage(ctx context.Context, from string, to []string, body []byte) (err error) {
	return SendMessage(ctx, from, to, body)
}

// Close closes the Mailer.  The Mailer may not be used after this is called.
func (m *Mailer) Close() {}

// SendMessage sends a single email message.
func SendMessage(ctx context.Context, from string, to []string, body []byte) (err error) {
	var (
		conf   aws.Config
		client *ses.Client
	)
	conf, _ = aconfig.LoadDefaultConfig(ctx, aconfig.WithCredentialsProvider(aws.CredentialsProviderFunc(func(_ context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     config.Get("sendmailAccessKey"),
			SecretAccessKey: config.Get("sendmailSecretKey"),
		}, nil
	})))
	client = ses.NewFromConfig(conf)
	cset := "serv-outgoing"
	_, err = client.SendRawEmail(ctx, &ses.SendRawEmailInput{RawMessage: &types.RawMessage{Data: body}, ConfigurationSetName: &cset})
	if err != nil {
		return fmt.Errorf("AWS SendRawEmail: %s", err)
	}
	return nil
}
