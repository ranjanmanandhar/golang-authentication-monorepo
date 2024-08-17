package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"gitlab-server.wlink.com.np/nettv/nettv-auth/authenticator/internal/server"
	"golang.org/x/sync/errgroup"
)

// @title   Service Name API documentation
// @version 1.0.0

// @contact.name  Prabesh Lamichhane Magar
// @contact.email prab.magar@gmail.com

// @host     localhost:8080
// @BasePath /api/v1
func main() {
	flag.Parse()

	os.Exit(start())
}

func start() int {
	ctx := context.Background()
	port := 8080
	server := server.New(server.NewServerOption{
		Port:              port,
		Host:              "",
		ReadTimeOut:       5,
		ReadHeaderTimeout: 5,
		WriteTimeout:      5,
		IdleTimeout:       5,
	})

	var eg errgroup.Group
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	eg.Go(func() error {
		<-ctx.Done()

		if err := server.Stop(); err != nil {
			return err
		}

		return nil
	})

	if err := server.Start(); err != nil {
		return 1
	}

	if err := eg.Wait(); err != nil {
		return 1
	}

	return 0
}
