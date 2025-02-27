import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Button, Modal, Table, Descriptions, message } from 'antd';
import instance from 'service/api';

interface Metrics {
    disk_total: number;
    disk_free: number;
    memory_total: number;
    memory_available: number;
    processor: string;
    os: string;
    has_password: boolean;
    minimum_password_lenght: number;
}

interface ProcessInfo {
    pid: number;
    name: string;
}

interface ServiceInfo {
    name: string;
    display_name: string;
    status: string;
}

interface AppsServices {
    processes: ProcessInfo[];
    services: ServiceInfo[];
}

export const ClientDetails: React.FC = () => {
    const { id } = useParams();
    const [metrics, setMetrics] = useState<Metrics | null>(null);
    const [appsServices, setAppsServices] = useState<AppsServices>({
        processes: [],
        services: []
    });
    const [cameraModalVisible, setCameraModalVisible] = useState(false);
    const [screenshotUrl, setScreenshotUrl] = useState('');
    // для карты — например, lat/lon из /api/clients/:id/map, если решили так
    // ...

    useEffect(() => {
        if (id) {
            fetchMetrics();
            fetchAppsServices();
        }
    }, [id]);

    const fetchMetrics = async () => {
        try {
            const res = await instance.get<Metrics>(`/api/clients/${id}/metrics`);
            setMetrics(res.data);
        } catch (err) {
            message.error('Ошибка получения метрик');
        }
    };

    const fetchAppsServices = async () => {
        try {
            const res = await instance.get<AppsServices>(`/api/clients/${id}/apps`);
            setAppsServices(res.data);
        } catch (err) {
            message.error('Ошибка получения списка приложений/служб');
        }
    };

    const sendCommand = async (cmd: string) => {
        try {
            await instance.post(`/api/clients/${id}/command`, { cmd });
            message.success(`Команда "${cmd}" отправлена`);
        } catch (err) {
            message.error(`Ошибка отправки команды "${cmd}"`);
        }
    };

    const showCameraModal = () => {
        setCameraModalVisible(true);
        // можно запустить setInterval для обновления /uploads/latest_frame.jpg
        // или организовать WebSocket
    };

    const handleCameraModalClose = () => {
        setCameraModalVisible(false);
    };

    const handleScreenshot = async () => {
        // отправляем команду на снятие скриншота
        await sendCommand('screenshot');
        // подождём чуть-чуть и обновим URL скриншота
        setTimeout(() => {
            setScreenshotUrl(`/uploads/latest_screenshot.jpg?t=${Date.now()}`);
        }, 1000);
    };

    // Примерно так можно встраивать карту Google (если у нас есть lat/lon)
    // const lat = 55.7558;
    // const lon = 37.6173;

    return (
        <div>
            <h1>Client {id}</h1>
            <div style={{ marginBottom: 16 }}>
                <Button onClick={() => sendCommand('metrics')}>Обновить метрики</Button>
                <Button onClick={() => sendCommand('list_apps_services')}>
                    Обновить процессы/службы
                </Button>
                <Button onClick={showCameraModal}>Просмотр камеры</Button>
                <Button onClick={handleScreenshot}>Сделать скриншот</Button>
                <Button onClick={() => sendCommand('vpn_create')}>Создать VPN</Button>
                <Button onClick={() => sendCommand('record_audio')}>Записать аудио</Button>
            </div>

            <Modal
                title="Камера"
                visible={cameraModalVisible}
                footer={null}
                onCancel={handleCameraModalClose}
            >
                {/* Вариант: img + автообновление */}
                <img
                    src={`/uploads/latest_frame.jpg?t=${Date.now()}`}
                    alt="camera"
                    style={{ width: '100%' }}
                />
            </Modal>

            <h2>Скриншот</h2>
            {screenshotUrl && (
                <img src={screenshotUrl} alt="screenshot" style={{ maxWidth: 300 }} />
            )}

            <h2>Метрики</h2>
            {metrics ? (
                <Descriptions bordered size="small">
                    <Descriptions.Item label="Disk total">{metrics.disk_total}</Descriptions.Item>
                    <Descriptions.Item label="Disk free">{metrics.disk_free}</Descriptions.Item>
                    <Descriptions.Item label="Memory total">
                        {metrics.memory_total}
                    </Descriptions.Item>
                    <Descriptions.Item label="Memory available">
                        {metrics.memory_available}
                    </Descriptions.Item>
                    <Descriptions.Item label="Processor">{metrics.processor}</Descriptions.Item>
                    <Descriptions.Item label="OS">{metrics.os}</Descriptions.Item>
                    <Descriptions.Item label="Has password">
                        {String(metrics.has_password)}
                    </Descriptions.Item>
                    <Descriptions.Item label="Min password length">
                        {metrics.minimum_password_lenght}
                    </Descriptions.Item>
                </Descriptions>
            ) : (
                <p>Нет данных</p>
            )}

            <h2>Приложения</h2>
            <Table
                dataSource={appsServices.processes}
                rowKey="pid"
                columns={[
                    { title: 'PID', dataIndex: 'pid' },
                    { title: 'Name', dataIndex: 'name' }
                ]}
            />

            <h2>Службы</h2>
            <Table
                dataSource={appsServices.services}
                rowKey="name"
                columns={[
                    { title: 'Name', dataIndex: 'name' },
                    { title: 'Display Name', dataIndex: 'display_name' },
                    { title: 'Status', dataIndex: 'status' }
                ]}
            />

            <h2>Карта</h2>
            {/* Либо iFrame, либо свой компонент c react-google-maps */}
            <div style={{ width: '600px', height: '400px', backgroundColor: '#eee' }}>
                {/* placeholder для карты */}
                {/* <iframe
          style={{ width: '100%', height: '100%', border: 0 }}
          src={`https://www.google.com/maps?q=${lat},${lon}&hl=es;z=14&output=embed`}
        ></iframe> */}
            </div>
        </div>
    );
};
