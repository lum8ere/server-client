import React, { useEffect, useState } from 'react';

interface AudioStreamProps {
    deviceId: string;
}

export const AudioStream: React.FC<AudioStreamProps> = ({ deviceId }) => {
    const [audioSrc, setAudioSrc] = useState<string>('');

    useEffect(() => {
        const ws = new WebSocket('ws://localhost:9000/ws');
        ws.onopen = () => {
            // Отправляем регистрацию, если нужно
            const registrationMessage = {
                action: 'register_frontend',
                device_key: deviceId
            };
            ws.send(JSON.stringify(registrationMessage));
        };
        ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                if (data.action === 'recorded_audio') {
                    // Формируем data URL для аудио
                    setAudioSrc(`data:audio/wav;base64,${data.payload}`);
                }
            } catch (e) {
                console.error('Error processing WS message', e);
            }
        };
        return () => ws.close();
    }, [deviceId]);

    return (
        <div>
            {audioSrc ? <audio controls src={audioSrc} /> : <p>Waiting for recorded audio...</p>}
        </div>
    );
};
