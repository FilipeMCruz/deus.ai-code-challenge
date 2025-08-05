// Package main is responsible for:
// - reading the flags passed to the program
// - registering the endpoints
// - starting the server in the defined port
package main

import (
	"context"
	"deus.ai-code-challenge/api"
	"deus.ai-code-challenge/infrastructure"
	"deus.ai-code-challenge/repository"
	"flag"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	err := start(ctx, stop, *port)
	if err != nil {
		log.Fatal(err)
	}
}

// start registers the handlers (wrapped with logging) in a ServeMux
// and calls infrastructure.Run to run the http Server
func start(ctx context.Context, stop func(), port int) error {
	repo := repository.NewVisitsInMemoryRepository()

	mux := http.NewServeMux()

	for url, handler := range api.Handlers(repo) {
		mux.Handle(url, infrastructure.Wrap(handler))
	}

	return infrastructure.Run(ctx, stop, port, mux)
}
