.PHONY: build run clean install test fmt vet

# Version info
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags "-X main.version=$(VERSION)"

# Development
build:
	go build $(LDFLAGS) -o .bin/hab main.go

run:
	go run main.go

test:
	go test ./...

clean:
	rm -f hab .bin/hab

install: build
	cp .bin/hab /usr/local/bin/

# Code quality
fmt:
	go fmt ./...

vet:
	go vet ./...