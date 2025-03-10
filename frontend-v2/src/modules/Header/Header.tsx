import React from 'react';
import { Dropdown, Menu, Avatar, Space } from 'antd';
import { MenuFoldOutlined, MenuUnfoldOutlined, UserOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';

interface AppHeaderProps {
    collapsed: boolean;
    toggleCollapsed: () => void;
}

export const AppHeader: React.FC<AppHeaderProps> = ({ collapsed, toggleCollapsed }) => {
    const navigate = useNavigate();

    const profileMenu = (
        <Menu
            items={[
                {
                    key: 'profile',
                    label: 'Профиль',
                    onClick: () => navigate('/profile')
                },
                {
                    key: 'settings',
                    label: 'Настройки',
                    onClick: () => navigate('/settings')
                },
                {
                    type: 'divider'
                },
                {
                    key: 'logout',
                    label: 'Выход',
                    onClick: () => navigate('/login')
                }
            ]}
        />
    );

    return (
        <div
            style={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                width: '100%',
                padding: '0 24px',
                background: '#fff',
                height: 64,
                // padding: 0,
                boxShadow: '0 2px 8px rgba(0,0,0,0.15)'
            }}
        >
            {/* Кнопка сворачивания Sider */}
            <div onClick={toggleCollapsed} style={{ cursor: 'pointer', fontSize: 18 }}>
                {collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            </div>

            {/* Иконка аккаунта, стилизованная в серый цвет */}
            <Dropdown overlay={profileMenu} trigger={['click']}>
                <Space style={{ cursor: 'pointer' }}>
                    <Avatar icon={<UserOutlined />} style={{ backgroundColor: '#8c8c8c' }} />
                </Space>
            </Dropdown>
        </div>
    );
};
