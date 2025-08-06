package storage

import "context"

type Storage interface {
	Upload(ctx context.Context, file UploadRequest) *Response
	Download(ctx context.Context, file DownloadRequest) *Response
	Delete(ctx context.Context, bucket string, name string) *Response
}
