run:
	@go run cmd/main.go -mode=production

run-dev:
	@swag init -g ./cmd/main.go -o ./internal/docs
	@go run cmd/main.go -mode=development

tidy:
	@go mod tidy
