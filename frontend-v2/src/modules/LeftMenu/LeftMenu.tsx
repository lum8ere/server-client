import React from 'react';
import { Menu } from 'antd';
import { useNavigate } from 'react-router-dom';
import {
    HomeOutlined,
    AppstoreOutlined,
    UserOutlined,
    SettingOutlined,
    WindowsOutlined
} from '@ant-design/icons';

interface LeftMenuProps {
    collapsed: boolean;
}

export const LeftMenu: React.FC<LeftMenuProps> = ({ collapsed }) => {
    const navigate = useNavigate();

    // Логотип с анимацией
    const logoBlock = (
        <div
            style={{
                height: 64,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                borderBottom: '1px solid rgba(255, 255, 255, 0.1)',
                transition: 'all 0.2s'
            }}
        >
            <span
                style={{
                    color: '#fff',
                    fontSize: collapsed ? 20 : 22,
                    fontWeight: 600,
                    transition: 'all 0.2s'
                }}
            >
                {collapsed ? '⚡' : 'MyApp'}
            </span>
        </div>
    );

    // Пункты меню
    const menuItems = [
        // {
        //     key: '/',
        //     icon: <HomeOutlined style={{ fontSize: 18, color: '#fff' }} />,
        //     label: <span style={{ color: '#fff' }}>Главная</span>
        // },
        {
            key: '/devices',
            icon: <WindowsOutlined style={{ fontSize: 18, color: '#fff' }} />,
            label: <span style={{ color: '#fff' }}>Devices</span>
        }
    ];

    return (
        <>
            {logoBlock}
            <Menu
                theme="dark"
                mode="inline"
                defaultSelectedKeys={['/']}
                inlineCollapsed={collapsed}
                items={menuItems}
                style={{
                    background: '#001529',
                    borderRight: 0,
                    padding: '8px 0'
                }}
                onClick={({ key }) => navigate(key)}
            />
        </>
    );
};
