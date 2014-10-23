SRC     := $(wildcard *.go) main/client.go
BIN     := bin
CLIENT  := ${BIN}/client
GOPATH  ?= ${PWD}

all: client

client: ${CLIENT}

${CLIENT}: ${SRC}
	@echo build $@
	@GOPATH=${GOPATH} go build -o $@ -- main/client.go

clean:
	@echo clean ${BIN}
	@rm -f ${BIN}/*

.PHONY: all client clean
