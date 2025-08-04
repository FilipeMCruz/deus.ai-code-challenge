// Package infrastructure is responsible for running the http server and defining a generic wrapper for logging
package infrastructure

import (
	"context"
	"deus.ai-code-challenge/infrastructure/logging"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

// Wrap wraps a handler with:
// - basic request info logging
func Wrap(next http.Handler) http.Handler {
	return logging.WrapLogging(next)
}

// Run runs an http server and ensures that it is gracefully shutdown:
// - in flight requests are answered
// - new requests are not accepted
func Run(ctx context.Context, stop func(), port int, handler http.Handler) error {
	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
		BaseContext: func(_ net.Listener) context.Context {
			return ongoingCtx
		},
	}

	go func() {
		log.Printf("deus.ai server starting on port %d", port)
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	<-ctx.Done()
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	defer stopOngoingGracefully()

	return httpServer.Shutdown(shutdownCtx)
}
