package rendered_label_storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
	"io"
	"time"
)

type S3 struct {
	Logger        *zap.Logger
	Client        *s3.Client
	PresignClient *s3.PresignClient
	BucketName    string
}

func (s *S3) Create(key string, data io.Reader) error {
	_, err := s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
		Body:   data,
	})
	if err != nil {
		s.Logger.Error("Failed creating label in S3.", zap.Error(err))
		return fmt.Errorf("label creation failed: %w", err)
	}

	return nil
}

func (s *S3) GetDownloadUrl(key string, validity time.Duration) (string, error) {
	expireReq := func(opts *s3.PresignOptions) {
		opts.Expires = validity
	}
	req, err := s.PresignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	}, expireReq)
	if err != nil {
		s.Logger.Error("Cannot create presigned URL for rendered image download.",
			zap.String("key", key),
			zap.Error(err))
		return "", fmt.Errorf("presigned URL creation failed: %w", err)
	}
	return req.URL, nil
}
