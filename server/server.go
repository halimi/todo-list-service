package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/halimi/todo-list-service/db"
	"github.com/halimi/todo-list-service/todolistpb"
)

// Server is implementing TodoListServiceServer interface
type Server struct {
	Postgres *db.Postgres
}

// CreateTodo request handler
func (s *Server) CreateTodo(ctx context.Context, req *todolistpb.CreateTodoRequest) (*todolistpb.CreateTodoResponse, error) {
	fmt.Println("Create Todo request")
	todo := req.GetTodo()

	id, err := s.Postgres.Insert(todo)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	return &todolistpb.CreateTodoResponse{
		Todo: &todolistpb.Todo{
			Id:      id,
			Title:   todo.GetTitle(),
			Note:    todo.GetNote(),
			DueDate: todo.GetDueDate(),
		},
	}, nil
}
