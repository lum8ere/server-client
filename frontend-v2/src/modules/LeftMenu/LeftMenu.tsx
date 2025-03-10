import React from 'react';
import { Menu } from 'antd';
import { useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { WindowsOutlined, UserOutlined } from '@ant-design/icons';
import { ReactComponent as MyLogoIcon } from 'assets/logo.svg';
import { RootState } from 'store';

interface LeftMenuProps {
    collapsed: boolean;
}

export const LeftMenu: React.FC<LeftMenuProps> = ({ collapsed }) => {
    const navigate = useNavigate();
    const { user } = useSelector((state: RootState) => state.auth);

    const logoBlock = (
        <div
            onClick={() => navigate('/devices')}
            style={{
                height: 40,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                borderBottom: '1px solid rgba(255,255,255,0.1)',
                transition: 'all 0.2s'
            }}
        >
            <MyLogoIcon style={{ height: collapsed ? 32 : 40, transition: 'all 0.2s' }} />
        </div>
    );

    // Всегда показываем Devices
    const menuItems = [
        {
            key: '/devices',
            icon: <WindowsOutlined style={{ fontSize: 18, color: '#fff' }} />,
            label: <span style={{ color: '#fff' }}>Devices</span>
        }
    ];

    // Добавляем пункт Users только если роль пользователя "admin" (без учета регистра)
    if (user && user.role.toLowerCase() === 'admin') {
        menuItems.push({
            key: '/users',
            icon: <UserOutlined style={{ fontSize: 18, color: '#fff' }} />,
            label: <span style={{ color: '#fff' }}>Users</span>
        });
    }

    return (
        <>
            {logoBlock}
            <Menu
                theme="dark"
                mode="inline"
                defaultSelectedKeys={['/devices']}
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
