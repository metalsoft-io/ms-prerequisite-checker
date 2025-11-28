BINARY = ms-prerequisite-check
VERSION = 7.0.0
GO_LDFLAGS = -s -X main.version=$(VERSION)
OS_ARCH = darwin/amd64 linux/386 linux/amd64 linux/arm

.DEFAULT_GOAL := build

.PHONY: certs build release clean

certs:
	mkdir -p certs/certs
	(cd certs/certs && go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost)

build: certs
	CGO_ENABLED=0 go build -ldflags "$(GO_LDFLAGS)" -o $(BINARY) ./cmd/cli

release: certs
	@mkdir -p bin
	@set -e; for target in $(OS_ARCH); do \
		os=$${target%/*}; arch=$${target#*/}; \
		out=bin/$${os}_$${arch}/$(BINARY); \
		mkdir -p bin/$${os}_$${arch}; \
		echo "Building $$out"; \
		GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -ldflags "$(GO_LDFLAGS)" -o $$out ./cmd/cli; \
	done

clean:
	rm -rf $(BINARY) bin certs/certs
