package bucket

import (
	"github.com/spf13/pflag"
	"github.com/vietanhduong/bucket-proxy/pkg/bucket/gcs"
)

const (
	namespace = "bucket"

	providerFlag = namespace + ".provider"
	nameFlag     = namespace + ".name"
)

func RegisterFlags(fs *pflag.FlagSet) {
	fs.String(providerFlag, string(GoogleCloudStorage), "Bucket provider must be one of: gcs")
	fs.String(nameFlag, "", "Bucket name")
	gcs.RegisterFlags(namespace, fs)
}
