BIN := gowebdav
SRC := $(wildcard *.go) cmd/gowebdav/main.go

all: test cmd

cmd: ${BIN}

${BIN}: ${SRC}
	go build -o $@ ./cmd/gowebdav

test:
	go test ./...

clean:
	@rm -f ${BIN}

.PHONY: all cmd clean test
