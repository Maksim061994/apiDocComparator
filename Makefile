.PHONY: dev
dev:
	go build -v ./cmd/apiserver

.PHONY: prod
prod:
	env GOOS=linux GOARCH=amd64 go build -v ./cmd/apiserver

.PHONY: testing
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := dev
