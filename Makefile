.PHONY: proto test test-cover build

proto:
	cd proto && buf lint . && buf generate .

test:
	go test ./...

build:
	go build -v ./...

test-cover:
	go test -v -coverprofile=coverage.out ./... -covermode=atomic -coverpkg=./...

cover: test-cover
	go tool cover -html=coverage.out