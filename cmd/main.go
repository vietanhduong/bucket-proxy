package main

import (
	"fmt"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vietanhduong/bucket-proxy/pkg/bucket"
	"github.com/vietanhduong/bucket-proxy/pkg/logging"
	"github.com/vietanhduong/bucket-proxy/pkg/proxy"
	"github.com/vietanhduong/bucket-proxy/pkg/server"
)

func newCommand() *cobra.Command {
	v := viper.New()
	cmd := &cobra.Command{
		Use:   "bucket-proxy",
		Short: "Proxy for cloud storage bucket",
		RunE: func(cmd *cobra.Command, args []string) error {
			if printVersion(cmd) {
				return nil
			}

			logging.InitFromViper(v)

			ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
			defer cancel()

			b, err := bucket.InitFromViper(v)
			if err != nil {
				return fmt.Errorf("init bucket: %w", err)
			}
			p := proxy.New(b)
			s := server.InitFromViper(v)
			s.RegisterHandler(p)
			return s.Run(ctx.Done())
		},
	}

	addFlags(v, cmd, bucket.RegisterFlags, logging.RegisterFlags, server.RegisterFlags)
	return cmd
}

func main() { execute(newCommand()) }
