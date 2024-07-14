package types

import (
	"io"
	"time"
)

type DownloadOptions struct {
	AcceptCompress bool
	Start, Offset  int64
}

type ObjectMetadata struct {
	Bucket             string
	Name               string
	ContentType        string
	ContentLanguage    string
	Size               int64
	ContentEncoding    string
	ContentDisposition string
	CacheControl       string
	Created            time.Time
	Updated            time.Time
}

type DownloadResponse struct {
	io.Reader
	Size            int64
	ContentEncoding string
}
