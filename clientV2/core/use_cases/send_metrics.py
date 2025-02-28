import json
from clientV2.core.entities.metric import Metric
from clientV2.core.services.logger_service import LoggerService

def send_metrics(ws_client, metric: Metric, logger: LoggerService):
    try:
        metrics_json = json.dumps(metric.to_dict())
        ws_client.send_message(metrics_json)
        logger.info("Metrics sent successfully.")
    except Exception as e:
        logger.error(f"Failed to send metrics: {e}")
