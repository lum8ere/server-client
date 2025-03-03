from clientV2.core.services.logger_service import LoggerService
from clientV2.adapters.devices import camera_adapter, microphone_adapter, screenshot_adapter
from clientV2.core.use_cases.vpn_connection import create_vpn_connection
from clientV2.core.use_cases.usb_ports import enable_usb_ports, disable_usb_ports

# Словарь команд и соответствующих обработчиков.
COMMAND_HANDLERS = {
    "start_camera": camera_adapter.start_camera_stream,
    "stop_camera": camera_adapter.stop_camera_stream,
    "capture_frame": camera_adapter.capture_frame,
    "record_audio": lambda: microphone_adapter.record_audio(duration=5),
    "screenshot": screenshot_adapter.take_screenshot,
    "create_vpn": create_vpn_connection,
    "enable_usb": enable_usb_ports,    # включение USB-портов
    "disable_usb": disable_usb_ports,  # отключение USB-портов
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