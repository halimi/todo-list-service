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
	Repo db.Repository
}

// CreateTodo request handler
func (s *Server) CreateTodo(ctx context.Context, req *todolistpb.CreateTodoRequest) (*todolistpb.CreateTodoResponse, error) {
	fmt.Println("Create Todo request")
	ctx = db.SetRepository(ctx, s.Repo)
	todo := req.GetTodo()

	id, err := db.Insert(ctx, todo)
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

// ReadTodo request handler
func (s *Server) ReadTodo(ctx context.Context, req *todolistpb.ReadTodoRequest) (*todolistpb.ReadTodoResponse, error) {
	fmt.Println("Read todo request")
	ctx = db.SetRepository(ctx, s.Repo)
	todoID := req.GetTodoId()

	if todoID == 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Could not get ID"),
		)
	}

	todo, err := db.Get(ctx, todoID)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	if todo.GetId() == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not found Todo with the specified ID: %v", todoID),
		)
	}

	return &todolistpb.ReadTodoResponse{
		Todo: todo,
	}, nil
}

// UpdateTodo request handler
func (s *Server) UpdateTodo(ctx context.Context, req *todolistpb.UpdateTodoRequest) (*todolistpb.UpdateTodoResponse, error) {
	fmt.Println("Update Todo request")
	ctx = db.SetRepository(ctx, s.Repo)
	todo := req.GetTodo()

	if todo.GetId() == 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Could not get ID"),
		)
	}

	todoNew, err := db.Update(ctx, todo)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	if todoNew.GetId() == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not found Todo with the specified ID: %v", todo.GetId()),
		)
	}

	return &todolistpb.UpdateTodoResponse{
		Todo: todoNew,
	}, nil
}

// DeleteTodo request handler
func (s *Server) DeleteTodo(ctx context.Context, req *todolistpb.DeleteTodoRequest) (*todolistpb.DeleteTodoResponse, error) {
	fmt.Println("Delete todo request")
	ctx = db.SetRepository(ctx, s.Repo)
	todoID := req.GetTodoId()

	if todoID == 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Could not get ID"),
		)
	}

	count, err := db.Delete(ctx, todoID)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	if count == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not found Todo with the specified ID: %v", todoID),
		)
	}

	return &todolistpb.DeleteTodoResponse{}, nil
}

// ListTodos request handler
func (s *Server) ListTodos(req *todolistpb.ListTodosRequest, stream todolistpb.TodoListService_ListTodosServer) error {
	fmt.Println("List todos request")
	ctx := db.SetRepository(context.Background(), s.Repo)

	todoList, err := db.List(ctx)
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	for _, todo := range todoList {
		stream.Send(&todolistpb.ListTodosResponse{Todo: todo})
	}

	return nil
}
