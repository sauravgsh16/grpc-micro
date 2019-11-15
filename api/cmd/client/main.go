package main

import (
	"context"
	"flag"
	"log"
	"time"

	v1 "github.com/sauravgsh16/api-grpc/pkg/api/v1"
	"google.golang.org/grpc"
)

const (
	apiVersion = "v1"
)

func main() {
	address := flag.String("server", "", "gRPC server in format 'host:port'")
	flag.Parse()

	// Connect to gRPC server
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to gRPC server: %v", err)
	}
	defer conn.Close()

	c := v1.NewToDoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	/*
		t := time.Now().In(time.UTC)
		reminder, _ := ptypes.TimestampProto(t)
		pfx := t.Format(time.RFC3339Nano)


		// Create call
		reqCreate := &v1.CreateRequest{
			Api: apiVersion,
			ToDo: &v1.ToDo{
				Title:       "title (" + pfx + ")",
				Description: "description (" + pfx + ")",
				Reminder:    reminder,
			},
		}
		respCreate, err := c.Create(ctx, reqCreate)
		if err != nil {
			log.Fatalf("Create request failed: %v", err)
		}
		log.Printf("Create Result: %v", respCreate.GetId())
	*/

	// Read call
	reqRead := &v1.ReadRequest{
		Api: apiVersion,
		Id:  1,
	}
	respRead, err := c.Read(ctx, reqRead)
	if err != nil {
		log.Fatalf("Read request failed: %v", err)
	}
	log.Printf("Read result: %v", respRead)
}
