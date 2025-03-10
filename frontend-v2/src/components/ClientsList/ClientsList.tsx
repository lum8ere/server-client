import React, { useEffect, useState } from 'react';
import { Table, Button, Space, Badge, Progress, Typography } from 'antd';
import { useNavigate } from 'react-router-dom';
import instance from 'service/api';

const { Title, Link } = Typography;

interface Device {
    id: string;
    device_identifier: string;
    user_id?: string;
    description?: string;
    status: string;
    last_seen: string; // в формате ISO string
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

// Для отображения в таблице объединяем данные устройства и его метрик.
interface ClientNode {
    device: Device;
    metric?: Metric;
}

export const ClientsList: React.FC = () => {
    const [nodes, setNodes] = useState<ClientNode[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const navigate = useNavigate();

    useEffect(() => {
        fetchDevicesAndMetrics();
    }, []);

    const fetchDevicesAndMetrics = async () => {
        setLoading(true);
        try {
            // Выполняем параллельные запросы:
            const [devicesRes, metricsRes] = await Promise.all([
                instance.get<Device[]>('/api/devices'),
                instance.get<Metric[]>('/api/metrics')
            ]);

            // Создаем словарь метрик по device_id
            const metricsMap: { [key: string]: Metric } = {};
            metricsRes.data.forEach((m) => {
                metricsMap[m.device_id] = m;
            });

            // Объединяем устройства с метриками
            const combined: ClientNode[] = devicesRes.data.map((device) => ({
                device,
                metric: metricsMap[device.id]
            }));

            setNodes(combined);
        } catch (err) {
            console.error('Ошибка получения данных', err);
        } finally {
            setLoading(false);
        }
    };

    // Функция для форматирования байтов
    const formatBytes = (bytes: number, decimals = 2) => {
        if (!bytes) return '0 Bytes';
        const k = 1024;
        const dm = decimals < 0 ? 0 : decimals;
        const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
    };

    // Определяем колонки таблицы
    const columns = [
        {
            title: 'Node',
            render: (_: unknown, record: ClientNode) => (
                <Link onClick={() => navigate(`/devices/${record.device.id}`)}>
                    {record.metric?.hostname || record.device.device_identifier || '---'}
                </Link>
            )
        },
        {
            title: 'Status',
            dataIndex: 'status',
            render: (_: unknown, record: ClientNode) => {
                const status = record.device.status.toLowerCase();
                const color = status === 'online' ? 'green' : 'gray';
                return <Badge color={color} text={status === 'online' ? 'Online' : 'Offline'} />;
            }
        },
        {
            title: 'Disk Space Available',
            render: (_: unknown, record: ClientNode) => {
                if (!record.metric) return '---';
                const total = record.metric.disk_total;
                const free = record.metric.disk_free;
                const percentFree = (free / total) * 100;
                return (
                    <Space>
                        <Progress
                            percent={Math.round(percentFree)}
                            size="default"
                            showInfo
                            strokeColor={
                                percentFree < 15 ? 'red' : percentFree < 30 ? 'orange' : 'green'
                            }
                            style={{ width: 80 }}
                        />
                        <div>
                            {formatBytes(free)} / {formatBytes(total)}
                        </div>
                    </Space>
                );
            }
        },
        {
            title: 'Memory Available',
            render: (_: unknown, record: ClientNode) => {
                if (!record.metric) return '---';
                const total = record.metric.memory_total;
                const available = record.metric.memory_available;
                const percentAvailable = (available / total) * 100;
                return (
                    <Space>
                        <Progress
                            percent={Math.round(percentAvailable)}
                            size="default"
                            showInfo
                            strokeColor={
                                percentAvailable > 80
                                    ? 'red'
                                    : percentAvailable < 40
                                      ? 'orange'
                                      : percentAvailable < 15
                                        ? 'green'
                                        : 'green'
                            }
                            style={{ width: 80 }}
                        />
                        <div>
                            {formatBytes(available)} / {formatBytes(total)}
                        </div>
                    </Space>
                );
            }
        },
        {
            title: 'Operating System',
            render: (_: unknown, record: ClientNode) => record.metric?.os_info || '---'
        }
        // При необходимости можно добавить другие колонки (CPU, network и т.д.)
    ];

    return (
        <div style={{ padding: 16 }}>
            <div style={{ marginBottom: 16 }}>
                <Button onClick={fetchDevicesAndMetrics}>Refresh Data</Button>
            </div>
            <Table
                rowKey={(record: ClientNode) => record.device.id}
                columns={columns}
                dataSource={nodes}
                loading={loading}
                pagination={{
                    pageSize: 5,
                    showSizeChanger: true,
                    showTotal: (total: number) => `Total ${total} items`
                }}
            />
        </div>
    );
};
