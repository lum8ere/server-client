import React from 'react';
import ReactDOM from 'react-dom/client';
import {App} from './App';
import 'antd/dist/reset.css'; // Сброс стилей antd (Ant Design v5)
import './index.css';

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
