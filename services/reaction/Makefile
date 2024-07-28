.PHONY: build build.image run.compose stop.compose run.dev

build.local:
	go build -o build/server "./cmd/server/"
	go build -o build/migrate "./cmd/migrate/"

build.image:
	docker build -f ./docker/app.Dockerfile -t nexa-token:latest .
	docker build -f ./docker/migrate.Dockerfile -t nexa-token-migrate:latest .

run.compose: pubkey.pem
	docker compose -f ./docker-compose.dev.yml \
	-p nexa-token up -d

stop.compose:
	docker compose down --remove-orphans

run.dev: run.compose build.local
	@./build/server

run: build.local
	@export NEXA_TOKEN_RELEASE=1; \
	./build/server

migrate: build.local
	@./build/migrate

seed: build.local
	@./build/seed

clean:
	@go clean
	@rm -fdR build