.PHONY: docker-up docker-down build
docker-up: build
	docker-compose -f docker-compose.yaml up

docker-down: ## Stop docker containers and clear artefacts.
	docker-compose -f docker-compose.yaml down
	docker system prune 

test:
	@go test ./... --cover

build:
	docker build -t lordrahl/scores-backend:main .
	