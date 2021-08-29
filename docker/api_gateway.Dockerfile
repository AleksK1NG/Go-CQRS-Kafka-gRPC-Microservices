FROM golang:1.17-alpine

WORKDIR /app

ENV CONFIG=docker

COPY .. /app

RUN go get github.com/githubnemo/CompileDaemon
RUN go mod download


ENTRYPOINT CompileDaemon --build="go build -o main api_gateway_service/cmd/main.go" --command=./main