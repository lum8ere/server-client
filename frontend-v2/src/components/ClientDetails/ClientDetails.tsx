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

const { TabPane } = Tabs;
const backendUrl = 'http://localhost:9000';

// --- Типы данных ---
interface Device {
    id: string;
    device_identifier: string;
    user_id?: string;
    description?: string;
    status: string;
    last_seen: string; // ISO string
    created_at: string;
    updated_at: string;
}

interface Metric {
    id: string;
    device_id: string;
    public_ip: string;
    hostname: string;
    os_info: string;
    disk_total: number;
    disk_used: number;
    disk_free: number;
    memory_total: number;
    memory_used: number;
    memory_available: number;
    process_count: number;
    cpu_percent: number;
    bytes_sent: number;
    bytes_recv: number;
    created_at: string;
}

interface ClientDetailsState {
    device: Device | null;
    metric: Metric | null;
}

interface Location {
    status: string;
    lat: number;
    lon: number;
}

type NotificationType = 'success' | 'info' | 'warning' | 'error';

// --- Компонент для видеостриминга ---
// При получении бинарных данных (JPEG) создается Blob и отображается в <img>
interface VideoStreamProps {
    deviceId: string;
}
const VideoStream: React.FC<VideoStreamProps> = ({ deviceId }) => {
    const [videoSrc, setVideoSrc] = useState<string>('');
    useEffect(() => {
        const ws = new WebSocket(`ws://localhost:9000/ws?role=frontend&device_id=${deviceId}`);
        ws.binaryType = 'arraybuffer';
        ws.onopen = () => {
            console.log('VideoStream WebSocket connected');
        };
        ws.onmessage = (event) => {
            const blob = new Blob([event.data], { type: 'image/jpeg' });
            const url = URL.createObjectURL(blob);
            setVideoSrc(url);
        };
        ws.onerror = (err) => {
            console.error('VideoStream WebSocket error:', err);
        };
        ws.onclose = () => {
            console.log('VideoStream WebSocket closed');
        };
        return () => {
            ws.close();
        };
    }, [deviceId]);
    return <img src={videoSrc} alt="Live Video" style={{ width: '100%' }} />;
};

