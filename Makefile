.PHONY: fmt lint test check build cover snap node-deps playwright-install run start gifgrep

GIFGREP_ARGS ?=

# Allow: `make gifgrep tui skynet` (extra make goals become args).
ifneq (,$(filter $(firstword $(MAKECMDGOALS)),gifgrep run start))
  ifeq (,$(GIFGREP_ARGS))
    GIFGREP_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  endif
  $(eval $(GIFGREP_ARGS):;@:)
endif

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

node-deps:
	npm install

playwright-install:
	npx playwright install chromium

run start gifgrep:
	node scripts/run-go.mjs -- $(GIFGREP_ARGS)
