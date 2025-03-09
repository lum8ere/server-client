import React, { useEffect, useRef } from 'react';

interface AudioStreamMicProps {
    deviceId: string | undefined;
}

export const AudioStreamMic: React.FC<AudioStreamMicProps> = ({ deviceId }) => {
    const audioContextRef = useRef<AudioContext | null>(null);
    const nextPlayTimeRef = useRef<number>(0);
    const wsRef = useRef<WebSocket | null>(null);

    useEffect(() => {
        // Создаем аудиоконтекст
        const audioContext = new AudioContext();
        audioContextRef.current = audioContext;
        nextPlayTimeRef.current = audioContext.currentTime;

        // Открываем WebSocket соединение (URL можно брать из настроек)
        const ws = new WebSocket('ws://localhost:9000/ws');
        wsRef.current = ws;
        ws.onopen = () => {
            console.info(`WS соединение для микрофонного стриминга открыто, deviceId: ${deviceId}`);
            // Регистрируем фронтенд-клиента
            const registrationMessage = {
                action: 'register_frontend',
                device_key: deviceId
            };
            ws.send(JSON.stringify(registrationMessage));
        };

        ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                if (data.action === 'audio_stream') {
                    // Получаем base64 аудио-чанка
                    const audioB64: string = data.payload;
                    // Декодируем base64 в бинарные данные
                    const binaryStr = window.atob(audioB64);
                    const len = binaryStr.length;
                    const bytes = new Uint8Array(len);
                    for (let i = 0; i < len; i++) {
                        bytes[i] = binaryStr.charCodeAt(i);
                    }
                    const arrayBuffer = bytes.buffer;
                    // Интерпретируем данные как 16-битный PCM
                    const int16Data = new Int16Array(arrayBuffer);
                    // Преобразуем в Float32Array (нормализуем в диапазоне -1...1)
                    const float32Data = new Float32Array(int16Data.length);
                    for (let i = 0; i < int16Data.length; i++) {
                        float32Data[i] = int16Data[i] / 32768;
                    }
                    // Создаем AudioBuffer с параметрами: 1 канал, sampleRate 44100
                    const audioBuffer = audioContext.createBuffer(1, float32Data.length, 44100);
                    audioBuffer.copyToChannel(float32Data, 0, 0);

                    // Создаем источник аудио и подключаем его к выходу
                    const source = audioContext.createBufferSource();
                    source.buffer = audioBuffer;
                    source.connect(audioContext.destination);

                    // Планируем воспроизведение
                    let startTime = nextPlayTimeRef.current;
                    const chunkDuration = audioBuffer.duration;
                    // Если предыдущий старт уже в прошлом, обновляем до currentTime
                    if (startTime < audioContext.currentTime) {
                        startTime = audioContext.currentTime;
                    }
                    source.start(startTime);
                    // Обновляем время для следующего чанка
                    nextPlayTimeRef.current = startTime + chunkDuration;
                }
            } catch (e) {
                console.error('Ошибка обработки аудио сообщения', e);
            }
        };

        ws.onerror = (e) => {
            console.error('WS ошибка в аудио стриминге', e);
        };

        ws.onclose = () => {
            console.info('WS соединение для аудио стриминга закрыто');
        };

        return () => {
            ws.close();
            audioContext.close();
        };
    }, [deviceId]);

    return <div>Микрофонный стриминг активен</div>;
};
