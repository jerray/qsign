simple-test:
	@go test

test:
	@go test -race -coverprofile=coverage.out -covermode=atomic

build:
	@go build -race

lint:
	@golint

fmt:
	@gofmt -s -w *.go

watch:
	@watchman-make -p '*.go' -t simple-test

all: test build
