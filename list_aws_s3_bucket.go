package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

func ListBuckets(ctx context.Context, client *s3.S3) (*s3.ListBucketsOutput, error) {
	// Get all list buckets
	res, err := client.ListBucketsWithContext(ctx, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Config
	s3Config := &aws.Config{
		Region:      aws.String("ap-northeast-1"), // set region aws
		Credentials: credentials.NewStaticCredentials(os.Getenv("KEY_ID"), os.Getenv("SECRET_KEY"), ""),
	}

	// Create new instance session
	sess := session.New(s3Config)

	// Create a instance from s3 use session above
	s3Client := s3.New(sess)

	// Create context
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Get list bucket use the `ListBuckets` above
	buckets, err := ListBuckets(ctx, s3Client)
	if err != nil {
		fmt.Printf("Couldn't retrieve list bucket: %v", err)
		return
	}

	// Show all buckets
	for _, bucket := range buckets.Buckets {
		fmt.Printf("Found bucket: %s, created at: %s\n", *bucket.Name, *bucket.CreationDate)
	}
}
