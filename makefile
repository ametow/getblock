API_KEY=<api_key>

all: cli service

cli:
	@CGO_ENABLED=0 go build -o getblock-cli cmd/cli/main.go

service:
	@CGO_ENABLED=0 go build -o getblockd cmd/daemon/main.go

run-cli: cli
	GETBLOCK_KEY=${API_KEY} ./getblock-cli -c etc/config.yaml

run-service: service
	GETBLOCK_KEY=${API_KEY} ./getblockd -c etc/config.yaml

image:
	docker build -t getblock:v1 --build-arg API_KEY=${API_KEY} .

container-cli: image
	docker run getblock:v1

container-service: image
	docker run -p 8080:8080 -d getblock:v1 /app/getblockd -c /app/config.yaml