package main

import (
	"context"
	"flag"
	"log"
	"time"

	v1 "github.com/sauravgsh16/api-grpc/pkg/api/v1"
	"google.golang.org/grpc"

	"github.com/golang/protobuf/ptypes"
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
		log.Fatalf("Create request failed: %v\n", err)
	}
	log.Printf("Create Result: %v\n", respCreate.GetId())

	// Read call
	reqRead := &v1.ReadRequest{
		Api: apiVersion,
		Id:  2,
	}
	respRead, err := c.Read(ctx, reqRead)
	if err != nil {
		log.Fatalf("Read request failed: %v\n", err)
	}
	log.Printf("Read result: %v\n", respRead)

	// Update call
	reqUpdate := &v1.UpdateRequest{
		Api: apiVersion,
		Todo: &v1.ToDo{
			Id:          respRead.Todo.Id,
			Title:       respRead.Todo.Title,
			Description: respRead.Todo.Description + " updated",
			Reminder:    respRead.Todo.Reminder,
		},
	}

	respUpdate, err := c.Update(ctx, reqUpdate)
	if err != nil {
		log.Fatalf("Failed to update: %v\n", err)
	}
	log.Printf("Updated result: %v\n", respUpdate)

	// ReadAll call
	reqReadAll := &v1.ReadAllRequest{
		Api: apiVersion,
	}
	respReadAll, err := c.ReadAll(ctx, reqReadAll)
	if err != nil {
		log.Fatalf("Read all request failed: %v\n", err)
	}
	log.Println("Read all response :")
	for _, todo := range respReadAll.GetToDos() {
		log.Printf("%v\n", todo)
	}

	// Delete call
	reqDel := &v1.DeleteRequest{
		Api: apiVersion,
		Id:  1,
	}
	respDel, err := c.Delete(ctx, reqDel)
	if err != nil {
		log.Fatalf("Delete request failed: %v\n", err)
	}
	log.Printf("Delete resp: %v\n", respDel)
}
