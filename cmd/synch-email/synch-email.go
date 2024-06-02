package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/store/listperson"
	"sunnyvaleserv.org/portal/util/config"
	"sunnyvaleserv.org/portal/util/log"
)

func main() {
	var (
		entry  *log.Entry
		data   []byte
		client *s3.Client
		err    error
	)
	// Find the database.
	switch os.Getenv("HOME") {
	case "/home/snyserv":
		if err := os.Chdir("/home/snyserv/sunnyvaleserv.org/data"); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	case "/Users/stever":
		if err := os.Chdir("/Users/stever/src/serv-portal/data"); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	}
	// Generate the JSON data block.
	entry = log.New("", "synch-email")
	store.Connect(context.Background(), entry, func(st *store.Store) {
		data = listperson.ListData(st)
	})
	// Connect to AWS S3.
	client = s3.New(s3.Options{
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
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}
