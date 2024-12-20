BIN := "./bin/previewer"

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/app

run:
	$(BIN) -config ./configs/config.yml

test:
	go test -race ./internal/... -count 100

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.2

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build test lint

up-build:
	cd deployments && \
	docker-compose --project-name="previewer" up --build

up:
	cd deployments && \
	docker-compose --project-name="previewer" up -d

down:
	cd deployments && \
	docker-compose --project-name="previewer" stop

integration-tests: up
	go test -race ./integrations/...
	make down