import { PropsWithChildren } from 'react';
import { Layout } from 'antd';
import { AppHeader } from 'modules/Header/Header';

const { Content } = Layout;

interface DefaultLayoutProps extends PropsWithChildren {}

export const DefaultLayout = ({ children }: DefaultLayoutProps) => {
    return (
        <Layout style={{ minHeight: '100vh' }}>
            <Layout style={{ position: 'relative' }}>
                <AppHeader />
                <Content>{children}</Content>
            </Layout>
        </Layout>
    );
};
