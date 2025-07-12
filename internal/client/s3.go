package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/VI-IM/im_backend_go/shared/logger"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucket        string
}

type S3ClientInterface interface {
	UploadFile(ctx context.Context, key string, body io.Reader) (string, error)
	GetFile(ctx context.Context, key string) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, key string) error
	GeneratePresignedURL(ctx context.Context, key string, operation string, duration time.Duration) (string, error)
}

func NewS3Client(bucket string) (S3ClientInterface, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	// If no bucket is provided, use the one from environment
	if bucket == "" {
		bucket = os.Getenv("AWS_S3_BUCKET")
		if bucket == "" {
			return nil, fmt.Errorf("AWS_S3_BUCKET environment variable is not set")
		}
	}

	// Check if we're using LocalStack
	endpointURL := os.Getenv("AWS_ENDPOINT_URL")
	var client *s3.Client

	if endpointURL != "" {
		// Configure for LocalStack
		client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpointURL)
			o.UsePathStyle = true // LocalStack requires path-style addressing
		})
		logger.Get().Info().Msgf("S3 client configured for LocalStack endpoint: %s", endpointURL)
	} else {
		// Use default AWS configuration
		client = s3.NewFromConfig(cfg)
	}

	presignClient := s3.NewPresignClient(client)

	return &S3Client{
		client:        client,
		presignClient: presignClient,
		bucket:        bucket,
	}, nil
}

func (c *S3Client) UploadFile(ctx context.Context, key string, body io.Reader) (string, error) {
	// Read all content from the reader
	content, err := io.ReadAll(body)
	if err != nil {
		logger.Get().Error().Err(err).Msg("failed to read file content")
		return "", err
	}

	// Upload to S3 with the content
	_, err = c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(content),
		ContentType: aws.String("application/pdf"),
	})
	if err != nil {
		logger.Get().Error().Err(err).Msg("failed to upload to S3")
		return "", err
	}

	// Return the appropriate URL based on environment
	endpointURL := os.Getenv("AWS_ENDPOINT_URL")
	if endpointURL != "" {
		// LocalStack URL format
		return fmt.Sprintf("%s/%s/%s", endpointURL, c.bucket, key), nil
	}

	// Return the public S3 URL for AWS
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", c.bucket, key), nil
}

func (c *S3Client) GetFile(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

func (c *S3Client) DeleteFile(ctx context.Context, key string) error {
	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (c *S3Client) GeneratePresignedURL(ctx context.Context, key string, operation string, duration time.Duration) (string, error) {

	switch operation {
	case "GET":
		request, err := c.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(c.bucket),
			Key:    aws.String(key),
		}, s3.WithPresignExpires(duration))
		if err != nil {
			logger.Get().Error().Err(err).Msg("failed to generate presigned GET URL")
			return "", err
		}
		return request.URL, nil

	case "PUT":
		request, err := c.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(c.bucket),
			Key:         aws.String(key),
			ContentType: aws.String("application/pdf"),
		}, s3.WithPresignExpires(duration))
		if err != nil {
			logger.Get().Error().Err(err).Msg("failed to generate presigned PUT URL")
			return "", err
		}
		return request.URL, nil

	default:
		return "", fmt.Errorf("unsupported operation: %s. Use 'GET' or 'PUT'", operation)
	}
}
