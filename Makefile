.PHONY: test run lint build 

test:
	go test -v ./...

start_server:
	go run ./cmd/server/

lint:
	golangci-lint run

build:
	go build -o bin/app ./cmd/server/


generate_docs:
	swag init --generalInfo cmd/server/main.go --output docs