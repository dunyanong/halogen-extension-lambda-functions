package main

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func handler(ctx context.Context, event events.S3Event) error {
	// Initialize AWS session
	sess := session.Must(session.NewSession())

	// Create an S3 service client
	svc := s3.New(sess)

	for _, record := range event.Records {
		// Retrieve bucket and object key from the S3 event
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.Key

		// Download the file from S3
		downloadedFile, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("failed to download file from S3 (bucket: %s, key: %s): %w", bucket, key, err)
		}

		// Create a buffer to store the contents of the downloaded file
		buffer := new(bytes.Buffer)
		if _, err := io.Copy(buffer, downloadedFile.Body); err != nil {
			return fmt.Errorf("failed to copy downloaded file to buffer (bucket: %s, key: %s): %w", bucket, key, err)
		}

		// Unzip the contents of the file
		zipReader, err := zip.NewReader(bytes.NewReader(buffer.Bytes()), int64(buffer.Len()))
		if err != nil {
			return fmt.Errorf("failed to create zip reader for file (bucket: %s, key: %s): %w", bucket, key, err)
		}

		// Extract and save each file from the zip archive
		for _, file := range zipReader.File {
			// Open the file from the zip archive
			zippedFile, err := file.Open()
			if err != nil {
				return fmt.Errorf("failed to open zipped file %s from archive (bucket: %s, key: %s): %w", file.Name, bucket, key, err)
			}
			defer zippedFile.Close()

			// Create the destination file
			destFilePath := filepath.Join("/tmp", file.Name)
			destFile, err := os.Create(destFilePath)
			if err != nil {
				return fmt.Errorf("failed to create destination file %s (bucket: %s, key: %s): %w", destFilePath, bucket, key, err)
			}
			defer destFile.Close()

			// Copy the contents of the zipped file to the destination file
			if _, err := io.Copy(destFile, zippedFile); err != nil {
				return fmt.Errorf("failed to copy contents to destination file %s (bucket: %s, key: %s): %w", destFilePath, bucket, key, err)
			}

			fmt.Printf("File extracted: %s\n", destFilePath)
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
