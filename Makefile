run:
	@echo "Запуск системы:"
	@go run cmd/main.go 

fmt:
	@go fmt ./...

vet:
	@go vet ./...

unit-tests: vet
	@echo "Запуск unit-тестов:"
	@go test -v ./internal/services/...

clean:
	@go clean -testcache