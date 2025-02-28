import pyaudio
import wave
import io
import time
import requests
from client.config import SERVER_URL, RECORD_DURATION
from client.logger_config import setup_logger

logger = setup_logger()

# Глобальная переменная для управления трансляцией аудио
mic_streaming_active = False

def capture_and_send_audio():
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
                logger.error("Ошибка захвата или отправки аудио: " + str(e))
            time.sleep(0.1)
        else:
            if stream is not None:
                stream.stop_stream()
                stream.close()
                stream = None
                logger.info("Микрофон освобождён (трансляция остановлена).")
            time.sleep(1)

def record_audio_snippet():
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
        for _ in range(0, int(44100 / 1024 * RECORD_DURATION)):
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
