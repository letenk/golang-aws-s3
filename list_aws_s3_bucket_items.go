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

func ListItems(ctx context.Context, client *s3.S3, bucketName string, prefix string) (*s3.ListObjectsV2Output, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix),
	}

	res, err := client.ListObjectsV2WithContext(ctx, input)

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

	// Bucket name
	bucketName := os.Getenv("BUCKET_NAME")
	// Prefix name
	prefixName := ""

	// Create context
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Get list item in buckets use `ListItems` above
	bucketObjects, err := ListItems(ctx, s3Client, bucketName, prefixName)
	if err != nil {
		fmt.Printf("Couldn't retrieve bucket items: %v", err)
		return
	}

	// Show all list item
	for _, item := range bucketObjects.Contents {
		fmt.Printf("Name: %s, Last Modified: %s\n", *item.Key, *item.LastModified)
	}

}
