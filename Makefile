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
# Запуск клиента
run-client:
	python -m clientV2.infrastructure.main

# Запуск сервера
run-server:
	cd backend-v2/apps/backend-api && go run .

# Запуск фронта
run-frontend:
	cd frontend-v2 && yarn dev