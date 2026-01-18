ifeq ($(OS),Windows_NT)
    BIN_EXT := .exe
else
    BIN_EXT :=
endif

ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

.PHONY: dev
dev: docs
	@if command -v air >/dev/null 2>&1; then \
		if [ "$(OS)" = "Windows_NT" ]; then \
			echo "Using Windows config: .air.win.toml"; \
			air -c .air.win.toml; \
		else \
			echo "Using Linux config: .air.toml"; \
			air -c .air.toml; \
		fi \
	else \
		echo "Air is not installed. Installing..."; \
		go install github.com/air-verse/air@latest; \
		if [ "$(OS)" = "Windows_NT" ]; then \
			echo "Using Windows config: .air.win.toml"; \
			air -c .air-win.toml; \
		else \
			echo "Using Linux config: .air.toml"; \
			air -c .air.toml; \
		fi \
	fi

.PHONY: build
build:
	go build -o bin/main$(BIN_EXT) ./cmd/api/main.go

.PHONY: run
run:
	go run ./cmd/api/main.go

.PHONY: preview
preview:
	./bin/main$(BIN_EXT)

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix ./...

.PHONY: docs
docs:
	@if command -v swag >/dev/null 2>&1; then \
		echo "Generating Swagger Docs..."; \
	else \
		echo "swag is not installed. Installing..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		echo "swag installed. Add swag into path if swag is not allowed command"; \
	fi; \
	swag init -g main.go -d cmd/app,internal/controller/restapi && swag fmt

.PHONY: test
test:
	go test -v ./...
