import React, { useEffect, useState } from 'react';
import { Table, Button, Space, Badge, Progress, Typography } from 'antd';
import { useNavigate } from 'react-router-dom';
import instance from 'service/api';

const { Title, Link } = Typography;

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

interface ClientNode {
    ID: string;
    IP: string;
    Status: string;
    Metrics: Metrics;
}

export const ClientsList: React.FC = () => {
    const [nodes, setNodes] = useState<ClientNode[]>([]);
    const [filteredNodes, setFilteredNodes] = useState<ClientNode[]>([]);
    const navigate = useNavigate();

    useEffect(() => {
        fetchNodes();
    }, []);

    const fetchNodes = async () => {
        try {
            const res = await instance.get<ClientNode[]>('/api/clients');
            setNodes(res.data);
            setFilteredNodes(res.data);
        } catch (err) {
            console.error('Ошибка получения списка клиентов', err);
        }
    };

    // Вспомогательная функция для форматирования байтов
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
            dataIndex: 'nodeName',
            render: (_: unknown, record: ClientNode) => {
                return (
                    <Link onClick={() => navigate(`/client/${record.ID}`)}>
                        {record.Metrics.pc_name || '---'}
                    </Link>
                );
            }
        },
        {
            title: 'Status',
            dataIndex: 'status',
            render: (_: unknown, record: ClientNode) => {
                const color = record.Status === 'online' ? 'green' : 'gray';
                return (
                    <Badge color={color} text={record.Status === 'online' ? 'Online' : 'Offline'} />
                );
            }
        },
        {
            title: 'Disk space available',
            // Используем formatBytes и Progress
            render: (_: unknown, record: ClientNode) => {
                // Допустим, здесь disk_total и disk_free приходят в байтах
                const total = record.Metrics.disk_total;
                const free = record.Metrics.disk_free;
                const used = total - free;
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
                            style={{ width: 80 }} // можно добавить ширину
                        />
                        <div>
                            {formatBytes(free)} / {formatBytes(total)}
                        </div>
                    </Space>
                );
            }
        },
        {
            title: 'Operating system',
            dataIndex: 'operatingSystem',
            ellipsis: true,
            render: (_: unknown, record: ClientNode) => record.Metrics.os || '---'
        }
    ];

    const rowSelection = {
        onChange: (selectedRowKeys: React.Key[], selectedRows: ClientNode[]) => {
            console.log('selectedRowKeys:', selectedRowKeys, 'selectedRows:', selectedRows);
        }
    };

    const downloadClient = () => {
        window.open('/download/client', '_blank');
    };

    return (
        <div style={{ padding: 16 }}>
            <div style={{ marginBottom: 16 }}>
                <Button onClick={downloadClient}>Add nodes</Button>
            </div>

            <Table
                rowKey="ID"
                rowSelection={rowSelection}
                columns={columns}
                dataSource={filteredNodes}
                pagination={{
                    pageSize: 5,
                    showSizeChanger: true,
                    showTotal: (total) => `Total ${total} items`
                }}
            />
        </div>
    );
};
