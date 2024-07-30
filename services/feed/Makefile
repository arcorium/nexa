.PHONY: build build.image.server
build.local:
	go build -o build/server "./cmd/server/"
	go build -o build/seed "./cmd/seed/"

build.image.server:
	docker build -f ./docker/app.Dockerfile -t nexa-feed:latest .

run.compose: pubkey.pem
	docker compose -f ./docker-compose.dev.yml -f ./docker-compose.dev.override.yml \
	-p feed up -d

stop.compose:
	docker compose down --remove-orphans

run: build
	@export NEXA_FEED_RELEASE=1; \
	./build/server

seed: build
	@./build/seed

clean:
	@go clean
	@rm -fdR build