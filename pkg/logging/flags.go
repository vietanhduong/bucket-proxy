package logging

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	namespace    = "log"
	levelFlag    = namespace + ".level"
	formatFlag   = namespace + ".format"
	disableQuote = namespace + ".disable-quote"
)

func RegisterFlags(fs *pflag.FlagSet) {
	fs.String(levelFlag, "info", "Log level")
	fs.String(formatFlag, "text", "Log format")
	fs.Bool(disableQuote, false, "Disable quote, this option only works for text format.")
}

func InitFromViper(v *viper.Viper) {
	SetLevel(v.GetString(levelFlag))
	SetFormatter(Formatter(v.GetString(formatFlag)), v.GetBool(disableQuote))
}
