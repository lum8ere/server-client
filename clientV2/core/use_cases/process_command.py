from clientV2.core.services.logger_service import LoggerService
from clientV2.adapters.devices import camera_adapter, screenshot_adapter
from clientV2.core.use_cases import vpn_connection, usb_ports, send_recorded_audio

# Словарь команд и соответствующих обработчиков.
COMMAND_HANDLERS = {
    "start_camera": camera_adapter.start_camera_stream,
    "stop_camera": camera_adapter.stop_camera_stream,
    "capture_frame": camera_adapter.capture_frame,
    "record_audio": send_recorded_audio.send_recorded_audio,
    "screenshot": screenshot_adapter.screenshot,
    "create_vpn": vpn_connection.create_vpn_connection,
    "enable_usb": usb_ports.enable_usb_ports,    # включение USB-портов
    "disable_usb": usb_ports.disable_usb_ports,  # отключение USB-портов
}

def process_command(command: str, logger: LoggerService):
    cmd = command.lower().strip()
    logger.info(f"Processing command: {cmd}")
    handler = COMMAND_HANDLERS.get(cmd)
    if handler:
        result = handler()  # Вызываем обработчик; если функция что-то возвращает, можно обработать результат.
        logger.info(f"Command '{cmd}' executed successfully.")
        # Если необходимо, можно отправить результат на сервер или выполнить дополнительную логику.
    else:
        logger.info(f"Unknown command received: {command}")