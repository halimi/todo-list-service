package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/halimi/todo-list-service/db"
	"github.com/halimi/todo-list-service/server"
	"github.com/halimi/todo-list-service/todolistpb"

	"github.com/kouhin/envflag"
)

func main() {
	// set the flags to get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	dbUser := flag.String("db-user", "postgres", "DB user name")
	dbPass := flag.String("db-pass", "postgres", "DB password")
	dbHost := flag.String("db-host", "localhost", "DB host name")
	dbPort := flag.String("db-port", "5432", "DB port number")

	envflag.Parse()

	config := &db.PostgresConfig{
		User:     *dbUser,
		Password: *dbPass,
		Host:     *dbHost,
		Port:     *dbPort,
	}

	postgres := &db.Postgres{db.Setup(config)}

	if postgres == nil {
		panic("postgres is nil")
	}

	lis, err := net.Listen("tcp", "0.0.0.0:5000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)

	todolistpb.RegisterTodoListServiceServer(s, &server.Server{postgres})
	reflection.Register(s)

	go func() {
		fmt.Println("Starting server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch

	fmt.Println("Closing Postgres connection")
	if err := postgres.Close(); err != nil {
		log.Fatalf("Error on closing the database: %v", err)
	}

	fmt.Println("Closing the listener")
	if err := lis.Close(); err != nil {
		log.Fatalf("Error on closing the listener: %v", err)
	}

	fmt.Println("Stopping the server")
	s.Stop()
}
