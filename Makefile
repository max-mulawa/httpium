# Example:
#   make build
#   make clean
#   make install-service
#   make uninstall-service

.PHONY: build
build: lint
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/httpium ./cmd/server

lint:
	golangci-lint run 

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
	docker build . -t mulawam/httpium:latest

docker-push: docker-build
	docker login
	docker push mulawam/httpium:latest

docker-run: docker-build
	$(info Run httpium on port 8080)
	docker run --rm --name httpium -p 8080:8080 mulawam/httpium:latest

setup-dev: install-git-hooks
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	./scripts/provision-static-dir.sh

install-git-hooks:
	cp scripts/git-hook-pre-commit.sh .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit
