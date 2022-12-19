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

func DeleteItem(ctx context.Context, sess *session.Session, bucketName string, fileName string) (*s3.DeleteObjectOutput, error) {
	// Creates a new instance of the S3 client with a session.
	svc := s3.New(sess)

	// Delete Object
	res, err := svc.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})

	if err != nil {
		return nil, err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})

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

	awsRegion := os.Getenv("AWS_REGION")

	// Config
	s3Config := &aws.Config{
		Region:      aws.String(awsRegion), // set region aws
		Credentials: credentials.NewStaticCredentials(os.Getenv("KEY_ID"), os.Getenv("SECRET_KEY"), ""),
	}

	// Create new instance session
	sess := session.New(s3Config)

	// Create context
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	bucketName := os.Getenv("BUCKET_NAME")
	fileName := "1.png"

	// Upload
	res, err := DeleteItem(ctx, sess, bucketName, fileName)

	if err != nil {
		fmt.Printf("Failed delete item, err: %v", res)
	}
	fmt.Println("Successfully delete file!")
}
