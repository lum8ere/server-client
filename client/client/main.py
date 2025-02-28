import time
import requests
import threading
from client.config import BASE_URL
from client.logger_config import setup_logger
from client.network.ws_client import control_ws
from client.services.video_service import capture_and_send_video

logger = setup_logger()

def main():
    logger.info("Клиент запустился")
    # Проверка доступности сервера
    try:
        response = requests.get(f"{BASE_URL}/api/time", timeout=5)
        if response.status_code == 200:
            logger.info(f"Сервер доступен. Серверное время: {response.text}")
        else:
            logger.error(f"Ошибка подключения к серверу: {response.status_code}")
            exit(1)
    except Exception as e:
        logger.error(f"Ошибка подключения к серверу: {e}")
        exit(1)

    # Запуск WebSocket-клиента и видеотрансляции в отдельных потоках
    ws_thread = threading.Thread(target=control_ws, daemon=True)
    ws_thread.start()

    video_thread = threading.Thread(target=capture_and_send_video, daemon=True)
    video_thread.start()

    # Здесь можно добавить запуск потоков для аудио, метрик и т.д.
    ws_thread.join()
    video_thread.join()

if __name__ == "__main__":
    main()
