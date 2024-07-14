package server

import (
	"context"
	"net/http"
	"time"

	"github.com/vietanhduong/bucket-proxy/pkg/logging"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var log = logging.WithField("pkg", "pkg/server")

type Server struct {
	listenAddress string
	drainTimeout  time.Duration

	readHeaderTimeout time.Duration
	readTimeout       time.Duration
	writeTimeout      time.Duration
	maxHeaderBytes    int

	mux *http.ServeMux
}

func New(opt ...Option) *Server {
	s := defaultServer()

	for _, o := range opt {
		o(s)
	}
	return s
}

type RegisterHandler interface {
	HttpHandler() (string, http.Handler)
}

func (s *Server) RegisterHandler(handler RegisterHandler) {
	s.mux.Handle(handler.HttpHandler())
}

func (s *Server) initHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/healtz", HealthzHandler())
	mux.Handle("/", s.mux)
	return h2c.NewHandler(LoggingMiddleware()(mux), &http2.Server{})
}

func (s *Server) Run(stop <-chan struct{}) error {
	srv := &http.Server{
		Addr:              s.listenAddress,
		Handler:           s.initHandler(),
		ReadHeaderTimeout: s.readHeaderTimeout,
		ReadTimeout:       s.readTimeout,
		WriteTimeout:      s.writeTimeout,
		MaxHeaderBytes:    s.maxHeaderBytes,
	}

	errCh := make(chan error, 1)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.WithError(err).Errorf("Failed to server")
			errCh <- err
		}
	}()

	stopFn := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), s.drainTimeout)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	}

	log.Infof("HTTP server is being started at address %s...", s.listenAddress)
	for {
		select {
		case err := <-errCh:
			return err
		case <-stop:
			log.Info("HTTP server is being stopped...")
			return stopFn()
		}
	}
}
