.PHONY: build-lb run-lb build-server1 build-server2 run-server1 run-server2 \
        build-testservers build-all fmt vet lint test

build-lb:
	go build -o bin/lb ./cmd/lb

run-lb:
	go run ./cmd/lb

build-server1:
	go build -o bin/backend1 ./testservers/backend1

build-server2:
	go build -o bin/backend2 ./testservers/backend2

run-server1:
	go run ./testservers/backend1

run-server2:
	go run ./testservers/backend2

build-testservers: build-server1 build-server2

build-all: build-lb build-testservers

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

lint:
	golangci-lint run