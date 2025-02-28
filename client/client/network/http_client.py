import requests
from client.logger_config import setup_logger

logger = setup_logger()

def send_post(url, data, headers):
    try:
        response = requests.post(url, data=data, headers=headers)
        return response
    except Exception as e:
        logger.error("HTTP POST ошибка: " + str(e))
        return None
