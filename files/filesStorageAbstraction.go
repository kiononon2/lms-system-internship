package files

import "context"

type FileStorage interface {
	UploadFile(ctx context.Context, filename string, data []byte) (string, error)
	DownloadFile(ctx context.Context, fileURL string) ([]byte, error)
}
