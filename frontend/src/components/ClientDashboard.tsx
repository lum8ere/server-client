import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Card, Button, Image, Row, Col, Space, message } from 'antd';
import { Metrics } from '../types';
import { fetchClientMetrics, fetchServerTime, sendCommand } from 'services/api';

const ClientDashboard: React.FC = () => {
    const { id: clientId } = useParams<{ id: string }>();
    const [serverTime, setServerTime] = useState<string>('');
    const [metrics, setMetrics] = useState<Metrics | null>(null);
    const [frameUrl, setFrameUrl] = useState<string>(`/uploads/latest_frame.jpg`);
    const [screenshotUrl, setScreenshotUrl] = useState<string>(`/uploads/latest_screenshot.jpg`);

    const updateData = async () => {
        try {
            const timeText = await fetchServerTime();
            setServerTime(new Date(timeText).toLocaleString());
            const metricsData = await fetchClientMetrics(clientId!);
            setMetrics(metricsData);
            setFrameUrl(`/uploads/latest_frame.jpg?t=${Date.now()}`);
            setScreenshotUrl(`/uploads/latest_screenshot.jpg?t=${Date.now()}`);
        } catch (err) {
            console.error('Ошибка обновления данных', err);
        }
    };

    useEffect(() => {
        const interval = setInterval(() => {
            updateData();
        }, 1000);
        return () => clearInterval(interval);
    }, [clientId]);

    const handleCommand = async (cmd: string) => {
        try {
            const result = await sendCommand(clientId!, cmd);
            message.success(`Команда отправлена: ${result}`);
        } catch (err) {
            message.error('Ошибка отправки команды');
            console.error(err);
        }
    };

    return (
        <div style={{ padding: 20 }}>
            <h1>Управление клиентом {clientId}</h1>
            <Space style={{ marginBottom: 20 }}>
                <Button type="primary" onClick={() => handleCommand('start')}>
                    Start Streaming
                </Button>
                <Button type="primary" onClick={() => handleCommand('stop')}>
                    Stop Streaming
                </Button>
                <Button onClick={() => handleCommand('screenshot')}>Take Screenshot</Button>
                <Button onClick={() => handleCommand('metrics')}>Get Metrics</Button>
                <Button onClick={() => handleCommand('download')}>Download Client</Button>
                <Button onClick={() => (window.location.href = `/map?client=${clientId}`)}>
                    Show on the Map
                </Button>
            </Space>
            <Card title="Server Time" style={{ marginBottom: 20 }}>
                {serverTime}
            </Card>
            <Row gutter={[16, 16]}>
                <Col xs={24} md={12}>
                    <Card title="Video Stream">
                        <Image
                            src={frameUrl}
                            alt="Video Stream"
                            preview={false}
                            style={{ width: '100%' }}
                        />
                    </Card>
                </Col>
                <Col xs={24} md={12}>
                    <Card title="Screenshot">
                        <Image
                            src={screenshotUrl}
                            alt="Screenshot Stream"
                            preview={false}
                            style={{ width: '100%' }}
                        />
                    </Card>
                </Col>
            </Row>
            <Card title="Метрики" style={{ marginTop: 20 }}>
                {metrics ? (
                    <div>
                        <p>
                            Дисковое пространство:{' '}
                            {metrics.disk_total
                                ? (metrics.disk_total / (1024 * 1024 * 1024)).toFixed(2)
                                : '0'}{' '}
                            GB (Total),{' '}
                            {metrics.disk_free
                                ? (metrics.disk_free / (1024 * 1024 * 1024)).toFixed(2)
                                : '0'}{' '}
                            GB (Free)
                        </p>
                        <p>
                            Оперативная память:{' '}
                            {metrics.memory_total
                                ? (metrics.memory_total / (1024 * 1024 * 1024)).toFixed(2)
                                : '0'}{' '}
                            GB (Total),{' '}
                            {metrics.memory_available
                                ? (metrics.memory_available / (1024 * 1024 * 1024)).toFixed(2)
                                : '0'}{' '}
                            GB (Available)
                        </p>
                        <p>Процессор: {metrics.processor}</p>
                        <p>ОС: {metrics.os}</p>
                    </div>
                ) : (
                    <p>Метрики отсутствуют</p>
                )}
            </Card>
        </div>
    );
};

export default ClientDashboard;
