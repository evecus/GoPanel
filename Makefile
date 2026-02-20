.PHONY: all build build-web build-go docker clean dev-web install

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS  = -ldflags "-s -w -X main.version=$(VERSION)"

all: build

build-web:
	cd web && npm install && npm run build

build-go:
	go mod tidy
	CGO_ENABLED=1 go build $(LDFLAGS) -o gopanel .

build: build-web build-go

build-linux-amd64: build-web
	go mod tidy
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build $(LDFLAGS) -o gopanel-linux-amd64 .

build-linux-arm64: build-web
	go mod tidy
	GOOS=linux GOARCH=arm64 CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc go build $(LDFLAGS) -o gopanel-linux-arm64 .

docker:
	docker build --build-arg VERSION=$(VERSION) -t gopanel:$(VERSION) -t gopanel:latest .

dev-web:
	cd web && npm run dev

clean:
	rm -f gopanel gopanel-linux-*
	rm -rf web/dist

install: build
	sudo install -m 755 gopanel /usr/local/bin/gopanel
