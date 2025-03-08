import React, { useEffect, useState, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
    Button,
    Space,
    Row,
    Col,
    Descriptions,
    Tabs,
    message,
    notification,
    Dropdown,
    Modal,
    Table
} from 'antd';
import type { MenuProps } from 'antd';
import instance from 'service/api';
import { MapContainer, Marker, Popup, TileLayer } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import { CameraStream } from 'components/CameraStream/CameraStream';
import { AudioStream } from 'components/AudioStream/AudioStream';

const { TabPane } = Tabs;

const backendUrl = 'http://localhost:4000';

// ---------- Типы данных ----------
interface Metrics {
    disk_total: number;
    disk_free: number;
    memory_total: number;
    memory_available: number;
    processor: string;
    os: string;
    has_password: boolean;
    minimum_password_lenght: number;
    pc_name: string;
}

interface ProcessInfo {
    Pid: number;
    Name: string;
}

interface ServiceInfo {
    Name: string;
    DisplayName: string;
    Status: string;
}

interface AppsServices {
    processes: ProcessInfo[];
    services: ServiceInfo[];
}

interface ClientNode {
    ID: string;
    IP: string;
    Status: string;
    Metrics: Metrics;
    AppsServices: AppsServices; // Уже объект
}

interface Location {
    status: string;
    lat: number;
    lon: number;
}

interface SendCommandBody {
    device_id: string;
    command: string;
}

type NotificationType = 'success' | 'info' | 'warning' | 'error';

