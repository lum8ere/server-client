import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Button, Layout, Space } from 'antd';
import './Header.scss';

export const AppHeader: React.FC = () => {
    const navigate = useNavigate();
    return (
        <Layout.Header className="header">
            <div className="header__center">
                <Space.Compact>
                    <Button
                        key={'home_button'}
                        type="text"
                        className="header__nav-button"
                        onClick={() => navigate('/')}
                    >
                        Home
                    </Button>
                </Space.Compact>
            </div>

            <div className="header__right">SERVICE TIME</div>
        </Layout.Header>
    );
};
