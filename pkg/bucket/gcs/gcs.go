package gcs

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/vietanhduong/bucket-proxy/pkg/bucket/types"
	"github.com/vietanhduong/bucket-proxy/pkg/config"
	"google.golang.org/api/option"
)

type Client struct {
	gcs    *storage.Client
	opts   []option.ClientOption
	bucket string
}

func NewClient(bucket string, opt ...Option) (*Client, error) {
	var c Client
	c.bucket = bucket
	for _, o := range opt {
		o(&c)
	}
	c.opts = append(c.opts, option.WithUserAgent(config.UserAgent()))
	var err error

	if c.gcs, err = storage.NewClient(context.Background(), c.opts...); err != nil {
		return nil, fmt.Errorf("new google cloud storage client: %w", err)
	}
	return &c, nil
}

func (c *Client) ObjectMetadata(ctx context.Context, path string) (*types.ObjectMetadata, error) {
	if path == "" {
		return nil, fmt.Errorf("empty path")
	}
	obj := c.gcs.Bucket(c.bucket).Object(path)
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("object attrs: %w", err)
	}

	if !attrs.Deleted.IsZero() {
		return nil, nil
	}

	return &types.ObjectMetadata{
		Bucket:             attrs.Bucket,
		Name:               attrs.Name,
		Size:               attrs.Size,
		ContentType:        attrs.ContentType,
		ContentLanguage:    attrs.ContentLanguage,
		ContentEncoding:    attrs.ContentEncoding,
		ContentDisposition: attrs.ContentDisposition,
		CacheControl:       attrs.CacheControl,
		Created:            attrs.Created,
		Updated:            attrs.Updated,
	}, nil
}

func (c *Client) Download(ctx context.Context, path string, opts types.DownloadOptions) (*types.DownloadResponse, error) {
	if path == "" {
		return nil, fmt.Errorf("empty path")
	}
	obj := c.gcs.Bucket(c.bucket).Object(path)
	if opts.AcceptCompress {
		obj = obj.ReadCompressed(true)
	}

	var r *storage.Reader
	var err error

	if opts.Start == 0 && opts.Offset <= 0 {
		r, err = obj.NewReader(ctx)
	} else {
		r, err = obj.NewRangeReader(ctx, opts.Start, opts.Offset)
	}
	if err != nil {
		return nil, fmt.Errorf("new reader: %w", err)
	}

	resp := &types.DownloadResponse{
		Reader:          r,
		Size:            r.Attrs.Size,
		ContentEncoding: r.Attrs.ContentEncoding,
	}
	return resp, nil
}
