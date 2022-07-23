# Example:
#   make build
#   make clean
#   make install-service
#   make uninstall-service

.PHONY: build
build: lint
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/httpium ./cmd/server

lint:
	golangci-lint run -E structcheck

test: build
	go test ./...

run:  build
	./bin/httpium

clean: 
	rm -rf bin*
	rm -rf static*

install-service: build
	sudo ./scripts/install-service.sh

uninstall-service:
	sudo ./scripts/uninstall-service.sh

docker-build:
	docker build . -t localhost/httpium:dev

docker-run: docker-build
	$(info Run httpium on port 8080)
	docker run --rm --name httpium -p 8080:8080 localhost/httpium:dev

setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	./scripts/provision-static-dir.sh