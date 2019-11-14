package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	v1 "github.com/sauravgsh16/api-grpc/pkg/api/v1"
)

// RunServer runs gRPC service to publish ToDo service
func RunServer(ctx context.Context, v1API v1.ToDoServiceServer, port string) error {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	v1.RegisterToDoServiceServer(server, v1API)

	// graceful shutdown
	c := make(chan os.Signal, 1)

	// Will relay all signal from os.Interrupt to chan c
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			log.Println("shutting down gRPC server....")
			server.GracefulStop()

			// Send done signal to context
			<-ctx.Done()
		}
	}()

	log.Println("Starting gRPC server....")
	return server.Serve(ln)
}
