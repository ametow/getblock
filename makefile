API_KEY=<api_key>

all: cli service

cli:
	go build -o getblock-cli cmd/cli/main.go

service:
	go build -o getblockd cmd/daemon/main.go

run-cli:
	GETBLOCK_KEY=${API_KEY} ./getblock-cli -c etc/config.yaml

run-service:
	GETBLOCK_KEY=${API_KEY} ./getblockd -c etc/config.yaml