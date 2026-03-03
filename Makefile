# Run tests
test:
	go test ./internal/... -v

# Run tests dengan coverage
test-cover:
	go test ./internal/... -coverprofile=coverage.out
	go tool cover -html=coverage.out

swagger:
	swag init -g cmd/server/main.go --output docs