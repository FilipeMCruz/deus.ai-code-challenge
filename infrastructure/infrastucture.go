// Package infrastructure is responsible for running the http server and defining a generic wrapper for logging, content type and panic recovery
package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"deus.ai-code-challenge/infrastructure/content"
	"deus.ai-code-challenge/infrastructure/logging"
	"deus.ai-code-challenge/infrastructure/recovery"
)

// Wrap wraps a handler with:
// - basic request info logging
// - basic content type header set to application/json
// - panic recovery, returns a 500
func Wrap(next http.Handler) http.Handler {
	next = recovery.WrapPanicRecovery(next)
	next = content.WrapJsonContentType(next)

	return logging.WrapLogging(next)
}

// Run runs an http server and ensures that it is gracefully shutdown:
// - in flight requests are answered
// - new requests are not accepted
func Run(ctx context.Context, stop func(), port int, handler http.Handler, started chan<- struct{}) error {
	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())
	defer stopOngoingGracefully()

	httpServer := &http.Server{
		Handler: handler,
		BaseContext: func(_ net.Listener) context.Context {
			return ongoingCtx
		},
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	started <- struct{}{}

	go func() {
		if err := httpServer.Serve(l); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	<-ctx.Done()
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return httpServer.Shutdown(shutdownCtx)
}
