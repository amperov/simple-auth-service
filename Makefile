SHELL := /bin/bash

DB_VERSION := 15
DB_USER := unflat
DB_PASSWORD := unflat
DB_NAME := auth-db

db-load:
	echo "Pulling Container"
	docker pull postgres:$(DB_VERSION)
db-run:
	echo "Running docker container"
	docker run --name=$(DB_NAME) \
	-e POSTGRES_USER=$(DB_USER) \
	-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
	-e POSTGRES_DB=$(DB_NAME) \
 	-p 5432:5432 -d --rm postgres:$(DB_VERSION)

db-check:
	docker stop 874ee59092a7


migrate-up:
	migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable" up
migrate-down:
	migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@localhost:5432/postgres?sslmode=disable" down
