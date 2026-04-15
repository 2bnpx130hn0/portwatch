BINARY  := portwatch
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
BUILD_FLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: all build test lint clean run

all: build

build:
	go build $(BUILD_FLAGS) -o bin/$(BINARY) ./cmd/portwatch

test:
	go test ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/

run: build
	./bin/$(BINARY) --config config.yaml

install: build
	cp bin/$(BINARY) /usr/local/bin/$(BINARY)
