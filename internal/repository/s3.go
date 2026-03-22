package repository

import (
	"bytes"
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const bucketName = "avatars"

type MinIOStorage struct {
	client *minio.Client
}

func NewMinIOStorage(endpoint, accessKey, secretKey string) *MinIOStorage {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		panic(fmt.Sprintf("minio init failed: %v", err))
	}

	policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::avatars/*"]}]}`
	_ = client.SetBucketPolicy(context.Background(), bucketName, policy)

	return &MinIOStorage{client: client}
}

func (s *MinIOStorage) Upload(ctx context.Context, userID string, data []byte, contentType string) (string, error) {
	objectName := fmt.Sprintf("%s/avatar.jpg", userID)
	_, err := s.client.PutObject(ctx, bucketName, objectName,
		bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/%s/%s", bucketName, objectName), nil
}

func (s *MinIOStorage) Delete(ctx context.Context, userID string) error {
	return s.client.RemoveObject(ctx, bucketName,
		fmt.Sprintf("%s/avatar.jpg", userID),
		minio.RemoveObjectOptions{},
	)
}
