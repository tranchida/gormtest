hello:
	@echo "Hello"

build:
	docker compose build

run:
	docker compose up --build

local-build:
	@go build -o bin/gormtest .

down:
	docker compose down

all: hello build run
