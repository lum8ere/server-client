COMPOSE_FILE=docker-compose.yml

# Команды для управления контейнерами
up:
	docker-compose -f $(COMPOSE_FILE) up -d

down:
	docker-compose -f $(COMPOSE_FILE) down

restart:
	@make down
	@make up

build:
	docker-compose -f $(COMPOSE_FILE) build

# Генерация типов через GORM
gen-type:
	cd backend-v2/apps/gen-type && go run .

# Команды для запуска сервисов
client:
	@echo "Launching the Python client..."
	cd client && python -m client.main

client-v2:
	python -m clientV2.infrastructure.main

server:
	@echo "Launching the Go server..."
	cd server && go run main.go

frontend:
	@echo "Launching the frontend..."
	cd frontend-v1 && npm start

# Цель all запускает все компоненты одновременно (в фоне)
all:
	@echo "Launching the server, client, and frontend..."
	(cd server && go run main.go) & \
	(cd client && python -m client.main) & \
	(cd frontend-v1 && npm start) & \
	wait
