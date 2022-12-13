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
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

func DownloadFile(ctx context.Context, downloader *s3manager.Downloader, bucketName string, key string) error {
	// Save file on directory
	file, err := os.Create(key)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create object input
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}

	_, err = downloader.DownloadWithContext(ctx, file, input)
	return err
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

	// Bucket name
	bucketName := os.Getenv("BUCKET_NAME")

	// Create new downloader with session
	downloader := s3manager.NewDownloader(sess)
	// Path download file location
	path := "1.png"

	// Create context
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Download
	err = DownloadFile(ctx, downloader, bucketName, path)
	if err != nil {
		fmt.Printf("Couldn't download file: %v", err)
		return
	}

	fmt.Println("Successfully downloaded file")

}
