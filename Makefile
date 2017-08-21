simple-test:
	@go test

test:
	@go test -race -coverprofile=coverage.out

build:
	@go build -race

watch:
	@watchman-make -p '*.go' -t simple-test

all: test build
