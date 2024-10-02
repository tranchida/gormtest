hello:
	@echo "Hello"

build:
	@go build -o bin/gormtest cmd/gormtest/main.go
	podman build -t localhost/gormtest:${version} .

run:
	podman run --rm -it -p 8080:8080 localhost/gormtest:${version}

all: hello build run