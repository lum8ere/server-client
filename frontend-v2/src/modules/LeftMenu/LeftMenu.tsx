import React from 'react';
import { Menu } from 'antd';
import { useNavigate } from 'react-router-dom';
import { WindowsOutlined } from '@ant-design/icons';
import { ReactComponent as MyLogoIcon } from 'assets/mini_logo.svg';
interface LeftMenuProps {
    collapsed: boolean;
}

export const LeftMenu: React.FC<LeftMenuProps> = ({ collapsed }) => {
    const navigate = useNavigate();

    const logoBlock = (
        <div
            onClick={() => navigate('/devices')}
            style={{
                height: 40,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                borderBottom: '1px solid rgba(255, 255, 255, 0.1)',
                transition: 'all 0.2s'
            }}
        >
            <MyLogoIcon
                style={{
                    height: collapsed ? 32 : 40,
                    transition: 'all 0.2s'
                }}
            />
        </div>
    );

    // Пункты меню
    const menuItems = [
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
