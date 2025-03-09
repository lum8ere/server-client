import psutil
import platform
import socket
import requests
from clientV2.core.entities.metric import Metric
from clientV2.core.services.logger_service import LoggerService

logger = LoggerService()

def get_local_ip() -> str:
    """Возвращает первый найденный IPv4 адрес, не являющийся loopback."""
    for interface, addrs in psutil.net_if_addrs().items():
        for addr in addrs:
            if addr.family == socket.AF_INET and not addr.address.startswith("127."):
                return addr.address
    return "127.0.0.1"

def collect_metrics() -> Metric:
    public_ip = ""
    try:
        response = requests.get("https://api.ipify.org?format=json", timeout=5)
        if response.ok:
            public_ip = response.json().get("ip", "")
    except Exception as e:
        logger.error(f"Error getting public IP: {e}")
        # Если запрос не удался, возвращаем локальный IP
        public_ip = get_local_ip()
    
    hostname = socket.gethostname()

    # Собираем расширенную информацию об ОС
    var_os_info = ""
    system_name = platform.system()
    if system_name == "Windows":
        # platform.win32_ver() возвращает (release, version, csd, ptype)
        win_ver = platform.win32_ver()
        var_os_info = f"Windows {win_ver[0]} (Build {win_ver[1]}, {win_ver[2] or 'no SP'})"
    else:
        # Для других ОС можно использовать platform.platform()
        var_os_info = platform.platform()

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
        os_info=var_os_info,
        disk_total=disk_total,
        disk_used=disk_used,
        disk_free=disk_free,
        memory_total=memory_total,
        memory_used=memory_used,
        memory_available=memory_available,
        process_count=process_count
    )
    return metric
