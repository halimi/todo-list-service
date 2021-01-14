package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/halimi/todo-list-service/todolistpb"
	_ "github.com/lib/pq" // Init pq library
)

const createTable = `
DROP SEQUENCE IF EXISTS todo_id;
DROP TABLE IF EXISTS todo;
CREATE SEQUENCE todo_id START 1;
CREATE TABLE todo (
	ID serial PRIMARY KEY,
	TITLE TEXT NOT NULL,
	NOTE TEXT,
	DUE_DATE TIMESTAMP WITH TIME ZONE
);
`

// Postgres sql interface
type Postgres struct {
	DB *sql.DB
}

// Close is closing the database connection
func (p *Postgres) Close() error {
	return p.DB.Close()
}

// Insert is inserting the data to the database
func (p *Postgres) Insert(todo *todolistpb.Todo) (int32, error) {
	query := `
	INSERT INTO todo (id, title, note, due_date)
	VALUES (nextval('todo_id'), $1, $2, $3)
	RETURNING id;
	`

	ts, err := ptypes.Timestamp(todo.GetDueDate())
	if err != nil {
		return -1, err
	}

	rows, err := p.DB.Query(query, todo.GetTitle(), todo.GetNote(), ts)
	if err != nil {
		return -1, err
	}

	var id int32
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return -1, err
		}
	}

	return id, nil
}

// Get is getting the data from the database
func (p *Postgres) Get(id int32) (*todolistpb.Todo, error) {
	query := `
	SELECT *
	FROM todo
	WHERE id = $1;
	`

	rows, err := p.DB.Query(query, id)
	if err != nil {
		return nil, err
	}

	var t todolistpb.Todo
	var ts time.Time
	for rows.Next() {
		if err := rows.Scan(&t.Id, &t.Title, &t.Note, &ts); err != nil {
			return nil, err
		}
	}

	t.DueDate, err = ptypes.TimestampProto(ts)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Update is updating the data in the database
func (p *Postgres) Update(todo *todolistpb.Todo) (*todolistpb.Todo, error) {
	query := `
	UPDATE todo
	SET title = $1, note = $2, due_date = $3
	WHERE id = $4
	RETURNING id, title, note, due_date;
	`

	ts, err := ptypes.Timestamp(todo.GetDueDate())
	if err != nil {
		return nil, err
	}

	rows, err := p.DB.Query(query, todo.GetTitle(), todo.GetNote(), ts, todo.GetId())
	if err != nil {
		return nil, err
	}

	var t todolistpb.Todo
	for rows.Next() {
		if err := rows.Scan(&t.Id, &t.Title, &t.Note, &ts); err != nil {
			return nil, err
		}
	}

	t.DueDate, err = ptypes.TimestampProto(ts)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// Delete is deleting the data from the database
func (p *Postgres) Delete(id int32) (int64, error) {
	query := `
	DELETE FROM todo
	WHERE id = $1;
	`

	res, err := p.DB.Exec(query, id)
	if err != nil {
		return -1, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}

	return count, nil
}

// List is listing the data
func (p *Postgres) List() ([]*todolistpb.Todo, error) {
	query := `
	SELECT *
	FROM todo
	ORDER BY id;
	`

	rows, err := p.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var todoList []*todolistpb.Todo
	for rows.Next() {
		var t todolistpb.Todo
		var ts time.Time

		if err := rows.Scan(&t.Id, &t.Title, &t.Note, &ts); err != nil {
			return nil, err
		}

		t.DueDate, err = ptypes.TimestampProto(ts)
		if err != nil {
			return nil, err
		}
		todoList = append(todoList, &t)
	}

	return todoList, nil
}

// Setup the databse
func Setup() *sql.DB {
	db, err := ConnectPostgres()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	if _, err = db.Exec(createTable); err != nil {
		log.Fatalf("Could not create the table: %v", err)
	}

	return db
}

// ConnectPostgres is connecting to a Postgres database
func ConnectPostgres() (*sql.DB, error) {
	connStr := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
