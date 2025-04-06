.PHONY: build run test clean debug

# Default Go build flags
GOFLAGS=-v

build:
	go build $(GOFLAGS) -o locket-server ./cmd/cli

run: build
	./locket-server

debug-build:
	go build -gcflags="all=-N -l" $(GOFLAGS) -o locket-server-debug ./cmd/cli

debug: debug-build
	dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./locket-server-debug

test:
	go test ./...

clean:
	rm -f locket-server locket-server-debug