.PHONY: proto test test-cover build
default: proto format build test

proto:
	cd proto && \
		buf lint . && \
		rm -fr ./gen && \
		buf generate .

test:
	go test ./...

build:
	go build -v -o ./bin/server ./cmd/server
	go build -v -o ./bin/client ./cmd/client
	go build -v -o ./bin/tukkid ./cmd/tukkid

	go build -v -o ./bin/test/csvdumper ./cmd/test/csvdumper

format:
	go fmt ./...

test-cover:
	go test -coverprofile=coverage.out ./... -coverpkg=./...

cover: test-cover
	go tool cover -html=coverage.out
