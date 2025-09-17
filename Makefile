NETWORK=infra-net

.PHONY: up down logs build

create-network:
	docker network create $(NETWORK) || true

up: create-network
	docker compose -f kafka/docker-compose.kafka.yaml up -d
	docker compose up -d

down:
	docker compose down
	docker compose -f kafka/docker-compose.kafka.yaml down

logs-app-docker:
	docker compose logs -f

logs-kafka-docker:
	docker compose -f kafka/docker-compose.kafka.yaml logs -f

build:
	docker compose build
	docker compose -f kafka/docker-compose.kafka.yaml build

clean:
	docker compose down -v
	docker network rm infra-net || true
