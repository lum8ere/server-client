import os

# Конфигурационные параметры приложения
SERVER_WS_URL = os.getenv("SERVER_WS_URL", "ws://127.0.0.1:9000/ws")
METRICS_INTERVAL = int(os.getenv("METRICS_INTERVAL", 60))
HEARTBEAT_INTERVAL = int(os.getenv("HEARTBEAT_INTERVAL", 30))
