test:
	go test -cover ./...

spin-up-dependencies:
	docker compose up -d

spin-up-service:
	go run ./cmd/server/main.go

start:
	make spin-up-dependencies &
	make spin-up-service

shutdown-dependencies:
	docker compose down