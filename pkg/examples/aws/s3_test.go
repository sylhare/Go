package localstack

import (
	"bytes"
	"context"
	"fmt"
	sdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
	"testing"
)

var (
	bucketName = "my-test-bucket"
	fileKey    = "testfile.txt"
)

func TestS3(t *testing.T) {
	ctx := context.Background()

	localstack, err := New()
	if err != nil {
		t.Fatalf("failed to start LocalStack: %s", err)
	}
	defer localstack.Terminate()
	staticCredentials := sdk.NewCredentialsCache(credentials.NewStaticCredentialsProvider("test", "test", ""))
	s3Client := s3.NewFromConfig(localstack.Config, func(o *s3.Options) {
		o.UsePathStyle = true
		o.Credentials = staticCredentials
	})

	t.Run("Create bucket", func(t *testing.T) {
		_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: &bucketName,
		})
		if err != nil {
			t.Fatalf("unable to create bucket, %v", err)
		}
		fmt.Printf("Successfully created bucket %s\n", bucketName)
	})

	t.Run("List buckets", func(t *testing.T) {
		result, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
		if err != nil {
			log.Fatalf("unable to list buckets, %v", err)
		}

		fmt.Println("Buckets:")
		for _, b := range result.Buckets {
			fmt.Printf("* %s\n", *b.Name)
		}
	})

	t.Run("Upload file", func(t *testing.T) {
		fileContent := "Hello, World!"

		_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: &bucketName,
			Key:    &fileKey,
			Body:   bytes.NewReader([]byte(fileContent)),
		})
		if err != nil {
			t.Fatalf("unable to upload file, %v", err)
		}
		fmt.Printf("Successfully uploaded file %s\n", fileKey)
	})

	t.Run("Retrieve file", func(t *testing.T) {
		getObjectOutput, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucketName,
			Key:    &fileKey,
		})
		if err != nil {
			t.Fatalf("unable to retrieve file, %v", err)
		}
		defer getObjectOutput.Body.Close()

		retrievedContent, err := io.ReadAll(getObjectOutput.Body)
		if err != nil {
			t.Fatalf("unable to read file content, %v", err)
		}
		fmt.Printf("Successfully retrieved file content: %s\n", string(retrievedContent))
	})
}
