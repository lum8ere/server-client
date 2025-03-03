# clientV2/adapters/communication/webrtc_client.py
import asyncio
import json
import threading
import cv2
from aiortc import RTCPeerConnection, RTCSessionDescription, VideoStreamTrack, RTCIceCandidate
from av import VideoFrame
import time
from clientV2.core.services.logger_service import LoggerService

logger = LoggerService()

# Класс для передачи видеопотока с камеры
class CameraVideoTrack(VideoStreamTrack):
    def __init__(self):
        super().__init__()  # Инициализация базового трека
        self.cap = cv2.VideoCapture(0)
        if not self.cap.isOpened():
            logger.error("Camera not available for WebRTC")
    
    async def recv(self):
        pts, time_base = await self.next_timestamp()
        ret, frame = self.cap.read()
        if not ret:
            logger.error("Failed to read frame from camera in WebRTC track")
            await asyncio.sleep(1/24)
            return await self.recv()
        # Преобразуем кадр OpenCV (BGR) в VideoFrame (RGB)
        frame = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
        video_frame = VideoFrame.from_ndarray(frame, format="rgb24")
        video_frame.pts = pts
        video_frame.time_base = time_base
        return video_frame

    def stop(self):
        if self.cap:
            self.cap.release()

# Класс для управления WebRTC соединением
class WebRTCClient:
    def __init__(self, ws_send_func):
        self.pc = RTCPeerConnection()
        self.ws_send = ws_send_func  # Функция для отправки сигналинговых сообщений через WebSocket
        self.video_track = None

        # Обработка ICE кандидатов
        @self.pc.on("icecandidate")
        async def on_icecandidate(event):
            if event.candidate:
                candidate = {
                    "candidate": event.candidate.component,
                    "sdpMid": event.candidate.sdp_mid,
                    "sdpMLineIndex": event.candidate.sdp_mline_index,
                    "candidateDescription": event.candidate.__dict__.get("candidate")
                }
                message = {
                    "action": "webrtc_ice",
                    "payload": candidate
                }
                self.ws_send(json.dumps(message))

    async def start(self):
        # Добавляем видео трек для трансляции камеры
        self.video_track = CameraVideoTrack()
        self.pc.addTrack(self.video_track)

        # Создаем SDP offer
        offer = await self.pc.createOffer()
        await self.pc.setLocalDescription(offer)

        # Отправляем offer через WS
        message = {
            "action": "webrtc_offer",
            "payload": {
                "sdp": self.pc.localDescription.sdp,
                "type": self.pc.localDescription.type
            }
        }
        self.ws_send(json.dumps(message))
        logger.info("Sent WebRTC offer.")

    async def handle_answer(self, payload):
        logger.info("Received WebRTC answer.")
        answer = RTCSessionDescription(sdp=payload["sdp"], type=payload["type"])
        await self.pc.setRemoteDescription(answer)

    async def add_ice_candidate(self, payload):
        logger.info("Received ICE candidate.")
        candidate = RTCIceCandidate(
            candidate=payload.get("candidateDescription"),
            sdpMid=payload.get("sdpMid"),
            sdpMLineIndex=payload.get("sdpMLineIndex")
        )
        await self.pc.addIceCandidate(candidate)

    async def stop(self):
        if self.video_track:
            self.video_track.stop()
        await self.pc.close()
        logger.info("WebRTC connection closed.")

# Глобальный экземпляр WebRTCClient (если требуется единственный на устройство)
_webrtc_client = None

def get_webrtc_client(ws_send_func):
    global _webrtc_client
    if _webrtc_client is None:
        _webrtc_client = WebRTCClient(ws_send_func)
    return _webrtc_client

# Для вызова асинхронных функций из синхронного кода создаем вспомогательный цикл
def run_async(coro):
    loop = asyncio.get_event_loop()
    if loop.is_running():
        asyncio.ensure_future(coro)
    else:
        loop.run_until_complete(coro)
