all: build run

run:
	./bin/golicense

build:
	go build -o bin/golicense main.go

release:
	go build -ldflags='-s -w' -o bin/golicense main.go

install:
	go install -ldflags='-s -w' .
