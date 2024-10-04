package gcs

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/vietanhduong/bucket-proxy/pkg/bucket/types"
	"github.com/vietanhduong/bucket-proxy/pkg/config"
	"github.com/vietanhduong/bucket-proxy/pkg/logging"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var log = logging.WithField("pkg", "pkg/bucket/gcs")

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
	// if there is no error and the object is not deleted, return the object metadata
	if err == nil {
		// if the object is deleted, return nil
		if !attrs.Deleted.IsZero() {
			log.WithField("path", path).Trace("object deleted")
			return nil, nil
		}
		var ret types.ObjectMetadata
		ret.FromObjectAttrs(attrs)
		return &ret, nil
	}

	// if the object does not exist, check if it is a directory
	if !errors.Is(err, storage.ErrObjectNotExist) {
		return nil, fmt.Errorf("object attrs: %w", err)
	}

	it := c.gcs.Bucket(c.bucket).Objects(ctx, &storage.Query{Prefix: path})

	// get the first object in the directory, if it exists, it is a directory
	// otherwise, the input path doesn't exist
	if _, err = it.Next(); err != nil {
		if err == iterator.Done {
			log.WithField("path", path).Trace("object not found")
			return nil, nil
		}
		return nil, fmt.Errorf("object attrs: %w", err)
	}
	return &types.ObjectMetadata{
		Bucket:      c.bucket,
		Name:        path,
		IsDirectory: true,
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

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("object attrs: %w", err)
	}

	var r *storage.Reader
	if opts.Start == 0 && opts.Offset <= 0 {
		r, err = obj.NewReader(ctx)
	} else {
		r, err = obj.NewRangeReader(ctx, opts.Start, opts.Offset)
	}
	if err != nil {
		return nil, fmt.Errorf("new reader: %w", err)
	}

	resp := &types.DownloadResponse{Reader: r}
	resp.ObjectMetadata.FromObjectAttrs(attrs)
	resp.Size = r.Attrs.Size
	resp.ContentEncoding = r.Attrs.ContentEncoding
	return resp, nil
}
