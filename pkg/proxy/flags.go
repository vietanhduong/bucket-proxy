package proxy

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vietanhduong/bucket-proxy/pkg/bucket"
)

const (
	namespace = "proxy"

	indexPageFlag    = namespace + ".index-page"
	notFoundPageFlag = namespace + ".not-found-page"
	webModeFlag      = namespace + ".web-mode"
)

func RegisterFlags(fs *pflag.FlagSet) {
	fs.Bool(webModeFlag, false, "Enable web mode. If enabled, the proxy will serve index and not found pages")
	fs.String(indexPageFlag, "index.html", "Index page. Only used when web mode is enabled")
	fs.String(notFoundPageFlag, "404.html", "Not found page. Only used when web mode is enabled")
}

func InitFromViper(bucket bucket.Interface, v *viper.Viper) *Proxy {
	return New(bucket,
		WithWebMode(v.GetBool(webModeFlag)),
		WithIndexPage(v.GetString(indexPageFlag)),
		WithNotFoundPage(v.GetString(notFoundPageFlag)),
	)
}
