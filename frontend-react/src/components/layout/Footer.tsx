/**
 * Footer 组件
 * 包含页脚 Logo、快速链接、产品链接、服务链接和版权信息
 */

import React from 'react';
import { useLanguage } from '../../contexts/LanguageContext';
import { useCMSData } from '../../contexts/CMSContext';

const Footer: React.FC = () => {
    const { currentLang, t } = useLanguage();
    const { getSiteSettings } = useCMSData();

    const siteSettings = getSiteSettings();
    const logoText = currentLang === 'zh' 
        ? siteSettings?.zhSiteLogo 
        : siteSettings?.enSiteLogo;
    const logoColor = siteSettings?.siteLogoColor || '#06a499';

    return (
        <footer>
            <div className="container">
                <div className="footer-grid">
                    {/* 关于我们 */}
                    <div className="footer-about">
                        <h3 id="footerLogo">
                            {logoText}
                            {currentLang === 'en' && logoText?.toUpperCase() === logoText && logoText.includes('MEDICAL') && (
                                <span id="footerLogoSuffix" style={{ color: logoColor }}>PRO</span>
                            )}
                        </h3>
                        <p>
                            {currentLang === 'zh' 
                                ? '中国专业医疗器械制造商和出口商。我们为世界各地客户提供高品质医疗设备、手术用品和全面的 OEM/ODM 服务。'
                                : 'Professional medical device manufacturer and exporter from China. We provide high-quality medical equipment, surgical supplies, and comprehensive OEM/ODM services to customers worldwide.'
                            }
                        </p>
                    </div>

                    {/* 快速链接 */}
                    <div className="footer-col">
                        <h4 data-i18n="footer.quickLinks">{t('footer.quickLinks')}</h4>
                        <a href="#home">{t('nav.home')}</a>
                        <a href="#about">{t('nav.about')}</a>
                        <a href="#products">{t('nav.products')}</a>
                        <a href="#contact">{t('nav.contact')}</a>
                    </div>

                    {/* 产品链接 */}
                    <div className="footer-col">
                        <h4 data-i18n="footer.products">{t('footer.products')}</h4>
                        <a href="#">Surgical Instruments</a>
                        <a href="#">Hospital Furniture</a>
                        <a href="#">Medical Disposables</a>
                        <a href="#">Rehabilitation Equipment</a>
                    </div>

                    {/* 服务链接 */}
                    <div className="footer-col">
                        <h4 data-i18n="footer.services">{t('footer.services')}</h4>
                        <a href="#" data-i18n="footer.oem">{t('footer.oem')}</a>
                        <a href="#" data-i18n="footer.wholesale">{t('footer.wholesale')}</a>
                        <a href="#" data-i18n="footer.shipping">{t('footer.shipping')}</a>
                        <a href="#" data-i18n="footer.support">{t('footer.support')}</a>
                    </div>
                </div>

                {/* 版权信息 */}
                <div className="footer-bottom">
                    <p data-i18n="footer.copyright">{t('footer.copyright')}</p>
                </div>
            </div>
        </footer>
    );
};

export default Footer;