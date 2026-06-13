/**
 * App 根组件
 * 整合所有页面区块组件
 */

import React, { useEffect } from 'react';
import { useLanguage } from './contexts/LanguageContext';
import { useCMSData } from './contexts/CMSContext';
import { initLazyLoading } from './utils/lazyLoad';
import { updateSiteTitle } from './services/i18n';

// 布局组件
import Header from './components/layout/Header';
import Footer from './components/layout/Footer';

// 页面区块组件
import Banner from './components/sections/Banner';
import Stats from './components/sections/Stats';
import About from './components/sections/About';
import Products from './components/sections/Products';
import Factory from './components/sections/Factory';
import Advantage from './components/sections/Advantage';
import Events from './components/sections/Events';
import Contact from './components/sections/Contact';

// 通用组件
import ErrorBoundary from './components/common/ErrorBoundary';
import Loading from './components/common/Loading';
import WhatsAppFloat from './components/common/WhatsAppFloat';

const AppContent: React.FC = () => {
    const { t, currentLang } = useLanguage();
    const { data, isLoading, error, refreshData } = useCMSData();

    // 语言切换时更新站点标题
    useEffect(() => {
        updateSiteTitle(data?.siteSettings);
    }, [currentLang, data?.siteSettings]);

    // 数据加载完成后初始化懒加载
    useEffect(() => {
        if (data && !isLoading) {
            const cleanup = initLazyLoading();
            return cleanup;
        }
    }, [data, isLoading]);

    if (isLoading) {
        return (
            <div className="app-loading">
                <Loading variant="spinner" text="Loading..." />
            </div>
        );
    }

    if (error && !data) {
        return (
            <div className="app-error">
                <p>Failed to load content. Please refresh the page.</p>
                <button onClick={refreshData}>Retry</button>
            </div>
        );
    }

    return (
        <>
            <Header />
            <main>
                <ErrorBoundary>
                    <Banner />
                    <Stats />
                    <About />
                    <Products />
                    <Factory />
                    <Advantage />
                    <Events />
                    <Contact />
                </ErrorBoundary>
            </main>
            <Footer />
            <WhatsAppFloat />
        </>
    );
};

function App() {
    return <AppContent />;
}

export default App;