/**
 * Banner 组件
 * 首页 Banner 区域，支持动态背景图片和多语言内容
 */

import React from 'react';
import { useLanguage } from '../../contexts/LanguageContext';
import { useCMSData } from '../../contexts/CMSContext';
import Button from '../common/Button';

const Banner: React.FC = () => {
    const { t } = useLanguage();
    const { getModule, getImageByCategory } = useCMSData();

    const bannerModule = getModule('banner');
    
    // 获取背景图片
    const backgroundImage = bannerModule?.imagePath 
        ? `/uploads/${bannerModule.imagePath}`
        : getImageByCategory('banner')?.filename;

    // 获取标题和副标题
    const title = bannerModule 
        ? (bannerModule.enTitle && bannerModule.zhTitle 
            ? (document.documentElement.lang === 'zh' ? bannerModule.zhTitle : bannerModule.enTitle)
            : t('banner.title'))
        : t('banner.title');

    const subtitle = bannerModule
        ? (bannerModule.enSubtitle && bannerModule.zhSubtitle
            ? (document.documentElement.lang === 'zh' ? bannerModule.zhSubtitle : bannerModule.enSubtitle)
            : t('banner.subtitle'))
        : t('banner.subtitle');

    return (
        <section 
            className="banner"
            id="home"
            style={{
                background: backgroundImage
                    ? `linear-gradient(135deg,rgba(10,92,173,0.9) 0%,rgba(6,58,117,0.85) 100%),url('/uploads/${backgroundImage}') center/cover no-repeat`
                    : 'linear-gradient(135deg,rgba(10,92,173,0.9) 0%,rgba(6,58,117,0.85) 100%)'
            }}
        >
            <div className="container banner-text">
                <h1 
                    data-i18n="banner.title"
                    data-i18n-html
                    dangerouslySetInnerHTML={{ __html: title }}
                />
                <p data-i18n="banner.subtitle">{subtitle}</p>
                <div className="banner-actions">
                    <Button href="#products" variant="primary">
                        {t('banner.viewProducts')}
                    </Button>
                    <Button href="#contact" variant="outline">
                        {t('banner.getQuote')}
                    </Button>
                </div>
            </div>
        </section>
    );
};

export default Banner;