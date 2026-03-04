# Run tests
test:
	go test ./internal/... -v

# Run tests dengan coverage
test-cover:
	go test ./internal/... -coverprofile=coverage.out
	go tool cover -html=coverage.out

dev:
	swag init -g cmd/server/main.go --output docs
	cd cmd/server && wire && cd ../..
	go run cmd/server/main.go cmd/server/wire_gen.go