package db

import (
	"context"

	"github.com/halimi/todo-list-service/todolistpb"
)

const keyRepository = "Repository"

// Repository interface
type Repository interface {
	Close() error
	Insert(*todolistpb.Todo) (int32, error)
	Get(int32) (*todolistpb.Todo, error)
	Update(*todolistpb.Todo) (*todolistpb.Todo, error)
	Delete(int32) (int64, error)
	List() ([]*todolistpb.Todo, error)
}

// SetRepository sets the repository
func SetRepository(ctx context.Context, repository Repository) context.Context {
	return context.WithValue(ctx, keyRepository, repository)
}

// getRepository returns with the repository
func getRepository(ctx context.Context) Repository {
	return ctx.Value(keyRepository).(Repository)
}

// Close is closing the database connection
func Close(ctx context.Context) error {
	return getRepository(ctx).Close()
}

// Insert is inserting the data to the database
func Insert(ctx context.Context, todo *todolistpb.Todo) (int32, error) {
	return getRepository(ctx).Insert(todo)
}

// Get is getting the data from the database
func Get(ctx context.Context, id int32) (*todolistpb.Todo, error) {
	return getRepository(ctx).Get(id)
}

// Update is updating the data in the database
func Update(ctx context.Context, todo *todolistpb.Todo) (*todolistpb.Todo, error) {
	return getRepository(ctx).Update(todo)
}

// Delete is deleting the data from the database
func Delete(ctx context.Context, id int32) (int64, error) {
	return getRepository(ctx).Delete(id)
}

// List is listing the data
func List(ctx context.Context) ([]*todolistpb.Todo, error) {
	return getRepository(ctx).List()
}
