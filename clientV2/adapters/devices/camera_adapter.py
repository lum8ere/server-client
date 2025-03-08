import cv2
import threading
import time
import base64
import json
from clientV2.core.services.logger_service import LoggerService
from clientV2.utils.device_id import get_device_id

logger = LoggerService()

# Глобальные переменные для стриминга камеры
_camera_cap = None
_streaming_thread = None
_streaming_active = False

# Глобальная переменная для WS клиента, которую нужно установить извне
ws_client = None

def set_ws_client(client):
    """
    Устанавливает глобальный WS-клиент,
    который будет использоваться для отправки кадров.
    """
    global ws_client
    ws_client = client

def start_camera_stream():
    """
    Запускает непрерывный стриминг с вебкамеры в отдельном потоке.
    Каждый кадр кодируется в JPEG, затем преобразуется в base64
    и отправляется на бэкенд через WebSocket.
    """
    global _camera_cap, _streaming_thread, _streaming_active
    if _streaming_active:
        logger.info("Camera streaming is already active.")
        return
    _camera_cap = cv2.VideoCapture(0)
    if not _camera_cap.isOpened():
        logger.error("Camera not available.")
        return
    _streaming_active = True
    _streaming_thread = threading.Thread(target=_stream_camera, daemon=True)
    _streaming_thread.start()
    logger.info("Camera streaming started.")

def _stream_camera():
    """
    Поток, который захватывает кадры с камеры, кодирует их и отправляет через WebSocket.
    """
    global _camera_cap, _streaming_active
    while _streaming_active:
        ret, frame = _camera_cap.read()
        if not ret:
            logger.error("Failed to read frame from camera.")
            continue
        ret, buffer = cv2.imencode('.jpg', frame)
        if not ret:
            logger.error("Failed to encode frame.")
            continue
        # Преобразуем JPEG-буфер в строку base64
        jpg_as_text = base64.b64encode(buffer).decode('utf-8')
        # Формируем сообщение для передачи
        message = {
            "action": "camera_frame",
            "device_key": get_device_id(),
            "payload": jpg_as_text
        }
        if ws_client is not None:
            try:
                ws_client.send_message(json.dumps(message))
            except Exception as e:
                logger.error(f"Error sending camera frame: {e}")
        else:
            logger.warn("ws_client is not set. Cannot send camera frame.")
        # Ограничиваем частоту кадров, например, до 15 FPS
        time.sleep(1/24)
    logger.info("Exiting camera streaming loop.")

def stop_camera_stream():
    """
    Останавливает стриминг камеры и освобождает ресурсы.
    """
    global _camera_cap, _streaming_active, _streaming_thread
    if not _streaming_active:
        logger.info("Camera streaming is not active.")
        return
    _streaming_active = False
    if _camera_cap is not None:
        _camera_cap.release()
        _camera_cap = None
    if _streaming_thread is not None:
        _streaming_thread.join(timeout=2)
        _streaming_thread = None
    logger.info("Camera streaming stopped.")

def capture_frame():
    """
    Захватывает один кадр с вебкамеры, кодирует его в JPEG и возвращает бинарные данные.
    """
    logger.info("Capturing a single frame from the camera...")
    cap = cv2.VideoCapture(0)
    ret, frame = cap.read()
    if not ret:
        logger.error("Failed to capture frame from camera.")
        cap.release()
        return None
    ret, jpeg = cv2.imencode(".jpg", frame)
    cap.release()
    if not ret:
        logger.error("Failed to encode frame.")
        return None
    return jpeg.tobytes()