// --- Основной компонент ClientDetails ---
export const ClientDetails: React.FC = () => {
    const navigate = useNavigate();
    const { id } = useParams<{ id: string }>();
    const [notificationApi, contextHolder] = notification.useNotification();

    // Состояние для данных устройства и метрик
    const [details, setDetails] = useState<ClientDetailsState>({ device: null, metric: null });
    const [position, setPosition] = useState<Location | null>(null);

    // Модальные окна для видео, скриншота и аудио
    const [webcamModalVisible, setWebcamModalVisible] = useState(false);
    const [screenshotModalVisible, setScreenshotModalVisible] = useState(false);
    const [audioModalVisible, setAudioModalVisible] = useState(false);

    // При открытии модального окна для вебкамеры мы отправляем команду на запуск видеостриминга
    useEffect(() => {
        if (id) {
            fetchDeviceDetails();
            fetchMetricDetails();
            fetchMap();
        }
    }, [id]);

    // Запрос для получения данных устройства
    const fetchDeviceDetails = async () => {
        try {
            const res = await instance.get<Device>(`/api/devices/${id}`);
            setDetails((prev) => ({ ...prev, device: res.data }));
        } catch (err) {
            message.error('Ошибка при получении данных устройства');
        }
    };

    // Запрос для получения метрик устройства
    const fetchMetricDetails = async () => {
        try {
            const res = await instance.get<Metric>(`/api/metrics/${id}`);
            setDetails((prev) => ({ ...prev, metric: res.data }));
        } catch (err) {
            message.error('Ошибка при получении метрик');
        }
    };

    // Запрос для получения координат
    const fetchMap = async () => {
        try {
            const res = await instance.get<Location>(`/api/map/${id}`);
            setPosition(res.data);
        } catch (err) {
            message.error('Ошибка при получении локации');
        }
    };

    // Функция отправки команды на сервер
    const sendCommand = async (cmd: string) => {
        try {
            await instance.post(`/command?cmd=${cmd}&id=${id}`);
            openNotificationWithIcon('success', cmd);
        } catch (err) {
            openNotificationWithIcon('error', cmd);
        }
    };

    const openNotificationWithIcon = (type: NotificationType, cmd: string) => {
        notificationApi[type]({
            message:
                type === 'success'
                    ? `Команда "${cmd}" успешно отправлена`
                    : `Ошибка при отправке команды "${cmd}"`
        });
    };

    const handleBack = () => {
        navigate(-1);
    };

    // Функция для форматирования размера (байты -> человекочитаемый формат)
    const formatBytes = (bytes: number, decimals = 2) => {
        if (!bytes) return '0 Bytes';
        const k = 1024;
        const dm = decimals < 0 ? 0 : decimals;
        const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
    };

    const mapCenter: [number, number] = position ? [position.lat, position.lon] : [51.505, -0.09];

    // Dropdown для управления USB
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

    // Обработка открытия вебкам-модального окна
    const handleOpenWebcamModal = async () => {
        await sendCommand('start_camera_stream');
        setWebcamModalVisible(true);
    };

    const handleStopWebcamModal = async () => {
        await sendCommand('stop_camera_stream');
        setWebcamModalVisible(false);
    };

    // Скриншот
    const handleScreenshot = async () => {
        await sendCommand('screenshot');
        setScreenshotModalVisible(true);
    };

    // Запись аудио
    const handleRecordAudio = async () => {
        await sendCommand('record_audio');
        setAudioModalVisible(true);
    };

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
                        <Button onClick={handleOpenWebcamModal}>View Webcam</Button>
                        <Button onClick={() => sendCommand('create_vpn')}>
                            Create VPN Connection
                        </Button>
                        <Button onClick={handleRecordAudio}>Record Audio</Button>
                        <Dropdown menu={{ items: usbItems }} placement="bottomLeft">
                            <Button>USB</Button>
                        </Dropdown>
                        <Button onClick={handleScreenshot}>Take Screenshot</Button>
                    </Space>
                </Col>
            </Row>

            {/* Системная информация */}
            <div style={{ marginBottom: 16, background: '#fff', padding: 16 }}>
                {details.device && details.metric ? (
                    <Descriptions title="System Information" bordered size="small">
                        <Descriptions.Item label="Hostname">
                            {details.metric.hostname}
                        </Descriptions.Item>
                        <Descriptions.Item label="OS">{details.metric.os_info}</Descriptions.Item>
                        <Descriptions.Item label="Disk Total">
                            {formatBytes(details.metric.disk_total)}
                        </Descriptions.Item>
                        <Descriptions.Item label="Disk Free">
                            {formatBytes(details.metric.disk_free)}
                        </Descriptions.Item>
                        <Descriptions.Item label="Memory Total">
                            {formatBytes(details.metric.memory_total)}
                        </Descriptions.Item>
                        <Descriptions.Item label="Memory Available">
                            {formatBytes(details.metric.memory_available)}
                        </Descriptions.Item>
                        <Descriptions.Item label="CPU Usage">
                            {details.metric.cpu_percent}%
                        </Descriptions.Item>
                        <Descriptions.Item label="Public IP">
                            {details.metric.public_ip}
                        </Descriptions.Item>
                        <Descriptions.Item label="Last Seen">
                            {new Date(details.device.last_seen).toLocaleString()}
                        </Descriptions.Item>
                    </Descriptions>
                ) : (
                    <p>No system information available</p>
                )}
            </div>

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
                                    Current position: {position.lat}, {position.lon}
                                </Popup>
                            </Marker>
                        )}
                    </MapContainer>
                </div>
            </div>

            {/* Webcam Modal - включает компонент VideoStream для получения ретранслированного видеопотока */}
            <Modal
                title="Webcam Streaming"
                open={webcamModalVisible}
                onCancel={handleStopWebcamModal}
                footer={<Button onClick={handleStopWebcamModal}>Stop Streaming</Button>}
                width="80%"
            >
                {id && <VideoStream deviceId={id} />}
            </Modal>

            {/* Screenshot Modal */}
            <Modal
                title="Screenshot"
                open={screenshotModalVisible}
                onCancel={() => setScreenshotModalVisible(false)}
                footer={null}
            >
                <img
                    src={`${backendUrl}/uploads/latest_screenshot.jpg?t=${Date.now()}`}
                    alt="Screenshot"
                    style={{ width: '100%', border: '1px solid #ccc' }}
                />
            </Modal>

            {/* Audio Modal */}
            <Modal
                title="Recorded Audio"
                open={audioModalVisible}
                onCancel={() => setAudioModalVisible(false)}
                footer={null}
            >
                <audio
                    src={`${backendUrl}/uploads/latest_recorded_audio.wav?t=${Date.now()}`}
                    controls
                    style={{ width: '100%' }}
                >
                    Your browser does not support the audio element.
                </audio>
            </Modal>
        </div>
    );
};

export default ClientDetails;
