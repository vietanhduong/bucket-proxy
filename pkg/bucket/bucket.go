package bucket

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"github.com/vietanhduong/bucket-proxy/pkg/bucket/gcs"
	"github.com/vietanhduong/bucket-proxy/pkg/bucket/types"
)

type Interface interface {
	Download(ctx context.Context, path string, opts types.DownloadOptions) (*types.DownloadResponse, error)
	ObjectMetadata(ctx context.Context, path string) (*types.ObjectMetadata, error)
}

type Provider string

const (
	GoogleCloudStorage Provider = "gcs"
)

func InitFromViper(v *viper.Viper) (Interface, error) {
	provider := Provider(v.GetString(providerFlag))
	bucket := v.GetString(nameFlag)
	if bucket == "" {
		return nil, fmt.Errorf("bucket name is required")
	}

	switch provider {
	case GoogleCloudStorage:
		return gcs.InitWithViper(bucket, v)
	default:
		return nil, fmt.Errorf("unknown provider %q", provider)
	}
}
