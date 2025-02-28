.PHONY: client server frontend all

client:
	@echo "Launching the Python client..."
	cd client && python -m client.main

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
