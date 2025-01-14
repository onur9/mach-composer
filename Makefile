build:
	go build -o mach-composer

lint:
	staticcheck ./...

format:
	go fmt ./...

test:
	go test -v ./...

coverage:
	go test -race -coverprofile=coverage.txt -covermode=atomic -coverpkg=./... ./...
	go tool cover -func=coverage.txt

docker:
	docker build -t docker.pkg.github.com/labd/mach-composer/mach:latest . --progress=plain
