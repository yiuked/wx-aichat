.PHONY: all build run gotool install clean help

BINARY_NAME=ghapi
GO_FILE:=main.go

all: gotool build

build:
	make clean
	@if [ ! -f go.mod ];then go mod init sfapi;fi
	@go env -w GOPROXY=https://goproxy.cn,direct
	@go mod tidy
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME}  ${GO_FILE}

run:
	@go run ./

gotool:
	go fmt ./
	go vet ./

install:
	make build

clean:
	@if [ -f ${BINARY_NAME} ] ; then rm ${BINARY_NAME} ; fi

