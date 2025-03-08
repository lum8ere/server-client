import React, { useEffect, useState } from 'react';

interface CameraStreamProps {
    id: string;
    // wsUrl: string;
}

export const CameraStream: React.FC<CameraStreamProps> = ({ id }) => {
    const [imageSrc, setImageSrc] = useState<string>('');

    debugger;

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
                if (data.action === 'camera_frame') {
                    setImageSrc(`data:image/jpeg;base64,${data.payload}`);
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
    }, [id]);

    return (
        <div>
            {imageSrc ? (
                <img src={imageSrc} alt="Camera Stream" style={{ width: '100%' }} />
            ) : (
                <p>Waiting for camera frames...</p>
            )}
        </div>
    );
};
