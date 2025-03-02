import json
from clientV2.core.entities.metric import Metric
from clientV2.core.services.logger_service import LoggerService
from clientV2.utils.device_id import get_device_id

def send_metrics(ws_client, metric: Metric, logger: LoggerService):
    try:
        metrics_message = {
            "action": "sent_metrics",
            "device_key": get_device_id(), # TODO: ПОКА КАК КОСТЫЛЬ, НУЖНО СДЕЛАТЬ ТАК ЧТОБ ОБЩЕЙ МЕССАДЖ ДЛЯ ВСЕХ ACTIONS 
            "payload": metric.to_dict()
        }
        ws_client.send_message(json.dumps(metrics_message))
        logger.info("Metrics sent successfully.")
    except Exception as e:
        logger.error(f"Failed to send metrics: {e}")
