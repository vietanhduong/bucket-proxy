package types

import (
	"errors"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

var ErrEmptyPath = errors.New("empty path")

type DownloadOptions struct {
	AcceptCompress bool
	Start, Offset  int64
}

type ObjectMetadata struct {
	Bucket             string
	Name               string
	IsDirectory        bool
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
	ObjectMetadata
}

func (o *ObjectMetadata) FromObjectAttrs(in *storage.ObjectAttrs) {
	if o == nil || in == nil {
		return
	}
	o.Bucket = in.Bucket
	o.Name = in.Name
	o.Size = in.Size
	o.ContentType = in.ContentType
	o.ContentLanguage = in.ContentLanguage
	o.ContentEncoding = in.ContentEncoding
	o.ContentDisposition = in.ContentDisposition
	o.CacheControl = in.CacheControl
	o.Created = in.Created
	o.Updated = in.Updated
}
