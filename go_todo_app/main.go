package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/koralle/go-web-application-development/go_todo_app/server"
)

func main() {

	if len(os.Args) != 2 {
		log.Printf("need port number\n")
		os.Exit(1)
	}

	port := os.Args[1]

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen port %s: %v", port, err)
	}

	if err := server.Run(context.Background(), l); err != nil {
		log.Printf("failed to terminate server: %v", err)
	}

}
