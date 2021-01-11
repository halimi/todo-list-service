package main

import (
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes"
	"github.com/halimi/todo-list-service/todolistpb"
	"google.golang.org/grpc"
)

func createTodo(c todolistpb.TodoListServiceClient) {
	fmt.Println("Creating Todo")
	todo := &todolistpb.Todo{
		Title:   "First Todo",
		Note:    "This is a test",
		DueDate: ptypes.TimestampNow(),
	}
	res, err := c.CreateTodo(context.Background(), &todolistpb.CreateTodoRequest{Todo: todo})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	resTodo := res.GetTodo()
	fmt.Println("Todo has been created:")
	fmt.Println("  Id:", resTodo.GetId())
	fmt.Println("  Title:", resTodo.GetTitle())
	fmt.Println("  Note:", resTodo.GetNote())
	fmt.Println("  Due date:", ptypes.TimestampString(resTodo.GetDueDate()))
}

func main() {
	fmt.Println("Todo List Client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:5000", opts)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := todolistpb.NewTodoListServiceClient(cc)

	createTodo(c)
}
