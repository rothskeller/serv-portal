package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/mail"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ses"
)

const (
	adminFrom = "SunnyvaleSERV.org <admin@mx.sunnyvaleserv.org>"
	adminTo   = "Steve Roth <sroth@sunnyvale.ca.gov>"
)

func main() {
	lambda.Start(handleMail)
}

func handleMail(ctx context.Context, input *SESInput) (err error) {
	var (
		conf      aws.Config
		s3Client  *s3.Client
		sesClient *ses.Client
		msgid     string
		hdr       mail.Header
		body      []byte
		raw       []byte
		listdata  AllListData
	)
	conf, _ = config.LoadDefaultConfig(ctx)
	s3Client = s3.NewFromConfig(conf)
	sesClient = ses.NewFromConfig(conf)
	msgid = input.Records[0].SES.Mail.MessageID
	if hdr, body, raw, err = readMail(ctx, s3Client, msgid); err != nil {
		return fmt.Errorf("reading message: %s", err)
	}
	if listdata, err = readListData(ctx, s3Client); err != nil {
		return fmt.Errorf("reading list data: %s", err)
	}
	for _, recip := range input.Records[0].SES.Receipt.Recipients {
		recip, _, _ = strings.Cut(recip, "@")
		if strings.HasSuffix(recip, ".mod") {
			err = handleModerationResponse(ctx, s3Client, sesClient, recip[:len(recip)-4], listdata, body)
		} else if ld := listdata[recip]; ld != nil {
			err = maybeSendMessageToList(
				ctx, sesClient, &input.Records[0].SES.Receipt, input.Records[0].SES.Mail.MessageID,
				hdr, body, raw, recip, ld)
		} else {
			err = unknownRecipient(ctx, sesClient, hdr, body, raw, recip)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func unknownRecipient(ctx context.Context, client *ses.Client, hdr mail.Header, body, raw []byte, recip string) (err error) {
	var (
		subject string
		comment string
	)
	if s := hdr.Get("Subject"); s != "" {
		subject = "ADMIN: " + s
	} else {
		subject = "ADMIN: (no subject)"
	}
	comment = fmt.Sprintf("ERROR: No such mailing list %q.", html.EscapeString(recip))
	return forwardMessage(ctx, client, hdr, body, raw, adminFrom, "", []string{adminTo}, subject, comment)
}

func readMail(ctx context.Context, client *s3.Client, msgid string) (hdr mail.Header, body, raw []byte, err error) {
	var (
		resp *s3.GetObjectOutput
		msg  *mail.Message
	)
	if resp, err = client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String("serv-mail"),
		Key:    &msgid,
	}); err != nil {
		return nil, nil, nil, err
	}
	defer resp.Body.Close()
	if raw, err = io.ReadAll(resp.Body); err != nil {
		return nil, nil, nil, err
	}
	if msg, err = mail.ReadMessage(bytes.NewReader(raw)); err != nil {
		return nil, nil, nil, err
	}
	if body, err = io.ReadAll(msg.Body); err != nil {
		return nil, nil, nil, err
	}
	return msg.Header, body, raw, nil
}

func readListData(ctx context.Context, client *s3.Client) (listdata AllListData, err error) {
	var resp *s3.GetObjectOutput

	if resp, err = client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String("serv-mail"),
		Key:    aws.String("list-data.json"),
	}); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&listdata); err != nil {
		return nil, err
	}
	return listdata, nil
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
		CommonHeaders CommonHeaders
		MessageID     string // also S3 object key
	}
	CommonHeaders struct {
		Subject string
		From    []string
		Sender  string
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
	AllListData map[string]*ListData
	ListData    struct {
		Senders    []string
		Moderators []string
		Receivers  []Receiver
	}
	Receiver struct {
		Name  string
		Addr  string
		Token string
	}
)
