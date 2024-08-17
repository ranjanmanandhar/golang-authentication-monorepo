package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"gitlab-server.wlink.com.np/nettv/nettv-auth/authenticator/internal/route"
)

type Server struct {
	address string
	server  *http.Server
}

type NewServerOption struct {
	Host              string
	Port              int
	ReadTimeOut       int
	ReadHeaderTimeout int
	WriteTimeout      int
	IdleTimeout       int
}

func New(opts NewServerOption) *Server {
	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))
	return &Server{
		address: address,
		server: &http.Server{
			Addr:              address,
			Handler:           route.SetupRoutes(),
			ReadTimeout:       time.Duration(opts.ReadTimeOut) * time.Second,
			ReadHeaderTimeout: time.Duration(opts.ReadHeaderTimeout) * time.Second,
			WriteTimeout:      time.Duration(opts.WriteTimeout) * time.Second,
			IdleTimeout:       time.Duration(opts.IdleTimeout) * time.Second,
		},
	}
}

func (s *Server) Start() error {

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("error staring server: %w", err)
	}

	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error stopping server: %w", err)
	}

	return nil
}
