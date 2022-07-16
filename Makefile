# Example:
#   make build
#   make clean
#   make install-service
#   make uninstall-service

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/httpium ./cmd/server

test:
	go test ./...

clean: 
	rm -rf bin*

install-service: build
	sudo ./scripts/install-service.sh

uninstall-service:
	sudo ./scripts/uninstall-service.sh

docker-build:
	docker build . -t localhost/httpium:dev

docker-run: docker-build
	$(info Run httpium on port 8080)
	docker run --rm --name httpium -p 8080:8080 localhost/httpium:dev