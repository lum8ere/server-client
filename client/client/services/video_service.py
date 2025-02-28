import cv2
import time
import requests
from client.config import SERVER_URL, FRAME_RATE
from client.logger_config import setup_logger

logger = setup_logger()

# Глобальная переменная для управления трансляцией видео
streaming_active = False

def capture_and_send_video():
    camera = None
    while True:
        if streaming_active:
            if camera is None:
                camera = cv2.VideoCapture(0)
                if not camera.isOpened():
                    logger.error("Камера не найдена. Попытка повторного открытия через 5 секунд.")
                    camera = None
                    time.sleep(5)
                    continue
                logger.info("Камера успешно запущена.")

            ret, frame = camera.read()
            if not ret:
                logger.error("Ошибка захвата кадра")
                time.sleep(1)
                continue

            _, buffer = cv2.imencode(".jpg", frame)
            frame_data = buffer.tobytes()

            try:
                logger.info("Отправка кадра на сервер...")
                response = requests.post(
                    f"{SERVER_URL}/",
                    data=frame_data,
                    headers={"Content-Type": "image/jpeg"}
                )
                if response.status_code == 200:
                    logger.info("Кадр успешно отправлен.")
                else:
                    logger.error(f"Ошибка отправки кадра: {response.status_code} - {response.text}")
            except Exception as e:
                logger.error("Ошибка отправки видео: " + str(e))
            time.sleep(1 / FRAME_RATE)
        else:
            if camera is not None:
                camera.release()
                camera = None
                logger.info("Камера освобождена (трансляция остановлена).")
            time.sleep(1)
