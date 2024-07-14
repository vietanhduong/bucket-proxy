package server

import (
	"net/http"
	"time"
)

type Option func(*Server)

func WithListenAddress(addr string) Option {
	return func(s *Server) {
		if addr != "" {
			s.listenAddress = addr
		}
	}
}

func WithDrainTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		if timeout > 0 {
			s.drainTimeout = timeout
		}
	}
}

func defaultServer() *Server {
	return &Server{
		listenAddress:     "0.0.0.0:8080",
		drainTimeout:      15 * time.Second,
		readHeaderTimeout: time.Second,
		readTimeout:       5 * time.Minute,
		writeTimeout:      5 * time.Minute,
		maxHeaderBytes:    1 << 20, // 1 MB
		mux:               http.NewServeMux(),
	}
}
