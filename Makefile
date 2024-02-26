build-all:
	docker-compose -f "docker-compose.yaml" up -d --build

.DEFAULT_GOAL := build-all