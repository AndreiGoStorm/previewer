BIN := "./bin/previewer"
DOCKER_IMG="previewer:develop"

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/app

run: build
	$(BIN) -config ./configs/config.yml

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.3

lint: install-lint-deps
	golangci-lint run ./...