.PHONY: proto test

proto:
	cd proto && buf lint . && buf generate .

test:
	go test -v ./...

build:
	go build -v ./...