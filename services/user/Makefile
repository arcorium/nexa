.PHONY: build
build:
	@go build -o build/server "./cmd/server/"
	@go build -o build/migrate "./cmd/migrate/"
	@go build -o build/seed_perms "./cmd/seed_perms/"

run: build
	@export NEXA_USER_RELEASE=1; \
	./build/server

migrate: build
	@./build/migrate

seed.perms: build
	@./build/seed_perms

clean:
	@go clean
	@rm -fdR build