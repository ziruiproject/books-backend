package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	s3sdk "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/rs/zerolog/log"
	"path/filepath"
	cfg "starter/config"
	"starter/internal/core/storage"
	"time"
)

type S3Storage struct {
	client        *s3sdk.Client
	presignClient *s3sdk.PresignClient
	region        string
}

func NewStorage() storage.Storage {
	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(cfg.AppConfig.StorageRootUser, cfg.AppConfig.StorageRootPassword, ""))
	config, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(cfg.AppConfig.StorageRegion))

	if err != nil {
		log.Panic().
			Err(err).
			Msg("unable to load AWS config")
		return nil
	}

	client := s3sdk.NewFromConfig(config, func(options *s3sdk.Options) {
		options.UsePathStyle = true
		options.BaseEndpoint = aws.String(fmt.Sprintf("http://%s:%s", cfg.AppConfig.StorageHost, cfg.AppConfig.StoragePortAPI))
	})
	presignClient := s3sdk.NewPresignClient(client)

	s3 := &S3Storage{
		client:        client,
		presignClient: presignClient,
		region:        cfg.AppConfig.StorageRegion,
	}

	return s3
}

func (s *S3Storage) Download(ctx context.Context, file storage.DownloadRequest) *storage.Response {
	if file.Name == "" {
		return nil
	}

	req := &s3sdk.GetObjectInput{
		Bucket:                     aws.String(file.Bucket),
		Key:                        aws.String(file.Name),
		ResponseContentDisposition: aws.String("inline"),
	}

	presignedURL, err := s.presignClient.PresignGetObject(ctx, req, s3sdk.WithPresignExpires(30*time.Minute))
	if err != nil {
		log.Error().Err(err).Msg("failed to generate presigned URL")
		return nil
	}

	return &storage.Response{
		Filename: file.Name,
		Link:     presignedURL.URL,
	}
}

func (s *S3Storage) Upload(ctx context.Context, file storage.UploadRequest) *storage.Response {
	if file.Name == "" {
		return nil
	}

	filename := fmt.Sprintf("book-%d%s", time.Now().UnixNano(), filepath.Ext(file.Name))
	_, err := s.client.PutObject(ctx, &s3sdk.PutObjectInput{
		Bucket: aws.String(file.Bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(file.Data),
		ACL:    types.ObjectCannedACLPrivate,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to upload file")
		return nil
	}

	return &storage.Response{
		Filename: filename,
	}
}

func (s *S3Storage) Delete(ctx context.Context, bucket, name string) *storage.Response {
	_, err := s.client.DeleteObject(ctx, &s3sdk.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to delete file")
		return nil
	}

	return &storage.Response{
		Filename: name,
	}
}
