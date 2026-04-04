VERSION := $(shell cat VERSION.txt)

.PHONY: build run test clean

build:
	go build -ldflags "-s -w -X main.version=$(VERSION)" -o bin/git-contrib ./cmd/

run:
	go run ./cmd/

test:
	go test ./... -v -cover

clean:
	rm -rf bin/ dist/
