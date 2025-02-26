import React, { useEffect, useState } from 'react';
import { Table, Space } from 'antd';
import { Link } from 'react-router-dom';
import { fetchClients } from 'services/api';
import { Client } from 'types';

const ClientList: React.FC = () => {
    const [clients, setClients] = useState<Client[]>([]);
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        const getClients = async () => {
            try {
                const data = await fetchClients();
                setClients(data);
            } catch (err) {
                console.error('Ошибка получения клиентов', err);
            } finally {
                setLoading(false);
            }
        };
        getClients();
    }, []);

    const columns = [
        {
            title: 'ID',
            dataIndex: 'ID',
            key: 'ID'
        },
        {
            title: 'Публичный IP',
            dataIndex: 'IP',
            key: 'IP'
        },
        {
            title: 'Действие',
            key: 'action',
            render: (_: any, record: Client) => (
                <Space size="middle">
                    <Link to={`/client/${record.ID}`}>Выбрать</Link>
                </Space>
            )
        }
    ];

    return (
        <div style={{ padding: 20 }}>
            <h1>Подключённые клиенты</h1>
            <Table dataSource={clients} columns={columns} rowKey="ID" loading={loading} />
        </div>
    );
};

export default ClientList;
