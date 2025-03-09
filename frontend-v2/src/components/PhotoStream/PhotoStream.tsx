import React, { useEffect, useState } from 'react';

interface WSMessage {
    action: string;
    device_key: string;
    payload: string;
}

interface MediaCaptureProps {
    id: string | undefined;
    mode: 'capture' | 'screenshot';
}

export const MediaCapture: React.FC<MediaCaptureProps> = ({ id, mode }) => {
    const [imageSrc, setImageSrc] = useState<string>('');

    useEffect(() => {
        const ws = new WebSocket('ws://localhost:9000/ws');

        ws.onopen = () => {
            console.info(`WS connection opened for frontend client with id: ${id}`);
            const registrationMessage = {
                action: 'register_frontend',
                device_key: id
            };
            ws.send(JSON.stringify(registrationMessage));
        };

        ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                if (mode === 'capture' && data.action === 'capture_frame') {
                    setImageSrc(`data:image/jpeg;base64,${data.payload}`);
                } else if (mode === 'screenshot' && data.action === 'screenshot') {
                    setImageSrc(`data:image/png;base64,${data.payload}`);
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

        return () => {
            ws.close();
        };
    }, []);

    return (
        <div>
            {imageSrc ? (
                <img src={imageSrc} alt={`${mode} image`} style={{ width: '100%' }} />
            ) : (
                <p>Waiting for {mode === 'capture' ? 'captured frame' : 'screenshot'}...</p>
            )}
        </div>
    );
};
