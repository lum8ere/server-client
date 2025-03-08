import base64
import json
from clientV2.core.services.logger_service import LoggerService
from clientV2.adapters.devices import microphone_adapter
from clientV2.utils.device_id import get_device_id

logger = LoggerService()

def send_recorded_audio():
    audio_bytes = microphone_adapter.record_audio(duration=5)
    if audio_bytes is None:
        logger.error("Audio recording failed.")
        return
    audio_base64 = base64.b64encode(audio_bytes).decode('utf-8')
    message = {
        "action": "recorded_audio",
        "device_key": get_device_id(),  # либо другой идентификатор
        "payload": audio_base64
    }
    if microphone_adapter.ws_client is not None:
        try:
            microphone_adapter.ws_client.send_message(json.dumps(message))
        except Exception as e:
            logger.error(f"Error sending recorded audio: {e}")
    else:
        logger.warn("ws_client is not set. Cannot send audio.")
