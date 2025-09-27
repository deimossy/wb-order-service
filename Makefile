NETWORK=infra-net

.PHONY: up down logs build clean test bench create-network send-orders

test:
	go test -v ./...

bench:
	go test ./internal/services/order_service -bench=. -benchmem

create-network:
	docker network create $(NETWORK) || true

up: create-network
	docker compose -f kafka/docker-compose.kafka.yaml up -d
	docker compose up -d

down:
	docker compose down
	docker compose -f kafka/docker-compose.kafka.yaml down

send-orders:
	for f in internal/services/order_service/testdata/*.json; do \
		jq -c . "$$f" | docker exec -i kafka /opt/bitnami/kafka/bin/kafka-console-producer.sh --bootstrap-server kafka:9092 --topic get_orders; \
	done

logs-app-docker:
	docker compose logs -f

logs-kafka-docker:
	docker compose -f kafka/docker-compose.kafka.yaml logs -f

logger:
	docker logs order_service --follow

build:
	docker compose build
	docker compose -f kafka/docker-compose.kafka.yaml build

clean:
	docker compose down -v
	docker network rm infra-net || true
