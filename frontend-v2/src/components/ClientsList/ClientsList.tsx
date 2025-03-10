import React, { useEffect, useState } from 'react';
import {
    Table,
    Button,
    Badge,
    Space,
    Input,
    Select,
    Modal,
    Form,
    message,
    Progress,
    ConfigProvider
} from 'antd';
import { useNavigate } from 'react-router-dom';
import instance from 'service/api';

interface Device {
    id: string;
    device_identifier: string;
    description?: string;
    status: string;
    last_seen: string;
    created_at: string;
    updated_at: string;
    group_id?: string;
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

interface ClientNode {
    device: Device;
    metric?: Metric;
}

interface DeviceGroup {
    id: string;
    name: string;
    description?: string;
    created_at: string;
}

const formatBytes = (bytes: number, decimals = 2): string => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
};

// Цветовая схема для дискового пространства (прямой процент свободного места)
const getProgressColor = (percent: number): string => {
    if (percent >= 60) {
        return '#52c41a'; // зеленый (остается как есть)
    } else if (percent >= 30) {
        return '#faad14'; // желтый
    } else {
        return '#f5222d'; // красный
    }
};

// Цветовая схема для памяти (чем выше использование, тем хуже)
const getMemoryProgressColor = (usedPercent: number): string => {
    if (usedPercent < 30) {
        return '#52c41a'; // green - low usage
    } else if (usedPercent < 70) {
        return '#faad14'; // yellow - moderate usage
    } else {
        return '#f5222d'; // red - high usage
    }
};

