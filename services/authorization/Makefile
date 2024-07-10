.PHONY: build build.image.server
build:
	go build -o build/server "./cmd/server/"
	go build -o build/migrate "./cmd/migrate/"
	go build -o build/seed "./cmd/seed/"

build.image.server:
	docker build -f ./docker/app.Dockerfile -t nexa-authorization:latest .

run: build
	@export NEXA_AUTHZ_RELEASE=1; \
	./build/server

migrate: build
	@./build/migrate

seed: build
	@./build/seed

clean:
	@go clean
	@rm -fdR build