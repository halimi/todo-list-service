package server_test

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/halimi/todo-list-service/db"
	"github.com/halimi/todo-list-service/server"
	"github.com/halimi/todo-list-service/todolistpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateTodo(t *testing.T) {
	s := server.Server{&db.MockDB{}}

	dd, err := ptypes.TimestampProto(time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local))
	if err != nil {
		log.Fatalf("Could not convert timestamp: %v", err)
	}

	todo := &todolistpb.Todo{
		Title:   "Create Todo test",
		Note:    "This is a test",
		DueDate: dd,
	}
	res, err := s.CreateTodo(context.Background(), &todolistpb.CreateTodoRequest{Todo: todo})
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
	gotTodo := res.GetTodo()

	wantTodo := &todolistpb.Todo{
		Id:      1,
		Title:   "Create Todo test",
		Note:    "This is a test",
		DueDate: dd,
	}

	if !reflect.DeepEqual(gotTodo, wantTodo) {
		t.Fatalf("Want: %v, Got: %v\n", wantTodo, gotTodo)
	}
}

func TestReadTodo(t *testing.T) {
	s := server.Server{&db.MockDB{}}

	res, err := s.ReadTodo(context.Background(), &todolistpb.ReadTodoRequest{TodoId: 1})
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}

	gotTodo := res.GetTodo()

	dd, err := ptypes.TimestampProto(time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local))
	if err != nil {
		log.Fatalf("Could not convert timestamp: %v", err)
	}

	wantTodo := &todolistpb.Todo{
		Id:      1,
		Title:   "Test Todo",
		Note:    "This is a test",
		DueDate: dd,
	}

	if !reflect.DeepEqual(gotTodo, wantTodo) {
		t.Fatalf("Want: %v, Got: %v\n", wantTodo, gotTodo)
	}
}

func TestUpdateTodo(t *testing.T) {
	s := server.Server{&db.MockDB{}}

	dd, err := ptypes.TimestampProto(time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local))
	if err != nil {
		log.Fatalf("Could not convert timestamp: %v", err)
	}

	wantTodo := &todolistpb.Todo{
		Id:      1,
		Title:   "Update Todo test",
		Note:    "This is a test",
		DueDate: dd,
	}

	res, err := s.UpdateTodo(context.Background(), &todolistpb.UpdateTodoRequest{Todo: wantTodo})
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}

	gotTodo := res.GetTodo()

	if !reflect.DeepEqual(gotTodo, wantTodo) {
		t.Fatalf("Want: %v, Got: %v\n", wantTodo, gotTodo)
	}
}

func TestDeleteTodo(t *testing.T) {
	s := server.Server{&db.MockDB{}}

	var todoID int32 = 1
	_, gotErr := s.DeleteTodo(context.Background(), &todolistpb.DeleteTodoRequest{TodoId: todoID})

	wantErr := status.Errorf(
		codes.NotFound,
		fmt.Sprintf("Could not found Todo with the specified ID: %v", todoID),
	)

	if !reflect.DeepEqual(gotErr, wantErr) {
		t.Fatalf("Want: %v, Got: %v\n", wantErr, gotErr)
	}
}
