package main

import (
	"context"
	"log"

	"github.com/koralle/go-web-application-development/go_todo_app/server"
)

func main() {

	if err := server.Run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
	}

}
