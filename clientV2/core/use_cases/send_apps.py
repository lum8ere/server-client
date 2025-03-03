import json
from clientV2.core.services.logger_service import LoggerService
from clientV2.utils.device_id import get_device_id
from clientV2.adapters.metrics.installed_apps import get_installed_applications_windows

logger = LoggerService()

def send_apps(ws_client, logger: LoggerService):
    """
    Отправляет список установленных приложений на сервер в сообщении с action='sent_apps'.
    """
    apps = get_installed_applications_windows()
    apps_dict = [app.to_dict() for app in apps]
    message = {
        "action": "sent_apps",
        "device_key": get_device_id(),
        "payload": apps_dict
    }
    try:
        ws_client.send_message(json.dumps(message))
        logger.info("Installed apps sent successfully.")
    except Exception as e:
        logger.error(f"Failed to send installed apps: {e}")
