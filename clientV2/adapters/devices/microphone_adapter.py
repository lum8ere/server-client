import io
import wave
import pyaudio
import base64
import json
import threading
import time
from clientV2.core.services.logger_service import LoggerService
from clientV2.utils.device_id import get_device_id

logger = LoggerService()

# Глобальные переменные для аудио стриминга
_audio_stream = None
_audio_streaming_thread = None
_audio_streaming_active = False

# Глобальная переменная для WS клиента, которую нужно установить извне
ws_client = None

def set_ws_client(client):
    """
    Устанавливает глобальный WS-клиент,
    который будет использоваться для отправки кадров.
    """
    global ws_client
    ws_client = client
    logger.info("Microphone ws_client is set.")

def start_audio_stream():
    """
    Запускает стриминг аудио с микрофона.
    Открывает аудиопоток с параметрами (формат 16bit, 1 канал, 44100 Гц, буфер 1024 фрейма)
    и запускает поток, который в цикле читает данные и отправляет их через WS.
    """
    global _audio_stream, _audio_streaming_thread, _audio_streaming_active
    if _audio_streaming_active:
        logger.info("Audio streaming is already active.")
        return
    p = pyaudio.PyAudio()
    try:
        _audio_stream = p.open(
            format=pyaudio.paInt16,
            channels=1,
            rate=44100,
            input=True,
            frames_per_buffer=1024
        )
    except Exception as e:
        logger.error(f"Error opening microphone: {e}")
        return
    _audio_streaming_active = True
    _audio_streaming_thread = threading.Thread(target=_stream_audio, args=(p,), daemon=True)
    _audio_streaming_thread.start()
    logger.info("Audio streaming started.")

def _stream_audio(p):
    """
    Поток для непрерывного чтения аудио чанков и отправки их через WebSocket.
    Каждый прочитанный блок кодируется в base64 и отправляется в JSON-сообщении.
    """
    global _audio_stream, _audio_streaming_active
    while _audio_streaming_active:
        try:
            # Читаем один блок аудио
            data = _audio_stream.read(1024, exception_on_overflow=False)
        except Exception as e:
            logger.error(f"Error reading audio data: {e}")
            continue
        # Кодируем полученные данные в base64
        audio_b64 = base64.b64encode(data).decode('utf-8')
        message = {
            "action": "audio_stream",
            "device_key": get_device_id(),
            "payload": audio_b64
        }
        if ws_client is not None:
            try:
                ws_client.send_message(json.dumps(message))
            except Exception as e:
                logger.error(f"Error sending audio chunk: {e}")
        else:
            logger.warn("ws_client is not set. Cannot send audio chunk.")
        # Немного задерживаем цикл (можно настроить частоту отправки)
        time.sleep(0.01)
    logger.info("Exiting audio streaming loop.")
    _audio_stream.stop_stream()
    _audio_stream.close()
    p.terminate()

def stop_audio_stream():
    """
    Останавливает стриминг аудио и освобождает ресурсы.
    """
    global _audio_streaming_active, _audio_streaming_thread, _audio_stream
    if not _audio_streaming_active:
        logger.info("Audio streaming is not active.")
        return
    _audio_streaming_active = False
    if _audio_streaming_thread is not None:
        _audio_streaming_thread.join(timeout=2)
        _audio_streaming_thread = None
    logger.info("Audio streaming stopped.")

def record_audio(duration=5):
    logger.info("Recording audio...")
    p = pyaudio.PyAudio()
    try:
        stream = p.open(format=pyaudio.paInt16, channels=1, rate=44100, input=True, frames_per_buffer=1024)
    except Exception as e:
        logger.error(f"Error opening microphone: {e}")
        return None
    frames = []
    for _ in range(0, int(44100 / 1024 * duration)):
        data = stream.read(1024, exception_on_overflow=False)
        frames.append(data)
    stream.stop_stream()
    stream.close()
    p.terminate()

    buffer = io.BytesIO()
    wf = wave.open(buffer, 'wb')
    wf.setnchannels(1)
    wf.setsampwidth(p.get_sample_size(pyaudio.paInt16))
    wf.setframerate(44100)
    wf.writeframes(b''.join(frames))
    wf.close()
    return buffer.getvalue()
