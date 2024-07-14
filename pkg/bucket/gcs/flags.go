package gcs

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	namespace           = "gcs"
	credentialsFileFlag = namespace + ".credentials-file"
)

func RegisterFlags(namespace string, fs *pflag.FlagSet) {
	if namespace != "" {
		namespace = fmt.Sprintf("%s.", strings.TrimPrefix(namespace, "."))
	}

	fs.String(namespace+credentialsFileFlag, "", "Path to the credentials file")
}

func InitWithViper(bucket string, v *viper.Viper) (*Client, error) {
	credentialsFile := v.GetString(credentialsFileFlag)
	return NewClient(bucket, WithCredentialsFile(credentialsFile))
}
