.PHONY: build run clean test

BINARY=bin/todo

build:
	go build -o $(BINARY) ./cmd/todo

run: build
	./$(BINARY)

clean:
	rm -rf bin/
	rm -f todos.csv

test:
	go test -v ./...

tidy:
	go mod tidy

install: build
	cp $(BINARY) ~/bin/todo
