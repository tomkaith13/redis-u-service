run:
	go run ./main.go
docker-build:
	docker build -t github.com/tomkaith13/redis-u-service .

docker-run:
	docker run -p 8080:8080 github.com/tomkaith13/redis-u-service

docker-image-clean:
	docker stop github.com/tomkaith13/redis-u-service && docker rm -f github.com/tomkaith13/redis-u-service
compose-up
	docker compose up