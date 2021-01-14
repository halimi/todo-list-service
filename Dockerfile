# build stage
FROM golang:1.15 AS build
WORKDIR /go/src/github.com/halimi/todo-list-service
COPY . .
ENV CGO_ENABLED=0
RUN go build

# final stage
FROM scratch
EXPOSE 5000/tcp
COPY --from=build /go/src/github.com/halimi/todo-list-service/todo-list-service .

ENTRYPOINT ["./todo-list-service"]
