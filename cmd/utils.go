package main

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vietanhduong/bucket-proxy/pkg/config"
)

const versionFlag = "version"

type RegisterFunc func(fs *pflag.FlagSet)

func addFlags(v *viper.Viper, cmd *cobra.Command, reg ...RegisterFunc) (*viper.Viper, *cobra.Command) {
	cmd.PersistentFlags().BoolP(versionFlag, "v", false, "Print version and exit")
	for _, r := range reg {
		r(cmd.Flags())
	}
	setupViper(v)
	v.BindPFlags(cmd.Flags())
	return v, cmd
}

func setupViper(v *viper.Viper) {
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
}

func execute(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func printVersion(cmd *cobra.Command) bool {
	if ok, _ := cmd.PersistentFlags().GetBool(versionFlag); ok {
		config.PrintVersion()
		return true
	}
	return false
}
