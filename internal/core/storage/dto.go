package storage

type UploadRequest struct {
	Bucket string
	Name   string
	Data   []byte
}

type DownloadRequest struct {
	Bucket string
	Name   string
}

type Response struct {
	Filename string `json:"file_name,omitempty"`
	Link     string `json:"link,omitempty"`
}
