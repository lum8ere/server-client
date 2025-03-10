import React, { useEffect, useState } from 'react';
import { Form, Input, Button, message, Card, Descriptions } from 'antd';
import instance from 'service/api';
import { useSelector, useDispatch } from 'react-redux';
import { RootState } from 'store';
import { setToken } from 'store/authSlice';

interface ProfileData {
    id: string;
    username: string;
    email: string;
    role_code: string;
    created_at: string;
    updated_at: string;
}

export const ProfilePage: React.FC = () => {
    const [profile, setProfile] = useState<ProfileData | null>(null);
    const [editMode, setEditMode] = useState<boolean>(false);
    const [form] = Form.useForm();
    const dispatch = useDispatch();
    const { token, user } = useSelector((state: RootState) => state.auth);

    useEffect(() => {
        // Fetch profile data on component mount
        const fetchProfile = async () => {
            try {
                const userId = user?.user_id || '';
                const res = await instance.get<ProfileData>('/api/profile', {
                    params: { user_id: userId }
                });
                setProfile(res.data);
                form.setFieldsValue({
                    username: res.data.username,
                    email: res.data.email
                });
            } catch (err) {
                message.error('Failed to fetch profile data');
            }
        };

        fetchProfile();
    }, [form, user]);

    const onFinish = async (values: any) => {
        try {
            const userId = user?.user_id || '';
            const payload = { user_id: userId, ...values };
            const res = await instance.put<ProfileData>('/api/profile', payload);
            message.success('Profile updated successfully');
            setProfile(res.data);
            setEditMode(false);
        } catch (err) {
            message.error('Failed to update profile');
        }
    };

    return (
        <Card
            title="Profile"
            extra={
                !editMode && (
                    <Button type="primary" onClick={() => setEditMode(true)}>
                        Edit Profile
                    </Button>
                )
            }
        >
            {editMode ? (
                <Form
                    form={form}
                    layout="vertical"
                    onFinish={onFinish}
                    initialValues={{ username: '', email: '' }}
                >
                    <Form.Item
                        name="username"
                        label="Username"
                        rules={[{ required: true, message: 'Please enter your username' }]}
                    >
                        <Input />
                    </Form.Item>
                    <Form.Item
                        name="email"
                        label="Email"
                        rules={[{ required: true, message: 'Please enter your email' }]}
                    >
                        <Input type="email" />
                    </Form.Item>
                    <Form.Item name="password" label="New Password">
                        <Input.Password placeholder="Leave blank if unchanged" />
                    </Form.Item>
                    <Form.Item
                        name="confirmPassword"
                        label="Confirm New Password"
                        dependencies={['password']}
                        rules={[
                            ({ getFieldValue }) => ({
                                validator(_, value) {
                                    if (!value || getFieldValue('password') === value) {
                                        return Promise.resolve();
                                    }
                                    return Promise.reject(new Error('Passwords do not match'));
                                }
                            })
                        ]}
                    >
                        <Input.Password placeholder="Confirm new password" />
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" htmlType="submit" style={{ marginRight: 8 }}>
                            Save
                        </Button>
                        <Button
                            onClick={() => {
                                setEditMode(false);
                                form.setFieldsValue({
                                    username: profile?.username,
                                    email: profile?.email
                                });
                            }}
                        >
                            Cancel
                        </Button>
                    </Form.Item>
                </Form>
            ) : (
                <>
                    {profile ? (
                        <Descriptions bordered column={1}>
                            <Descriptions.Item label="Username">
                                {profile.username}
                            </Descriptions.Item>
                            <Descriptions.Item label="Email">{profile.email}</Descriptions.Item>
                            <Descriptions.Item label="Role">{profile.role_code}</Descriptions.Item>
                            <Descriptions.Item label="Created At">
                                {new Date(profile.created_at).toUTCString()}
                            </Descriptions.Item>
                            <Descriptions.Item label="Updated At">
                                {new Date(profile.updated_at).toUTCString()}
                            </Descriptions.Item>
                        </Descriptions>
                    ) : (
                        <p>No profile data available.</p>
                    )}
                </>
            )}
        </Card>
    );
};
