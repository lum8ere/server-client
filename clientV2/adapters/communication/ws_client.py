import threading
import time
import websocket
from clientV2.config import settings
from clientV2.core.services.logger_service import LoggerService
from clientV2.utils.device_id import get_device_id

logger = LoggerService()

class WSClient:
    def __init__(self, url=settings.SERVER_WS_URL, heartbeat_interval=settings.HEARTBEAT_INTERVAL):
        self.url = url
        self.heartbeat_interval = heartbeat_interval
        self.ws = None
        self.connected = False
        self.lock = threading.Lock()
        self.stop_event = threading.Event()
        self.thread = None
        self.on_message_callback = None  # Функция обратного вызова для входящих сообщений

    def connect(self):
        while not self.stop_event.is_set():
            try:
                device_id = get_device_id()
                logger.info(f"Connecting to {self.url} with device_id: {device_id}")
                self.ws = websocket.WebSocketApp(
                    self.url,
                    header=[f"X-Device-Identifier: {device_id}"],
                    on_open=self.on_open,
                    on_message=self.on_message,
                    on_error=self.on_error,
                    on_close=self.on_close
                )
                self.ws.run_forever()
            except Exception as e:
                logger.error(f"WebSocket connection error: {e}")
            logger.info("Connection lost. Reconnecting in 5 seconds...")
            time.sleep(5)

    def on_open(self, ws):
        with self.lock:
            self.connected = True
        logger.info(f"WebSocket connection established")

    def on_message(self, ws, message):
        logger.info(f"Message received: {message}")
        if self.on_message_callback:
            self.on_message_callback(message)

    def on_error(self, ws, error):
        logger.error(f"WebSocket error: {error}")

    def on_close(self, ws, close_status_code, close_msg):
        with self.lock:
            self.connected = False
        logger.info("WebSocket connection closed.")

    def send_message(self, message):
        with self.lock:
            if self.connected and self.ws:
                try:
                    self.ws.send(message)
                    logger.info(f"Sent message: {message}")
                except Exception as e:
                    logger.error(f"Failed to send message: {e}")
            else:
                logger.warn("Not connected. Message not sent.")

    def start(self):
        self.thread = threading.Thread(target=self.connect, daemon=True)
        self.thread.start()
        threading.Thread(target=self.heartbeat, daemon=True).start()

    def heartbeat(self):
        while not self.stop_event.is_set():
            if self.connected:
                try:
                    self.ws.send("ping")
                except Exception as e:
                    logger.error(f"Heartbeat error: {e}")
            time.sleep(self.heartbeat_interval)

    def stop(self):
        self.stop_event.set()
        if self.ws:
            self.ws.close()
