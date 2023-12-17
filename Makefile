.PHONY: proto test test-cover build

proto:
	cd proto && buf lint . && buf generate .

test:
	go test ./...

build:
	go build -v ./...

test-cover:
	go test -coverprofile=coverage.out ./... -coverpkg=./...

cover: test-cover
	go tool cover -html=coverage.out