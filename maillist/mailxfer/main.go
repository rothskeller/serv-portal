package main

import (
	"bytes"
	"cmp"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"sunnyvaleserv.org/portal/maillist/private"
)

const transferURL = "https://sunnyvaleserv.org/mailrecv.cgi"

func main() {
	lambda.Start(transferMail)
}

func transferMail(ctx context.Context, input *SESInput) (err error) {
	var (
		conf     aws.Config
		s3Client *s3.Client
	)
	conf, _ = config.LoadDefaultConfig(ctx)
	s3Client = s3.NewFromConfig(conf)
	for _, record := range input.Records {
		if err := transferOneMail(ctx, s3Client, record); err != nil {
			return err
		}
	}
	return nil
}

func transferOneMail(ctx context.Context, s3Client *s3.Client, record SESInputRecord) (err error) {
	var (
		msgid string
		raw   []byte
		buf   bytes.Buffer
		mpw   *multipart.Writer
		req   *http.Request
		resp  *http.Response
	)
	msgid = record.SES.Mail.MessageID
	log.Printf("transferMail(%s):", msgid)
	if raw, err = readMail(ctx, s3Client, msgid); err != nil {
		return fmt.Errorf("reading message: %s", err)
	}
	mpw = multipart.NewWriter(&buf)
	if err := cmp.Or(
		writeFormField(mpw, "sesID", msgid),
		writeFormField(mpw, "sesHash", private.ComputeHash(msgid)),
		writeFormField(mpw, "dkim", record.SES.Receipt.DKIMVerdict.Status),
		writeFormField(mpw, "dmarc", record.SES.Receipt.DMARCVerdict.Status),
		writeFormField(mpw, "spf", record.SES.Receipt.SPFVerdict.Status),
		writeFormField(mpw, "spam", record.SES.Receipt.SpamVerdict.Status),
		writeFormField(mpw, "virus", record.SES.Receipt.VirusVerdict.Status)); err != nil {
		return err
	}
	for _, recip := range record.SES.Receipt.Recipients {
		if err = writeFormField(mpw, "recipient", recip); err != nil {
			return err
		}
	}
	if w, err := mpw.CreateFormField("message"); err != nil {
		return fmt.Errorf("CreateFormField message: %w", err)
	} else if _, err = w.Write(raw); err != nil {
		return fmt.Errorf("CreateFormField message: copy: %w", err)
	}
	if err = mpw.Close(); err != nil {
		return fmt.Errorf("multipart.Writer.Close: %w", err)
	}
	if req, err = http.NewRequestWithContext(ctx, http.MethodPost, transferURL, &buf); err != nil {
		return fmt.Errorf("NewRequest: %w", err)
	}
	req.Header.Set("Content-Type", mpw.FormDataContentType())
	if resp, err = http.DefaultClient.Do(req); err != nil {
		return fmt.Errorf("POST mail: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		buf, _ := io.ReadAll(resp.Body)
		if len(buf) != 0 {
			return fmt.Errorf("POST mail: %s", string(buf))
		} else {
			return fmt.Errorf("POST mail: %s", resp.Status)
		}
	}
	// TODO: delete mail from S3
	return nil
}

func writeFormField(mpw *multipart.Writer, name, value string) (err error) {
	if w, err := mpw.CreateFormField(name); err != nil {
		return fmt.Errorf("CreateFormField %s: %w", name, err)
	} else {
		io.WriteString(w, value)
	}
	return nil
}

func readMail(ctx context.Context, client *s3.Client, msgid string) (raw []byte, err error) {
	var resp *s3.GetObjectOutput
	if resp, err = client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String("serv-mail"),
		Key:    &msgid,
	}); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if raw, err = io.ReadAll(resp.Body); err != nil {
		return nil, err
	}
	return raw, nil
}

type (
	SESInput struct {
		Records []SESInputRecord
	}
	SESInputRecord struct {
		SES SESData
	}
	SESData struct {
		Mail    SESMail
		Receipt SESReceipt
	}
	SESMail struct {
		MessageID string // also S3 object key
	}
	SESReceipt struct {
		VirusVerdict Verdict
		DMARCVerdict Verdict
		Recipients   []string // bare addresses, names removed
		SpamVerdict  Verdict
		SPFVerdict   Verdict
		DKIMVerdict  Verdict
	}
	Verdict struct {
		Status string
	}
)
