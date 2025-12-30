.PHONY: fmt lint test check build cover snap

fmt:
	gofumpt -w .

lint:
	golangci-lint run

test:
	go test ./... -cover

check:
	go vet ./...

build:
	go build ./...

cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

snap:
	node scripts/ghostty-web-snap.mjs
