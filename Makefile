gen:
	protoc todolistpb/todolist.proto --go_out=plugins=grpc:.

docker:
	docker build -t halimi/todo-list-service .

test:
	docker run -d --name todolist-postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432  postgres
	sleep 5
	go test ./...
	docker stop todolist-postgres
	docker rm todolist-postgres
