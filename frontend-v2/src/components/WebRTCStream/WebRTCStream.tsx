import React, { useEffect, useRef } from 'react';

interface WebRTCStreamProps {
    deviceId: string;
    wsUrl: string; // URL сигнального сервера, например, "ws://localhost:9000/ws"
}

export const WebRTCStream: React.FC<WebRTCStreamProps> = ({ deviceId, wsUrl }) => {
    const videoRef = useRef<HTMLVideoElement>(null);
    const pcRef = useRef<RTCPeerConnection | null>(null);
    const wsRef = useRef<WebSocket | null>(null);

    useEffect(() => {
        // Создаем WS-подключение для сигналинга, можно передать query-параметры (например, role и deviceId)
        const ws = new WebSocket(`${wsUrl}?role=frontend&device_id=${deviceId}`);
        wsRef.current = ws;

        // Конфигурация для RTCPeerConnection (STUN сервер для базового соединения)
        const configuration = {
            iceServers: [{ urls: 'stun:stun.l.google.com:19302' }]
        };

        const pc = new RTCPeerConnection(configuration);
        pcRef.current = pc;

        // Отправка ICE кандидатов на сервер
        pc.onicecandidate = (event) => {
            if (event.candidate && wsRef.current) {
                const message = {
                    action: 'webrtc_ice',
                    device_key: deviceId,
                    payload: {
                        candidate: event.candidate.candidate,
                        sdpMid: event.candidate.sdpMid,
                        sdpMLineIndex: event.candidate.sdpMLineIndex
                    }
                };
                wsRef.current.send(JSON.stringify(message));
            }
        };

        // Обработка входящего медиапотока
        pc.ontrack = (event) => {
            if (videoRef.current) {
                videoRef.current.srcObject = event.streams[0];
            }
        };

        // Обработка сигналинговых сообщений по WS
        ws.onmessage = async (event) => {
            try {
                const data = JSON.parse(event.data);
                if (data.action === 'webrtc_offer') {
                    // Получили offer, устанавливаем remoteDescription и создаем answer
                    const offer = data.payload;
                    await pc.setRemoteDescription(new RTCSessionDescription(offer));
                    const answer = await pc.createAnswer();
                    await pc.setLocalDescription(answer);

                    const answerMessage = {
                        action: 'webrtc_answer',
                        device_key: deviceId,
                        payload: {
                            sdp: answer.sdp,
                            type: answer.type
                        }
                    };
                    ws.send(JSON.stringify(answerMessage));
                } else if (data.action === 'webrtc_ice') {
                    // Добавляем ICE кандидат, если получен от отправителя
                    const candidate = data.payload;
                    try {
                        await pc.addIceCandidate(new RTCIceCandidate(candidate));
                    } catch (err) {
                        console.error('Error adding received ICE candidate', err);
                    }
                }
            } catch (e) {
                console.error('Error processing WS message', e);
            }
        };

        ws.onerror = (e) => {
            console.error('WS error:', e);
        };

        ws.onclose = () => {
            console.log('WS connection closed');
        };

        // Чистим ресурсы при размонтировании компонента
        return () => {
            if (wsRef.current) {
                wsRef.current.close();
            }
            if (pcRef.current) {
                pcRef.current.close();
            }
        };
    }, [deviceId, wsUrl]);

    return <video ref={videoRef} autoPlay playsInline style={{ width: '100%' }} />;
};
