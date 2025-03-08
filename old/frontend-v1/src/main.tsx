import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import { App } from 'App';
import './index.css';

const rootElement = document.getElementById('root');

const NewClientApp: React.FC = () => (
    <StrictMode>
        <BrowserRouter>
            <App />
        </BrowserRouter>
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