export const ClientDetails: React.FC = () => {
    const navigate = useNavigate();
    const [api, contextHolder] = notification.useNotification();
    const { id } = useParams();

    // Состояния
    const [node, setNode] = useState<ClientNode | null>(null);
    const [position, setPosition] = useState<Location | null>(null);

    // Модалки
    const [webcamModalVisible, setWebcamModalVisible] = useState(false);
    const [screenshotModalVisible, setScreenshotModalVisible] = useState(false);
    const [audioModalVisible, setAudioModalVisible] = useState(false);

    // Пути к файлам
    const [webcamUrl, setWebcamUrl] = useState(`${backendUrl}/uploads/latest_frame.jpg`);
    const [screenshotUrl, setScreenshotUrl] = useState(
        `${backendUrl}/uploads/latest_screenshot.jpg`
    );
    const [audioUrl, setAudioUrl] = useState(`${backendUrl}/uploads/latest_recorded_audio.wav`);

    // ID таймера (снимки вебкамеры)
    const webcamIntervalRef = useRef<number | null>(null);

    // Загрузка данных
    useEffect(() => {
        if (id) {
            fetchById();
            fetchMap();
        }
    }, [id]);

    // Получаем ClientNode с бэка
    const fetchById = async () => {
        try {
            const res = await instance.get<ClientNode>(`/api/clients/${id}`);
            setNode(res.data);
        } catch (err) {
            message.error('Ошибка при получении ноды');
        }
    };

    // Получаем координаты
    const fetchMap = async () => {
        try {
            const res = await instance.get<Location>(`/api/map/${id}`);
            setPosition(res.data);
        } catch (err) {
            message.error('Ошибка при получении локации');
        }
    };

    // Отправка команд
    const sendCommand = async (cmd: string) => {
        try {
            let body: SendCommandBody = {
                command: cmd,
                device_id: id || ''
            };

            await instance.post(`/send_command`, body);
            openNotificationWithIcon('success', cmd);
        } catch (err) {
            openNotificationWithIcon('error', cmd);
        }
    };

    // Нотификации
    const openNotificationWithIcon = (type: NotificationType, cmd: string) => {
        if (type === 'success') {
            api[type]({
                message: `The command "${cmd}" has been sent successfully`
            });
        } else if (type === 'error') {
            api[type]({
                message: `Error while sending command "${cmd}"`
            });
        } else {
            api[type]({
                message: `Command "${cmd}" status: ${type}`
            });
        }
    };

    // Кнопка назад
    const handleBack = () => {
        navigate(-1);
    };

    // Форматирование байтов
    const formatBytes = (bytes: number, decimals = 2) => {
        if (!bytes) return '0 Bytes';
        const k = 1024;
        const dm = decimals < 0 ? 0 : decimals;
        const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
    };

    // Координаты карты
    const mapCenter: [number, number] = position ? [position.lat, position.lon] : [51.505, -0.09];

    // Dropdown USB
    const usbItems: MenuProps['items'] = [
        {
            key: 'on',
            label: 'Turn ON USB',
            onClick: () => sendCommand('usb_on')
        },
        {
            key: 'off',
            label: 'Turn OFF USB',
            onClick: () => sendCommand('usb_off')
        }
    ];

    // Вебкамера
    const handleOpenWebcamModal = async () => {
        await sendCommand('start_camera');
        setWebcamModalVisible(true);

        // Каждую секунду обновляем картинку
        webcamIntervalRef.current = window.setInterval(() => {
            setWebcamUrl(`${backendUrl}/uploads/latest_frame.jpg?t=${Date.now()}`);
        }, 1000);
    };

    const handleStopWebcamModal = async () => {
        await sendCommand('stop_camera');
        setWebcamModalVisible(false);

        if (webcamIntervalRef.current !== null) {
            clearInterval(webcamIntervalRef.current);
            webcamIntervalRef.current = null;
        }
    };

    // Скриншот
    const handleScreenshot = async () => {
        await sendCommand('screenshot');
        setTimeout(() => {
            setScreenshotUrl(`${backendUrl}/uploads/latest_screenshot.jpg?t=${Date.now()}`);
            setScreenshotModalVisible(true);
        }, 1000);
    };

    // Аудио
    const handleRecordAudio = async () => {
        await sendCommand('record_audio');
        setTimeout(() => {
            setAudioUrl(`${backendUrl}/uploads/latest_recorded_audio.wav?t=${Date.now()}`);
            setAudioModalVisible(true);
        }, 2000);
    };

    // Колонки для процессов
    const processColumns = [
        { title: 'PID', dataIndex: 'pid' },
        { title: 'Name', dataIndex: 'name' }
    ];

    // Колонки для служб
    const serviceColumns = [
        { title: 'Name', dataIndex: 'name' },
        { title: 'Display Name', dataIndex: 'display_name' },
        { title: 'Status', dataIndex: 'status' }
    ];

    return (
        <div style={{ padding: 16 }}>
            {contextHolder}
            <Row justify="space-between" align="middle" style={{ marginBottom: 16 }}>
                <Col>
                    <Space>
                        <Button onClick={handleBack}>{'< Back'}</Button>
                        <span style={{ fontSize: 18, fontWeight: 'bold' }}>{id || '---'}</span>
                    </Space>
                </Col>
                <Col>
                    <Space>
                        <Button onClick={handleOpenWebcamModal}>View the webcam</Button>
                        <Button onClick={() => sendCommand('create_vpn')}>
                            Create VPN connection
                        </Button>
                        <Button onClick={handleRecordAudio}>Record Audio</Button>
                        <Dropdown menu={{ items: usbItems }} placement="bottomLeft">
                            <Button>USB</Button>
                        </Dropdown>
                        <Button onClick={handleScreenshot}>Take Screenshot</Button>
                    </Space>
                </Col>
            </Row>

            {/* Информация о системе */}
            <div style={{ marginBottom: 16, background: '#fff', padding: 16 }}>
                {node ? (
                    <Descriptions title="Information about the system" bordered size="small">
                        <Descriptions.Item label="Disk total">
                            {formatBytes(node.Metrics.disk_total)}
                        </Descriptions.Item>
                        <Descriptions.Item label="Disk free">
                            {formatBytes(node.Metrics.disk_free)}
                        </Descriptions.Item>
                        <Descriptions.Item label="OS">{node.Metrics.os}</Descriptions.Item>
                        <Descriptions.Item label="Total memory">
                            {formatBytes(node.Metrics.memory_total)}
                        </Descriptions.Item>
                        <Descriptions.Item label="Processor">
                            {node.Metrics.processor}
                        </Descriptions.Item>
                        <Descriptions.Item label="Has password">
                            {node.Metrics.has_password ? 'Yes' : 'No'}
                        </Descriptions.Item>
                    </Descriptions>
                ) : (
                    <p>Нет данных о метриках</p>
                )}
            </div>

            {/* Табы */}
            <Tabs defaultActiveKey="details">
                {/* <TabPane tab="Details" key="details">
                    <p>Здесь можно вывести дополнительную информацию о деталях.</p>
                </TabPane> */}

                {/* Software: processes */}
                <TabPane tab="Software" key="software">
                    <Table
                        dataSource={node?.AppsServices?.processes || []}
                        columns={processColumns}
                        rowKey="pid"
                        pagination={{ pageSize: 5 }}
                    />
                </TabPane>

                {/* Service: services */}
                <TabPane tab="Service" key="service">
                    <Table
                        dataSource={node?.AppsServices?.services || []}
                        columns={serviceColumns}
                        rowKey="name"
                        pagination={{ pageSize: 5 }}
                    />
                </TabPane>

                <TabPane tab="Scripts" key="scripts">
                    <p>Какие-то запросы/логи.</p>
                </TabPane>
            </Tabs>

            {/* Карта */}
            <div style={{ marginTop: 16, padding: 16, background: '#fff' }}>
                <div style={{ width: '100%', height: '400px' }}>
                    <MapContainer
                        center={mapCenter}
                        zoom={13}
                        style={{ height: '100%', width: '100%' }}
                    >
                        <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
                        {position && (
                            <Marker position={[position.lat, position.lon]}>
                                <Popup>
                                    Текущая позиция: {position.lat}, {position.lon}
                                </Popup>
                            </Marker>
                        )}
                    </MapContainer>
                </div>
            </div>

            {/* Webcam Modal */}
            <Modal
                title="Webcam streaming"
                open={webcamModalVisible}
                closable={false}
                maskClosable={false}
                footer={<Button onClick={handleStopWebcamModal}>Stop streaming</Button>}
            >
                <CameraStream id={id || ''} />
            </Modal>

            {/* Screenshot Modal */}
            <Modal
                title="Screenshot"
                open={screenshotModalVisible}
                onCancel={() => setScreenshotModalVisible(false)}
                footer={null}
            >
                <img
                    src={screenshotUrl}
                    alt="Screenshot"
                    style={{ width: '100%', border: '1px solid #ccc' }}
                />
            </Modal>

            {/* Audio Modal */}
            <Modal
                title="Recorded audio"
                open={audioModalVisible}
                onCancel={() => setAudioModalVisible(false)}
                footer={null}
            >
                <AudioStream deviceId={id || ''} />
            </Modal>
        </div>
    );
};
