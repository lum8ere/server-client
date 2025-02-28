import pyaudio
import wave
import io
from clientV2.core.services.logger_service import LoggerService

logger = LoggerService()

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
