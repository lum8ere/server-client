import React, { useEffect, useState } from 'react';
import { Table, Button, Space, Badge, Select, Input, Row, Col, Progress, Typography } from 'antd';
import { useNavigate } from 'react-router-dom';
import instance from 'service/api';

const { Option } = Select;
const { Title } = Typography;

// Типы данных
interface ClientNode {
    id: string; // Уникальный ID
    nodeName: string; // 'LABORATO_STAND', 'LABORATO-THINK' и т.д.
    team: string | null; // 'lab_test' или null/undefined
    status: 'online' | 'offline';
    issues: number; // Количество каких-то "ошибок"/"предупреждений"
    diskSpaceAvailable: number; // Свободное место (GB), для примера
    diskSpaceTotal: number; // Всего места (GB), чтобы считать % занято
    operatingSystem: string; // 'Майкрософт Windows 11 ...'
    privateIpAddress: string; // '192.168.88.20'
}

export const ClientsList: React.FC = () => {
    const [nodes, setNodes] = useState<ClientNode[]>([]);
    const [filteredNodes, setFilteredNodes] = useState<ClientNode[]>([]);
    const [selectedTeam, setSelectedTeam] = useState('All nodes');
    const [searchValue, setSearchValue] = useState('');
    const navigate = useNavigate();

    // Получаем список узлов
    useEffect(() => {
        fetchNodes();
    }, []);

    const fetchNodes = async () => {
        try {
            // Допустим, ваш бэкенд теперь возвращает массив ClientNode
            const res = await instance.get<ClientNode[]>('/api/clients');
            setNodes(res.data);
            setFilteredNodes(res.data);
        } catch (err) {
            console.error('Ошибка получения списка клиентов', err);
        }
    };

    // Хэндлеры фильтра/поиска
    useEffect(() => {
        // Фильтрация по team
        let data = [...nodes];
        if (selectedTeam !== 'All nodes') {
            data = data.filter((item) => item.team === selectedTeam);
        }
        // Поиск по имени/Id/IP
        if (searchValue.trim()) {
            const lower = searchValue.toLowerCase();
            data = data.filter(
                (item) =>
                    item.nodeName.toLowerCase().includes(lower) ||
                    item.id.toLowerCase().includes(lower) ||
                    item.privateIpAddress.toLowerCase().includes(lower)
            );
        }
        setFilteredNodes(data);
    }, [nodes, selectedTeam, searchValue]);

    // Колонки
    const columns = [
        {
            title: 'Node',
            dataIndex: 'nodeName',
            render: (val: string) => val || '---'
        },
        {
            title: 'Team',
            dataIndex: 'team',
            render: (val: string | null) => val || '---'
        },
        {
            title: 'Status',
            dataIndex: 'status',
            render: (value: 'online' | 'offline') => {
                const color = value === 'online' ? 'green' : 'gray';
                return <Badge color={color} text={value === 'online' ? 'Online' : 'Offline'} />;
            }
        },
        {
            title: 'Disk space available',
            // У нас есть diskSpaceAvailable и diskSpaceTotal. Можно отображать просто цифры,
            // а можно сделать цветовую индикацию (Progress) с процентами.
            render: (_: any, record: ClientNode) => {
                const used = record.diskSpaceTotal - record.diskSpaceAvailable;
                const percent = (record.diskSpaceAvailable / record.diskSpaceTotal) * 100;
                return (
                    <Space direction="horizontal">
                        {/* Отображаем доступное место, например "30 GB" */}
                        <div>{record.diskSpaceAvailable} GB</div>
                        {/* Прогресс-бар (как индикатор) */}
                        <Progress
                            percent={Math.round(percent)}
                            size="small"
                            showInfo={false}
                            strokeColor={percent < 15 ? 'red' : percent < 30 ? 'orange' : 'green'}
                        />
                    </Space>
                );
            }
        },
        {
            title: 'Operating system',
            dataIndex: 'operatingSystem',
            ellipsis: true // на всякий случай, если строка длинная
        }
        // {
        //     title: 'Private IP address',
        //     dataIndex: 'privateIpAddress',
        //     render: (val: string | undefined) => val || '---'
        // },
    ];

    const rowSelection = {
        onChange: (selectedRowKeys: React.Key[], selectedRows: ClientNode[]) => {
            // пример, если захотите что-то делать при выборе строк
            console.log('selectedRowKeys:', selectedRowKeys, 'selectedRows:', selectedRows);
        }
    };

    const downloadClient = () => {
        // переход по ссылке /download/client
        window.open('/download/client', '_blank');
    };

    return (
        <div style={{ padding: 16 }}>
            {/* Заголовок (All teams) + кнопка Add nodes */}
            {/* <Row justify="space-between" align="middle" style={{ marginBottom: 16 }}>
                <Col>
                    <Title level={4} style={{ margin: 0 }}>
                        All teams
                    </Title>
                </Col>
                <Col></Col>
            </Row> */}

            {/* Кол-во нод и кнопка Export */}
            {/* <Row style={{ marginBottom: 16 }}>
                <Col>
                    <Space>
                        <span>{filteredNodes.length} nodes</span>
                        <a href="#">Export nodes</a>
                    </Space>
                </Col>
            </Row> */}

            {/* Фильтр по команде (All nodes / lab_test / и т.д.) и поиск */}
            {/* <Space style={{ marginBottom: 16 }}>
                <Select
                    style={{ width: 150 }}
                    value={selectedTeam}
                    onChange={(val) => setSelectedTeam(val)}
                >
                    <Option value="All nodes">All nodes</Option>
                    <Option value="lab_test">lab_test</Option>
                    <Option value="other_team">other_team</Option>
                </Select>
                <Input
                    style={{ width: 250 }}
                    placeholder="Search name, nodename, UUID"
                    value={searchValue}
                    onChange={(e) => setSearchValue(e.target.value)}
                />
            </Space> */}

            <div style={{ marginBottom: 16 }}>
                <Button onClick={downloadClient}>Add nodes</Button>
            </div>

            <Table
                rowKey="id"
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
