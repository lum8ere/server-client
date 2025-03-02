import time
import threading
from clientV2.config import settings
from clientV2.adapters.communication.ws_client import WSClient
from clientV2.adapters.metrics.collector import collect_metrics
from clientV2.core.use_cases.send_metrics import send_metrics
from clientV2.core.use_cases.process_command import process_command
from clientV2.core.services.logger_service import LoggerService

logger = LoggerService()

def main():
    logger.info("Starting clientv2 application...")

    # Создаем экземпляр WebSocket клиента (он отправит регистрацию при установке соединения)
    ws_client = WSClient()
    # Устанавливаем callback для входящих сообщений (команд)
    ws_client.on_message_callback = lambda msg: process_command(msg, logger)
    ws_client.start()

    # Запускаем поток для периодической отправки метрик
    def metrics_loop():
        while True:
            metric = collect_metrics()
            # Формируем сообщение с action "sent_metrics" и payload с метриками
            send_metrics(ws_client, metric, logger)
            time.sleep(settings.METRICS_INTERVAL)
    
    metrics_thread = threading.Thread(target=metrics_loop, daemon=True)
    metrics_thread.start()

    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        logger.info("Shutting down client...")
        ws_client.stop()

if __name__ == "__main__":
    main()
