import psutil
import platform
import socket
import requests
from clientV2.core.entities.metric import Metric
from clientV2.core.services.logger_service import LoggerService

logger = LoggerService()

def collect_metrics() -> Metric:
    public_ip = ""
    try:
        response = requests.get("https://api.ipify.org?format=json", timeout=5)
        if response.ok:
            public_ip = response.json().get("ip", "")
    except Exception as e:
        logger.error(f"Error getting public IP: {e}")
    
    hostname = socket.gethostname()
    os_info = f"{platform.system()} {platform.release()}"
    
    disk_total = disk_used = disk_free = 0
    try:
        disk = psutil.disk_usage("/")
        disk_total = disk.total
        disk_used = disk.used
        disk_free = disk.free
    except Exception as e:
        logger.error(f"Disk info error: {e}")
    
    memory_total = memory_used = memory_available = 0
    try:
        mem = psutil.virtual_memory()
        memory_total = mem.total
        memory_used = mem.used
        memory_available = mem.available
    except Exception as e:
        logger.error(f"Memory info error: {e}")
    
    process_count = 0
    try:
        process_count = len(psutil.pids())
    except Exception as e:
        logger.error(f"Process count error: {e}")

    metric = Metric(
        public_ip=public_ip,
        hostname=hostname,
        os_info=os_info,
        disk_total=disk_total,
        disk_used=disk_used,
        disk_free=disk_free,
        memory_total=memory_total,
        memory_used=memory_used,
        memory_available=memory_available,
        process_count=process_count
    )
    return metric
