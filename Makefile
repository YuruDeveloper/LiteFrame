.PHONY: all test clean lint

all: test lint

test:
	@echo "Testing..."
	@go test -v ./...

clean:
	@echo "Cleaning..."
	@go clean -testcache

lint:
	@echo "Linting..."
	@golangci-lint run
