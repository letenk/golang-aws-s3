package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

func UploadFile(ctx context.Context, uploader *s3manager.Uploader, filePath string, bucketName string, fileName string) (*s3manager.UploadOutput, error) {
	// Read file in path
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Create object input
	input := &s3manager.UploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(file),
		ContentType: aws.String("image/jpg"),
	}

	// Upload to aws with context
	res, err := uploader.UploadWithContext(ctx, input)
	return res, err
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

	// Bucket name
	bucketName := os.Getenv("BUCKET_NAME")
	// Create new uploadert with session
	uploader := s3manager.NewUploader(sess)
	// Path file location
	path := "1.png"
	// File name
	filename := "1.png"

	// Create context
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Upload
	res, err := UploadFile(ctx, uploader, path, bucketName, filename)
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
		return
	}

	fmt.Println("Successfully uploaded file!")
	fmt.Printf("Response upload: %v", res)
}
