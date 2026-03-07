run-server:
	go run cmd/server/main.go

run-ingestor:
	go run cmd/ingestor/main.go

seed:
	go run cmd/seed/main.go

docker-up:
	docker compose up -d

docker-down:
	docker compose down

test:
	go test ./...
