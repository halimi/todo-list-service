# Todo List Service

Simple gRPC service to can manage TODO lists. It provides API for CRUD operations.

The data is persisted in a database.

## API definition

The full definition of the `Todo` type and the services can be found in the [proto file](todolistpb/todolist.proto)

`Todo` definition:
```protobuf
message Todo {
    int32 id = 1;
    string title = 2;
    string note = 3;
    google.protobuf.Timestamp due_date = 4;
}
```

Service definitions:
```protobuf
service TodoListService {
    rpc CreateTodo(CreateTodoRequest) returns (CreateTodoResponse);
    rpc ReadTodo(ReadTodoRequest) returns (ReadTodoResponse);  // return NOT_FOUND if not found
    rpc UpdateTodo(UpdateTodoRequest) returns (UpdateTodoResponse);  // return NOT_FOUND if not found
    rpc DeleteTodo(DeleteTodoRequest) returns (DeleteTodoResponse);  // return NOT_FOUND if not found
    rpc ListTodos(ListTodosRequest) returns (stream ListTodosResponse);
}
```

## Data persistence

It can store the data in any type of repository that implements the interface.

```go
type Repository interface {
	Close() error
	Insert(*todolistpb.Todo) (int32, error)
	Get(int32) (*todolistpb.Todo, error)
	Update(*todolistpb.Todo) (*todolistpb.Todo, error)
	Delete(int32) (int64, error)
	List() ([]*todolistpb.Todo, error)
}
```

At the moment it has implementation for PostgresSQL database but it can easily extensible for other databases too.

## Run the service

The easiest way to run the application is to use `docker`.

Clone the repository:
```
git clone https://github.com/halimi/todo-list-service.git
```

Run the service with docker-compose:
```
docker-compose up
```

It brings up the todo-list-service and postgres docker containers.
The service is accessible on localhost port number 5000.

To test the service you can use the sample client implementation in the [client](client) directory.
```
cd client
go run client.go
```

## Run the tests

To run the test use the `test` `make` target.
```
make test
```

It will bring up a postgres container and run the tests. After that it will destroy the postgres container.

## Kubernetes deployment

To can deploy the service in Kubernetes the project contains Kubernetes manifest files in the [kubernetes](kubernetes) directory.

The easiest way to set up a Kubernetes cluster is to use [minikube](https://minikube.sigs.k8s.io/docs/start/)

Start the minikube cluster:
```
minikube start
```

To keep the database user name and password in secret use the [postgres-secret.yaml](kubernetes/postgres-secret.yaml) file to create a `secret`:
```
kubectl apply -f postgres-secret.yaml
```

To bring up a Postgres database service use the [postgres-deployment.yaml](kubernetes/postgres-deployment.yaml) file.
```
kubectl apply -f postgres-deployment.yaml
```

The Todo list service should reach the Postgres database. To be able to do that use the [postgres-configmap.yaml](kubernetes/postgres-configmap.yaml) file to set up the database URL.
```
kubectl apply -f postgres-configmap.yaml
```

To deploy the Todo list service use the [todolist-deployment.yaml](kubernetes/todolist-deployment.yaml) file.
```
kubectl apply -f todolist-deployment.yaml
```

To can reach the service outside from the network:
```
minikube service todolist-service
```

To test the service run the sample client:
```
go run client.go -host <IP> -port <Port>
```

## Deployment strategy

This deployment is using the `RollingUpdate` strategy. When a new version is deployed then it will slowly rolling out the new version by replacing the instances one after the other until all the instances are rolled out.

Pros and Cons of this strategy:

Pros:
 - version is slowly released across instances
 - convenient for stateful applications that can handle rebalancing of the data

Cons:
 - rollout/rollback can take time
 - supporting multiple APIs is hard
 - no control over traffic

Further information about the strategies in this [blog post](https://blog.container-solutions.com/kubernetes-deployment-strategies)
