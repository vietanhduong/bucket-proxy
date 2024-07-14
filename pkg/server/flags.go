package server

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	namespace = "server"

	addressFlag      = namespace + ".address"
	drainTimeoutFlag = namespace + ".drain-timeout"
)

func RegisterFlags(fs *pflag.FlagSet) {
	defaultValues := defaultServer()
	fs.String(addressFlag, defaultValues.listenAddress, "Server listen address")
	fs.Duration(drainTimeoutFlag, defaultValues.drainTimeout, "Server drain timeout")
}

func InitFromViper(v *viper.Viper) *Server {
	return New(
		WithListenAddress(v.GetString(addressFlag)),
		WithDrainTimeout(v.GetDuration(drainTimeoutFlag)),
	)
}
