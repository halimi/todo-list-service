gen:
	protoc todolistpb/todolist.proto --go_out=plugins=grpc:.

docker:
	docker build -t halimi/todo-list-service .
