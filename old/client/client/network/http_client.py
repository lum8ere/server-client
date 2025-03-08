import requests
from client.logger_config import setup_logger

logger = setup_logger()

def send_post(url, data, headers, timeout=10):
    try:
        response = requests.post(url, data=data, headers=headers, timeout=timeout)
        return response
    except Exception as e:
        logger.error("HTTP POST ошибка: " + str(e))
        return None
