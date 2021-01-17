package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/golang/protobuf/ptypes"
	"github.com/halimi/todo-list-service/todolistpb"
	"google.golang.org/grpc"
)

func printTodo(t *todolistpb.Todo) {
	fmt.Println("Todo:")
	fmt.Println("  Id:", t.GetId())
	fmt.Println("  Title:", t.GetTitle())
	fmt.Println("  Note:", t.GetNote())
	fmt.Println("  Due date:", ptypes.TimestampString(t.GetDueDate()))
}

func createTodo(c todolistpb.TodoListServiceClient) int32 {
	fmt.Println("Creating Todo")
	todo := &todolistpb.Todo{
		Title:   "First Todo",
		Note:    "This is a test",
		DueDate: ptypes.TimestampNow(),
	}
	res, err := c.CreateTodo(context.Background(), &todolistpb.CreateTodoRequest{Todo: todo})
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
	resTodo := res.GetTodo()
	printTodo(resTodo)

	return resTodo.GetId()
}

func readTodo(c todolistpb.TodoListServiceClient, id int32) {
	fmt.Println("Reading Todo")

	res, err := c.ReadTodo(context.Background(), &todolistpb.ReadTodoRequest{TodoId: id})
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
	resTodo := res.GetTodo()
	printTodo(resTodo)
}

func updateTodo(c todolistpb.TodoListServiceClient, id int32) {
	fmt.Println("Updating Todo")
	todo := &todolistpb.Todo{
		Id:      id,
		Title:   "Updated Todo",
		Note:    "This is an updated test",
		DueDate: ptypes.TimestampNow(),
	}

	res, err := c.UpdateTodo(context.Background(), &todolistpb.UpdateTodoRequest{Todo: todo})
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
	resTodo := res.GetTodo()
	printTodo(resTodo)
}

func deleteTodo(c todolistpb.TodoListServiceClient, id int32) {
	fmt.Println("Deleting Todo")

	_, err := c.DeleteTodo(context.Background(), &todolistpb.DeleteTodoRequest{TodoId: id})
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}

	fmt.Println("Successfully deleted:", id)
}

func listTodos(c todolistpb.TodoListServiceClient) {
	fmt.Println("Listing Todos")

	stream, err := c.ListTodos(context.Background(), &todolistpb.ListTodosRequest{})
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Stream error: %v", err)
		}
		printTodo(res.GetTodo())
	}
}

func main() {
	fmt.Println("Todo List Client")

	host := flag.String("host", "localhost", "Service host name")
	port := flag.String("port", "5000", "Service port number")

	flag.Parse()

	connStr := fmt.Sprintf("%v:%v", *host, *port)
	opts := grpc.WithInsecure()

	cc, err := grpc.Dial(connStr, opts)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := todolistpb.NewTodoListServiceClient(cc)

	id := createTodo(c)

	readTodo(c, id)

	updateTodo(c, id)

	listTodos(c)

	deleteTodo(c, id)
}
