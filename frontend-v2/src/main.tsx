import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import { Provider } from 'react-redux';
import { store } from 'store';
import { ConfigProvider } from 'antd';
import { App } from 'App';
import './index.css';

const rootElement = document.getElementById('root');

const NewClientApp: React.FC = () => (
    <StrictMode>
        <ConfigProvider
            theme={{
                token: {
                    colorPrimary: '#7cccab'
                }
            }}
        >
            <Provider store={store}>
                <BrowserRouter>
                    <App />
                </BrowserRouter>
            </Provider>
        </ConfigProvider>
    </StrictMode>
);

if (rootElement) {
    const container = document.getElementById('root');
    const root = createRoot(container!);
    root.render(<NewClientApp />);
} else {
    console.error('Root element not found!');
}

export default NewClientApp;
