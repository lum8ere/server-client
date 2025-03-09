import mss
import base64
import json
import time
from mss.tools import to_png
from clientV2.core.services.logger_service import LoggerService
from clientV2.utils.device_id import get_device_id

logger = LoggerService()

# Глобальный объект WebSocket-клиента, который должен быть установлен из основного кода
ws_client = None

def set_ws_client(client):
    """
    Устанавливает глобальный WS-клиент, который используется для отправки сообщений.
    """
    global ws_client
    ws_client = client
    logger.info("WS клиент для screenshot_adapter установлен.")

def screenshot():
    """
    Захватывает скриншот экрана, кодирует его в PNG, преобразует в строку base64
    и отправляет через WebSocket с действием 'screenshot_response'.
    """
    logger.info("Начало захвата скриншота экрана...")
    try:
        with mss.mss() as sct:
            # Берем первый монитор. Если нужно, можно выбрать другой.
            monitor = sct.monitors[1]
            img = sct.grab(monitor)
            png_bytes = to_png(img.rgb, img.size)
    except Exception as e:
        logger.error(f"Ошибка при захвате скриншота: {e}")
        return

    try:
        screenshot_b64 = base64.b64encode(png_bytes).decode('utf-8')
    except Exception as e:
        logger.error(f"Ошибка при кодировании скриншота в base64: {e}")
        return

    message = {
        "action": "screenshot",
        "device_key": get_device_id(),
        "payload": screenshot_b64
    }

    if ws_client:
        try:
            ws_client.send_message(json.dumps(message))
            logger.info("Скриншот успешно отправлен через WebSocket.")
        except Exception as e:
            logger.error(f"Ошибка при отправке скриншота: {e}")
    else:
        logger.warn("WS клиент не установлен. Невозможно отправить скриншот.")
