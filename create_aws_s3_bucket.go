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

func CreateBucket(ctx context.Context, client *s3.S3, bucketName string) (*s3.CreateBucketOutput, error) {
	bucket := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}

	res, err := client.CreateBucketWithContext(ctx, bucket)
	return res, err
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

	// Set name for create bucket
	bucketName := "learn-golang-aws-s3"

	// Create bucket use function `CreateBucket` above
	res, err := CreateBucket(ctx, s3Client, bucketName)
	if err != nil {
		fmt.Printf("Couldn't create bucket: %v", err)
		return
	}

	fmt.Println("Successfully created bucket")
	fmt.Printf("New bucket: %v", res)
}
