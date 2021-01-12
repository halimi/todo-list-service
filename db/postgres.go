package db

import (
	"database/sql"
	"log"

	"github.com/golang/protobuf/ptypes"
	"github.com/halimi/todo-list-service/todolistpb"
	_ "github.com/lib/pq" // Init pq library
)

const createTable = `
DROP TABLE IF EXISTS todo;
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
