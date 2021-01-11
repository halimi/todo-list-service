package server

import (
	"context"
	"fmt"

	"github.com/halimi/todo-list-service/todolistpb"
)

type Server struct {
	// Implementing TodoListServiceServer interface
}

func (*Server) CreateTodo(ctx context.Context, req *todolistpb.CreateTodoRequest) (*todolistpb.CreateTodoResponse, error) {
	fmt.Println("Create Todo request")
	todo := req.GetTodo()

	return &todolistpb.CreateTodoResponse{
		Todo: &todolistpb.Todo{
			Id:      1,
			Title:   todo.GetTitle(),
			Note:    todo.GetNote(),
			DueDate: todo.GetDueDate(),
		},
	}, nil
}
