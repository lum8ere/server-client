import React, { useEffect, useState } from 'react';
import { Table, Button, Modal, Form, Input, message } from 'antd';
import instance from 'service/api';

interface DeviceGroup {
    id: string;
    name: string;
    description?: string;
    created_at: string;
}

export const DeviceGroupsPage: React.FC = () => {
    const [groups, setGroups] = useState<DeviceGroup[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [modalVisible, setModalVisible] = useState<boolean>(false);
    const [editingGroup, setEditingGroup] = useState<DeviceGroup | null>(null);
    const [form] = Form.useForm();

    const fetchGroups = async () => {
        setLoading(true);
        try {
            const res = await instance.get<DeviceGroup[]>('/api/device-groups');
            setGroups(res.data);
        } catch (err) {
            message.error('Failed to fetch device groups');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchGroups();
    }, []);

    const onFinish = async (values: any) => {
        try {
            if (editingGroup) {
                const payload = { id: editingGroup.id, ...values };
                await instance.put<DeviceGroup>('/api/device-groups', payload);
                message.success('Device group updated successfully');
            } else {
                await instance.post<DeviceGroup>('/api/device-groups', values);
                message.success('Device group created successfully');
            }
            setModalVisible(false);
            form.resetFields();
            fetchGroups();
        } catch (err) {
            message.error('Failed to save device group');
        }
    };

    const handleEdit = (group: DeviceGroup) => {
        setEditingGroup(group);
        form.setFieldsValue({
            name: group.name,
            description: group.description
        });
        setModalVisible(true);
    };

    const handleDelete = async (group: DeviceGroup) => {
        try {
            await instance.delete('/api/device-groups', { params: { id: group.id } });
            message.success('Device group deleted successfully');
            fetchGroups();
        } catch (err) {
            message.error('Failed to delete device group');
        }
    };

    const columns = [
        {
            title: 'Group Name',
            dataIndex: 'name'
        },
        {
            title: 'Description',
            dataIndex: 'description'
        },
        {
            title: 'Created At',
            dataIndex: 'created_at',
            render: (date: string) => new Date(date).toLocaleString()
        },
        {
            title: 'Actions',
            render: (_: any, record: DeviceGroup) => (
                <>
                    <Button type="link" onClick={() => handleEdit(record)}>
                        Edit
                    </Button>
                    <Button type="link" danger onClick={() => handleDelete(record)}>
                        Delete
                    </Button>
                </>
            )
        }
    ];

    return (
        <div style={{ padding: 24 }}>
            <h2>Device Groups</h2>
            <Button
                type="primary"
                onClick={() => {
                    setEditingGroup(null);
                    setModalVisible(true);
                }}
                style={{ marginBottom: 16 }}
            >
                Add Group
            </Button>
            <Table
                rowKey="id"
                columns={columns}
                dataSource={groups}
                loading={loading}
                pagination={{ pageSize: 5 }}
            />
            <Modal
                title={editingGroup ? 'Edit Device Group' : 'Add Device Group'}
                visible={modalVisible}
                onCancel={() => setModalVisible(false)}
                footer={null}
                destroyOnClose
            >
                <Form
                    form={form}
                    layout="vertical"
                    onFinish={onFinish}
                    initialValues={{ name: '', description: '' }}
                >
                    <Form.Item
                        name="name"
                        label="Group Name"
                        rules={[{ required: true, message: 'Please enter the group name' }]}
                    >
                        <Input />
                    </Form.Item>
                    <Form.Item name="description" label="Description">
                        <Input />
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" htmlType="submit">
                            {editingGroup ? 'Update' : 'Create'}
                        </Button>
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    );
};
