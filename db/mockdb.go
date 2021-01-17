package db

import (
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/halimi/todo-list-service/todolistpb"
)

// MockDB interface
type MockDB struct {
}

// Close is closing the database connection
func (m *MockDB) Close() error {
	return nil
}

// Insert is inserting the data to the database
func (m *MockDB) Insert(todo *todolistpb.Todo) (int32, error) {
	id := todo.GetId() + 1
	return id, nil
}

// Get is getting the data from the database
func (m *MockDB) Get(id int32) (*todolistpb.Todo, error) {
	return getTestTodo(id, "Test Todo"), nil
}

// Update is updating the data in the database
func (m *MockDB) Update(todo *todolistpb.Todo) (*todolistpb.Todo, error) {
	return todo, nil
}

// Delete is deleting the data from the database
func (m *MockDB) Delete(id int32) (int64, error) {
	return 0, nil
}

// List is listing the data
func (m *MockDB) List() ([]*todolistpb.Todo, error) {
	var tl []*todolistpb.Todo
	tl = append(tl, getTestTodo(1, "Test Todo"))
	tl = append(tl, getTestTodo(2, "Test Todo"))
	return tl, nil
}

func getTestTodo(id int32, title string) *todolistpb.Todo {
	dd, err := ptypes.TimestampProto(time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local))
	if err != nil {
		log.Fatalf("Could not convert timestamp: %v", err)
	}

	return &todolistpb.Todo{
		Id:      id,
		Title:   title,
		Note:    "This is a test",
		DueDate: dd,
	}
}
