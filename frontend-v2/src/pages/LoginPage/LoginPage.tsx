import React, { useState } from 'react';
import { Form, Input, Button, message } from 'antd';
import { useDispatch, useSelector } from 'react-redux';
import { AppDispatch, RootState } from 'store';
import { useNavigate } from 'react-router-dom';
import { loginUser } from 'store/authSlice';
import { ReactComponent as Logo } from 'assets/logo.svg';

export const LoginPage: React.FC = () => {
    const navigate = useNavigate();
    const dispatch = useDispatch<AppDispatch>();
    const { loading, error } = useSelector((state: RootState) => state.auth);

    const onFinish = async (values: { email: string; password: string }) => {
        try {
            const resultAction = await dispatch(loginUser(values));
            if (loginUser.fulfilled.match(resultAction)) {
                message.success('Logged in successfully');
                navigate('/');
            } else {
                message.error(resultAction.payload as string);
            }
        } catch (err) {
            message.error('Login failed');
        }
    };

    return (
        <div
            style={{
                width: '350px',
                margin: '0 auto',
                padding: 24,
                textAlign: 'center'
            }}
        >
            <Logo style={{ width: '350px', height: '100px', marginBottom: 16 }} />
            <h2>Login</h2>
            <Form onFinish={onFinish}>
                <Form.Item
                    name="email"
                    rules={[{ required: true, message: 'Please enter your email!' }]}
                >
                    <Input placeholder="Email" />
                </Form.Item>
                <Form.Item
                    name="password"
                    rules={[{ required: true, message: 'Please enter your password!' }]}
                >
                    <Input.Password placeholder="Password" />
                </Form.Item>
                <Form.Item>
                    <Button type="primary" htmlType="submit" loading={loading} block>
                        Login
                    </Button>
                </Form.Item>
            </Form>
            {error && <p style={{ color: 'red' }}>{error}</p>}
        </div>
    );
};
