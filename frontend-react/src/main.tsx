/**
 * React 应用主入口
 * 使用 React 18 Concurrent 模式
 */

import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import { LanguageProvider } from './contexts/LanguageContext';
import { CMSProvider } from './contexts/CMSContext';
import './styles/main.css';

// 创建根容器
const rootElement = document.getElementById('root');

if (!rootElement) {
    throw new Error('Failed to find root element');
}

// 使用 React 18 createRoot API
const root = ReactDOM.createRoot(rootElement);

root.render(
    <React.StrictMode>
        <LanguageProvider>
            <CMSProvider>
                <App />
            </CMSProvider>
        </LanguageProvider>
    </React.StrictMode>
);