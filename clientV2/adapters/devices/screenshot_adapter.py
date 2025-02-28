import mss
from mss.tools import to_png
from clientV2.core.services.logger_service import LoggerService

logger = LoggerService()

def take_screenshot():
    logger.info("Taking screenshot...")
    try:
        with mss.mss() as sct:
            monitor = sct.monitors[1]
            img = sct.grab(monitor)
            screenshot_bytes = to_png(img.rgb, img.size)
            return screenshot_bytes
    except Exception as e:
        logger.error(f"Screenshot error: {e}")
        return None
