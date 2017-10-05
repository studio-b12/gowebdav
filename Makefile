SRC     := $(wildcard *.go) main/client.go
BIN     := bin
CLIENT  := ${BIN}/client

all: test client

client: ${CLIENT}

${CLIENT}: ${SRC}
	@echo build $@
	go build -o $@ -- main/client.go

test:
	go test

clean:
	@echo clean ${BIN}
	@rm -f ${BIN}/*

.PHONY: all client clean
