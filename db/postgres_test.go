package db_test

import (
	"database/sql"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/halimi/todo-list-service/db"
	"github.com/halimi/todo-list-service/todolistpb"
)

func setupDB() *sql.DB {
	config := &db.PostgresConfig{
		User:     "postgres",
		Password: "postgres",
		Host:     "localhost",
		Port:     "5432",
	}
	return db.Setup(config)
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

func TestInsert(t *testing.T) {
	postgres := &db.Postgres{setupDB()}
	defer postgres.Close()

	got, err := postgres.Insert(getTestTodo(0, "Test Todo"))
	if err != nil {
		t.Fatal(err)
	}

	var want int32 = 1

	if got != want {
		t.Fatalf("Want: %v, Got: %v\n", want, got)
	}
}

func TestGet(t *testing.T) {
	postgres := &db.Postgres{setupDB()}
	defer postgres.Close()

	id, err := postgres.Insert(getTestTodo(0, "Test Todo"))
	if err != nil {
		t.Fatal(err)
	}

	got, err := postgres.Get(id)
	if err != nil {
		t.Fatal(err)
	}

	want := getTestTodo(1, "Test Todo")

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Want: %v, Got: %v\n", want, got)
	}
}

func TestUpdate(t *testing.T) {
	postgres := &db.Postgres{setupDB()}
	defer postgres.Close()

	id, err := postgres.Insert(getTestTodo(0, "Test Todo"))
	if err != nil {
		t.Fatal(err)
	}

	want := getTestTodo(id, "Update Test Todo")

	got, err := postgres.Update(want)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Want: %v, Got: %v\n", want, got)
	}
}

func TestDelete(t *testing.T) {
	postgres := &db.Postgres{setupDB()}
	defer postgres.Close()

	id, err := postgres.Insert(getTestTodo(0, "Test Todo"))
	if err != nil {
		t.Fatal(err)
	}

	got, err := postgres.Delete(id)
	if err != nil {
		t.Fatal(err)
	}

	var want int64 = 1
	if got != want {
		t.Fatalf("Want: %v, Got: %v\n", want, got)
	}
}

func TestList(t *testing.T) {
	postgres := &db.Postgres{setupDB()}
	defer postgres.Close()

	id1, err := postgres.Insert(getTestTodo(0, "Test Todo"))
	if err != nil {
		t.Fatal(err)
	}

	id2, err := postgres.Insert(getTestTodo(0, "Test Todo"))
	if err != nil {
		t.Fatal(err)
	}

	got, err := postgres.List()
	if err != nil {
		t.Fatal(err)
	}

	var want []*todolistpb.Todo
	want = append(want, getTestTodo(id1, "Test Todo"))
	want = append(want, getTestTodo(id2, "Test Todo"))

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Want: %v, Got: %v\n", want, got)
	}
}
