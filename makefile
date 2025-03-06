.PHONY: swag
swag:
	swag init -g cmd/costco/main.go -o docs/swag --parseDependency

.PHONY: dev
dev:
	go run cmd/costco/main.go