export const ClientsList: React.FC = () => {
    const [nodes, setNodes] = useState<ClientNode[]>([]);
    const [groups, setGroups] = useState<DeviceGroup[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
    const [assignModalVisible, setAssignModalVisible] = useState<boolean>(false);
    const [form] = Form.useForm();
    const [groupFilter, setGroupFilter] = useState<string>('all');
    const [searchText, setSearchText] = useState<string>('');
    const navigate = useNavigate();

    const fetchDevices = async () => {
        setLoading(true);
        try {
            const [devicesRes, metricsRes] = await Promise.all([
                instance.get<Device[]>('/api/devices'),
                instance.get<Metric[]>('/api/metrics')
            ]);

            const metricsMap: { [key: string]: Metric } = {};
            metricsRes.data.forEach((m) => {
                metricsMap[m.device_id] = m;
            });

            const combined: ClientNode[] = devicesRes.data.map((device) => ({
                device,
                metric: metricsMap[device.id]
            }));

            setNodes(combined);
        } catch (err) {
            message.error('Failed to fetch devices');
        } finally {
            setLoading(false);
        }
    };

    const fetchGroups = async () => {
        try {
            const res = await instance.get<DeviceGroup[]>('/api/device-groups');
            setGroups(res.data);
        } catch (err) {
            message.error('Failed to fetch device groups');
        }
    };

    useEffect(() => {
        fetchDevices();
        fetchGroups();
    }, []);

    // Фильтрация по группе и поиску по названию (hostname или device_identifier)
    const filteredNodes = nodes.filter((node) => {
        const matchesGroup = groupFilter === 'all' ? true : node.device.group_id === groupFilter;
        const deviceName = node.metric?.hostname || node.device.device_identifier || '';
        const matchesSearch = deviceName.toLowerCase().includes(searchText.toLowerCase());
        return matchesGroup && matchesSearch;
    });

    const handleRowSelection = {
        selectedRowKeys,
        onChange: (selectedKeys: React.Key[]) => {
            setSelectedRowKeys(selectedKeys);
        }
    };

    const handleAssign = async () => {
        if (selectedRowKeys.length !== 1) {
            message.warning('Please select exactly one device to assign.');
            return;
        }
        setAssignModalVisible(true);
    };

    const onAssignFinish = async (values: any) => {
        const payload = {
            device_id: selectedRowKeys[0],
            group_id: values.group_id
        };
        try {
            await instance.post('/api/device-groups/assign', payload);
            message.success('Device assigned to group successfully');
            setAssignModalVisible(false);
            setSelectedRowKeys([]);
            fetchDevices();
        } catch (err) {
            message.error('Failed to assign device to group');
        }
    };

    const columns = [
        {
            title: 'Device',
            render: (_: any, record: ClientNode) => (
                <a onClick={() => navigate(`/devices/${record.device.id}`)}>
                    {record.metric?.hostname || record.device.device_identifier || '---'}
                </a>
            )
        },
        {
            title: 'Status',
            render: (_: any, record: ClientNode) => {
                const status = record.device.status.toLowerCase();
                const color = status === 'online' ? 'green' : 'default';
                return <Badge color={color} text={status === 'online' ? 'Online' : 'Offline'} />;
            }
        },
        {
            title: 'Group',
            render: (_: any, record: ClientNode) => {
                const group = groups.find((g) => g.id === record.device.group_id);
                return group ? group.name : 'Unassigned';
            }
        },
        {
            title: 'Disk Space',
            render: (_: any, record: ClientNode) => {
                if (!record.metric) return '---';
                const free = record.metric.disk_free;
                const total = record.metric.disk_total;
                const percent = total ? Math.round((free / total) * 100) : 0;
                return (
                    <div>
                        <Progress
                            percent={percent}
                            size="small"
                            strokeColor={getProgressColor(percent)}
                        />
                        <div style={{ fontSize: 12 }}>
                            {formatBytes(free)} / {formatBytes(total)}
                        </div>
                    </div>
                );
            }
        },
        {
            title: 'Memory',
            render: (_: any, record: ClientNode) => {
                if (!record.metric) return '---';
                const total = record.metric.memory_total;
                const available = record.metric.memory_available;
                // Используем процент использования памяти: чем больше используется, тем "хуже"
                const usedPercent = total ? Math.round(100 - (available / total) * 100) : 0;
                return (
                    <div>
                        <Progress
                            percent={usedPercent}
                            size="small"
                            strokeColor={getMemoryProgressColor(usedPercent)}
                        />
                        <div style={{ fontSize: 12 }}>
                            {formatBytes(available)} free / {formatBytes(total)} total
                        </div>
                    </div>
                );
            }
        }
    ];

    return (
        <div style={{ padding: 24 }}>
            <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
                <Space>
                    <Select
                        style={{ width: 200 }}
                        value={groupFilter}
                        onChange={(value) => setGroupFilter(value)}
                    >
                        <Select.Option value="all">All Groups</Select.Option>
                        {groups.map((group) => (
                            <Select.Option key={group.id} value={group.id}>
                                {group.name}
                            </Select.Option>
                        ))}
                    </Select>
                    <Input
                        placeholder="Search by device name"
                        style={{ width: 200 }}
                        allowClear
                        value={searchText}
                        onChange={(e) => setSearchText(e.target.value)}
                    />
                </Space>
                <Space>
                    <Button
                        type="primary"
                        onClick={handleAssign}
                        disabled={selectedRowKeys.length !== 1}
                    >
                        Assign Device to Group
                    </Button>
                    <Button type="primary" onClick={fetchDevices}>
                        Refresh
                    </Button>
                </Space>
            </div>
            <Table
                rowKey={(record: ClientNode) => record.device.id}
                columns={columns}
                dataSource={filteredNodes}
                loading={loading}
                rowSelection={handleRowSelection}
                pagination={{
                    pageSize: 5,
                    showSizeChanger: true,
                    showTotal: (total: number) => `Total ${total} items`
                }}
            />
            <Modal
                title="Assign Device to Group"
                visible={assignModalVisible}
                onCancel={() => setAssignModalVisible(false)}
                footer={null}
                destroyOnClose
            >
                <Form form={form} layout="vertical" onFinish={onAssignFinish}>
                    <Form.Item
                        name="group_id"
                        label="Select Group"
                        rules={[{ required: true, message: 'Please select a group' }]}
                    >
                        <Select placeholder="Select a group">
                            {groups.map((group) => (
                                <Select.Option key={group.id} value={group.id}>
                                    {group.name}
                                </Select.Option>
                            ))}
                        </Select>
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" htmlType="submit">
                            Assign
                        </Button>
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    );
};
