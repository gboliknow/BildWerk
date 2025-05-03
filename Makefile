# .PHONY is used to declare targets that are not actual files
.PHONY: run build test runserver

# Default target to build the application
build:
	@go build -o bin/api ./cmd/app

# Run the application, requires build to be up-to-date
run: build
	@./bin/api

# Run the server directly using `go run` and initialize Swagger documentation
runserver:
	@swag init -g cmd/app/main.go
	@go run cmd/app/main.go

# Run tests with verbose output
test:
	@go test -v ./...
