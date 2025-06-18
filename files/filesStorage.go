package files

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
)

type MinIOStorage struct {
	Client     *minio.Client
	BucketName string
}

func NewMinIOStorage(endpoint, accessKey, secretKey, bucket string) (*MinIOStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &MinIOStorage{
		Client:     client,
		BucketName: bucket,
	}, nil
}

func (s *MinIOStorage) UploadFile(ctx context.Context, filename string, data []byte) (string, error) {
	reader := bytes.NewReader(data)
	_, err := s.Client.PutObject(ctx, s.BucketName, filename, reader, int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://localhost:9000/%s/%s", s.BucketName, filename), nil
}

func (s *MinIOStorage) DownloadFile(ctx context.Context, objectName string) ([]byte, error) {
	object, err := s.Client.GetObject(ctx, s.BucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from MinIO: %w", err)
	}
	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}

	return data, nil
}
