package phys

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"sunnyvaleserv.org/portal/util/config"
)

// UploadEmailListData sends the email list data to our AWS instance.
func UploadEmailListData(storer Storer, data []byte) (err error) {
	var olddata string

	// Get the last-uploaded data for comparison.
	SQL(storer, "SELECT data FROM list_data", func(stmt *Stmt) {
		if stmt.Step() {
			olddata = stmt.ColumnText()
		}
	})
	if olddata == string(data) {
		return nil
	}
	// Connect to AWS S3.
	var client = s3.New(s3.Options{
		Region: "us-east-1",
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			config.Get("listDataUploadAccessKey"),
			config.Get("listDataUploadSecretKey"),
			"",
		)),
	})
	// Upload the JSON data block.
	if _, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("serv-mail"),
		Key:    aws.String("list-data.json"),
		Body:   bytes.NewReader(data),
	}); err != nil {
		storer.Problems().AddError(err)
		return err
	}
	// Save the newly uploaded data.
	storer.AsStore().Transaction(func() {
		SQL(storer, "UPDATE list_data SET data=?", func(stmt *Stmt) {
			stmt.BindText(string(data))
			stmt.Step()
		})
	})
	return nil
}
