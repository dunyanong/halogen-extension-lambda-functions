package main

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"log"
	"path/filepath"

	"github.com/yuin/gopher-lua"
)

func handler(ctx context.Context, s3Event events.S3Event) error {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	for _, record := range s3Event.Records {
		s3Entity := record.S3
		bucket := s3Entity.Bucket.Name
		key := s3Entity.Object.Key

		// Download the file from S3
		resp, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			log.Fatalf("Failed to get object: %v", err)
			return err
		}
		defer resp.Body.Close()

		// Unzip the file
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
		if err != nil {
			log.Fatalf("Failed to create zip reader: %v", err)
			return err
		}

		// Extracting and Reading the Lua File
		var luaCode string
		for _, f := range zr.File {
			if filepath.Ext(f.Name) == ".lua" {
				rc, err := f.Open() // Opens and reads the Lua file.
				if err != nil {
					log.Fatalf("Failed to open file inside zip: %v", err)
					return err
				}
				defer rc.Close()

				luaBytes, err := io.ReadAll(rc)
				if err != nil {
					log.Fatalf("Failed to read lua file: %v", err)
					return err
				}
				// Stores the Lua code as a string
				luaCode = string(luaBytes)
			}
		}

		// Logs an error if no Lua file is found :(
		if luaCode == "" {
			log.Fatalf("No Lua file found in the zip")
			return fmt.Errorf("no Lua file found in the zip")
		}

		// Compile Lua code using GopherLua
		L := lua.NewState()
		defer L.Close()
		if err := L.DoString(luaCode); err != nil {
			log.Fatalf("Failed to compile lua code: %v", err)
			return err
		}

		// Create byte code
		byteCode := L.Dump(func(proto *lua.Prototype) []byte {
			// This is a simplistic example of dumping byte code.
			// In practice, you may need a proper implementation.
			return proto.Source
		})

		// Upload byte code back to S3
		_, err = svc.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key + ".bytes"),
			Body:   bytes.NewReader(byteCode),
		})
		if err != nil {
			log.Fatalf("Failed to upload byte code: %v", err)
			return err
		}
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
