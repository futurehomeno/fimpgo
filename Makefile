
build-osx:
	go build -ldflags="-s -w" -o fimpcli cli/client.go

build-arm:
	GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="-s -w" -o fimpcli cli/client.go

build-amd:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o fimpcli cmd/client.go

run :
	go run cmd/main.go -c testdata/var/config.json

.phony : clean
