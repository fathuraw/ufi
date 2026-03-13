VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X github.com/fathuraw/ufi/cmd.Version=$(VERSION)

.PHONY: build install clean

build:
	go build -ldflags "$(LDFLAGS)" -o ufi .

install:
	go install -ldflags "$(LDFLAGS)" .

clean:
	rm -f ufi
