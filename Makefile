run:
	go run ./main.go
docker-build:
	docker build -t github.com/tomkaith13/redis-u-service .

clean:
	docker compose down
up:
	make docker-build && docker compose up
