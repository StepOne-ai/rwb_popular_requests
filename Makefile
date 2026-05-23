BINARY=popular-requests

.PHONY: test bench up down

test:
	go test ./... --race -coverprofile=cover.out && go tool cover -html=cover.out -o=cover.html && open cover.html

bench:
	go test -bench=. -benchmem -benchtime=3s ./internal/repository/...

loadtest:
	go run ./cmd/loadtest --rps=10000 --duration=10s

up:
	docker compose up --build

down:
	docker compose down
