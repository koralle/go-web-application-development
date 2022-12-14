package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/koralle/go-web-application-development/go_todo_app/config"
	"github.com/koralle/go-web-application-development/go_todo_app/mux"
	"github.com/koralle/go-web-application-development/go_todo_app/server"
)

func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}

	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

	mux, cleanup, err := mux.NewMux(ctx, cfg)

	if err != nil {
		return err
	}

	defer cleanup()

	s := server.NewServer(l, mux)

	return s.Run(ctx)
}

func main() {

	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}

}
