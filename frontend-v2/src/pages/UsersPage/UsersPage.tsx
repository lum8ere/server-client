import React, { useEffect, useState } from 'react';
import { useSelector } from 'react-redux';
import { Table, Button, Modal, Form, Input, Select, message } from 'antd';
import instance from 'service/api';
import { RootState } from 'store';

interface User {
    id: string;
    username: string;
    email: string;
    role_code: string; // предполагается, что в этом поле хранится role code или id
    created_at: string;
}

interface Role {
    id: string;
    name: string;
    code: string;
}

export const UsersPage: React.FC = () => {
    const token = useSelector((state: RootState) => state.auth.token);
    const [users, setUsers] = useState<User[]>([]);
    const [roles, setRoles] = useState<Role[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [modalVisible, setModalVisible] = useState<boolean>(false);
    const [form] = Form.useForm();

    const fetchUsers = async () => {
        setLoading(true);
        try {
            const res = await instance.get<User[]>('/api/users', {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            });
            setUsers(res.data);
        } catch (error) {
            message.error('Failed to fetch users');
        } finally {
            setLoading(false);
        }
    };

    const fetchRoles = async () => {
        try {
            const res = await instance.get<Role[]>('/api/dicts/roles');
            setRoles(res.data);
        } catch (error) {
            message.error('Failed to fetch roles');
        }
    };

    useEffect(() => {
        fetchUsers();
        fetchRoles();
    }, []);

    const handleAddUser = async (values: any) => {
        try {
            // Отправляем запрос на создание пользователя
            await instance.post('/api/auth/register', values, {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            });
            message.success('User added successfully');
            setModalVisible(false);
            form.resetFields();
            fetchUsers();
        } catch (error) {
            message.error('Failed to add user');
        }
    };

    const columns = [
        { title: 'Username', dataIndex: 'username' },
        { title: 'Email', dataIndex: 'email' },
        {
            title: 'Role',
            dataIndex: 'role_code',
            render: (roleCode: string) => {
                // Поиск роли в списке по коду (или id, если нужно)
                const roleObj = roles.find((r) => r.code === roleCode || r.id === roleCode);
                return roleObj ? roleObj.name : roleCode;
            }
        },
        {
            title: 'Created At',
            dataIndex: 'created_at',
            render: (dateStr: string) => {
                const date = new Date(dateStr);
                return date.toUTCString(); // можно изменить формат, если нужно
            }
        }
    ];

    return (
        <div style={{ padding: 24 }}>
            <h2>User Management</h2>
            <Button
                type="primary"
                onClick={() => setModalVisible(true)}
                style={{ marginBottom: 16 }}
            >
                Add User
            </Button>
            <Table rowKey="id" columns={columns} dataSource={users} loading={loading} />
            <Modal
                title="Add New User"
                visible={modalVisible}
                onCancel={() => setModalVisible(false)}
                footer={null}
            >
                <Form form={form} layout="vertical" onFinish={handleAddUser}>
                    <Form.Item
                        name="username"
                        label="Username"
                        rules={[{ required: true, message: 'Please enter a username' }]}
                    >
                        <Input />
                    </Form.Item>
                    <Form.Item
                        name="email"
                        label="Email"
                        rules={[{ required: true, message: 'Please enter an email' }]}
                    >
                        <Input type="email" />
                    </Form.Item>
                    <Form.Item
                        name="password"
                        label="Password"
                        rules={[{ required: true, message: 'Please enter a password' }]}
                    >
                        <Input.Password />
                    </Form.Item>
                    <Form.Item
                        name="role"
                        label="Role"
                        rules={[{ required: true, message: 'Please select a role' }]}
                    >
                        <Select placeholder="Select a role">
                            {roles.map((role) => (
                                <Select.Option key={role.code} value={role.code}>
                                    {role.name}
                                </Select.Option>
                            ))}
                        </Select>
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" htmlType="submit">
                            Add User
                        </Button>
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    );
};
