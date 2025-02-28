import cv2
import threading
import time
from clientV2.core.services.logger_service import LoggerService

logger = LoggerService()

# Глобальные переменные для стриминга камеры
_camera_cap = None
_streaming_thread = None
_streaming_active = False

def start_camera_stream():
    """Запускает непрерывный стриминг с камеры в отдельном потоке."""
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
    """Поток, считывающий кадры с камеры, пока активен стриминг."""
    global _camera_cap, _streaming_active
    while _streaming_active:
        ret, frame = _camera_cap.read()
        if not ret:
            logger.error("Failed to read frame from camera.")
            break
        # Здесь можно реализовать отправку кадра через WebSocket или другое действие.
        # Для демонстрации просто ждём для соблюдения частоты кадров.
        time.sleep(1 / 10)  # например, 10 FPS
    logger.info("Exiting camera streaming loop.")

def stop_camera_stream():
    """Останавливает стриминг камеры и освобождает ресурсы."""
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
    """Считывает один кадр с камеры (на случай, если требуется одиночный снимок)."""
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
