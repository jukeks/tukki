.PHONY: proto test

proto:
	cd proto && buf lint . && buf generate .

test:
	go test -v ./...

build:
	go build -v ./...


test-cover:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html