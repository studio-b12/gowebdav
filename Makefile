SRC     := $(wildcard *.go) main/client.go
BIN     := bin
CLIENT  := ${BIN}/client
GOPATH  ?= ${PWD}

all: test client

client: ${CLIENT}

${CLIENT}: ${SRC}
	@echo build $@
	@GOPATH=${GOPATH} go build -o $@ -- main/client.go

test:
	@GOPATH=${GOPATH} go test

clean:
	@echo clean ${BIN}
	@rm -f ${BIN}/*

.PHONY: all client clean
