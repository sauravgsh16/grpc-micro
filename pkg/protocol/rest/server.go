package rest

import (
	"context"
	"github.com/sauravgsh16/api-grpc/pkg/api/v1"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RunServer runs HTTP/REST gateway
func RunServer(ctx context.Context, grpcPort, httpPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := v1.RegisterToDoServiceHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts); err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: mux,
	}

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {

		}

		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		srv.Shutdown(ctx)
	}()

	log.Println("Starting HTTP/REST server")
	return srv.ListenAndServe()
}
