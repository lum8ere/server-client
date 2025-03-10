import React, { PropsWithChildren, useState } from 'react';
import { Layout } from 'antd';
import { LeftMenu } from 'modules/LeftMenu/LeftMenu';
import { AppHeader } from 'modules/Header/Header';

const { Header, Sider, Content } = Layout;

export const DefaultLayout: React.FC<PropsWithChildren> = ({ children }) => {
    const [collapsed, setCollapsed] = useState(false);
    const toggleCollapsed = () => setCollapsed(!collapsed);

    return (
        // style={{ height: '100vh' }}
        <Layout style={{ minHeight: '100vh' }}>
            <Sider
                width="15%"
                style={{ background: '#001529', color: '#fff', padding: '16px' }}
                collapsed={collapsed}
            >
                <LeftMenu collapsed={collapsed} />
            </Sider>
            <Layout>
                <AppHeader collapsed toggleCollapsed={toggleCollapsed} />
                <Content style={{ padding: 24, background: '#f0f2f5', overflow: 'auto' }}>
                    {children}
                </Content>
            </Layout>
        </Layout>
    );
};

export default DefaultLayout;
