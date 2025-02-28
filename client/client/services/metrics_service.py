import json
import socket
import requests
import psutil
import platform
from client.system.sys_info import (
    current_user_has_password,
    get_min_password_length,
    get_running_processes,
    get_running_services,
)
from client.logger_config import setup_logger

logger = setup_logger()

def get_init_info() -> dict:
    ip = ""
    try:
        response = requests.get("https://api.ipify.org?format=json", timeout=10)
        ip = response.json().get("ip")
    except Exception as e:
        logger.error("Ошибка получения публичного IP: " + str(e))
    metrics = {}
    try:
        disk_usage = psutil.disk_usage("C:\\")
        mem = psutil.virtual_memory()
        metrics = {
            "disk_total": disk_usage.total,
            "disk_free": disk_usage.free,
            "memory_total": mem.total,
            "memory_available": mem.available,
            "processor": platform.processor(),
            "os": f"{platform.system()} {platform.release()}",
            "has_password": current_user_has_password(),
            "minimum_password_length": get_min_password_length(),
            "pc_name": socket.gethostname()
        }
    except Exception as e:
        logger.error("Ошибка сбора метрик: " + str(e))
        metrics["error"] = str(e)
    return {
        "ip": ip,
        "metrics": metrics,
        "apps_services": {
            "processes": get_running_processes(),
            "services": get_running_services()
        }
    }

def collect_and_send_metrics(ws):
    metrics = {}
    try:
        disk_usage = psutil.disk_usage("C:\\")
        mem = psutil.virtual_memory()
        metrics = {
            "disk_total": disk_usage.total,
            "disk_free": disk_usage.free,
            "memory_total": mem.total,
            "memory_available": mem.available,
            "processor": platform.processor(),
            "os": f"{platform.system()} {platform.release()}",
            "has_password": current_user_has_password(),
            "minimum_password_length": get_min_password_length()
        }
    except Exception as e:
        logger.error("Ошибка сбора метрик: " + str(e))
        metrics["error"] = str(e)
    try:
        ws.send(json.dumps({"command": "metrics", "data": metrics}))
        logger.info("Метрики отправлены на сервер")
    except Exception as e:
        logger.error("Ошибка отправки метрик: " + str(e))
