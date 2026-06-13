/**
 * About 组件
 * 关于我们模块，包含公司介绍、认证列表和图片展示
 */

import React from 'react';
import { useLanguage } from '../../contexts/LanguageContext';
import { useCMSData } from '../../contexts/CMSContext';

const About: React.FC = () => {
    const { t } = useLanguage();
    const { getModule, getImageByCategory } = useCMSData();

    const aboutModule = getModule('about');

    // 获取关于图片
    const aboutImage = getImageByCategory('about')?.filename || '/uploads/1778806950157011000_IMG_4330.png';

    return (
        <section id="about">
            <div className="container">
                <div className="section-title">
                    <h2 id="aboutTitle" data-i18n="about.title">{t('about.title')}</h2>
                    <p id="aboutSubtitle" data-i18n="about.subtitle">{t('about.subtitle')}</p>
                </div>
                
                <div className="about-wrap">
                    <div className="about-text">
                        <h3 id="aboutLeftTitle" data-i18n="about.leftTitle">{t('about.leftTitle')}</h3>
                        
                        <div id="aboutContentContainer">
                            <p id="aboutParagraph1">
                                {aboutModule?.enContent || 
                                    'Founded in 2010, we are a professional manufacturer and exporter dedicated to the research, development, production and sales of high-standard medical devices. Our facility has passed ISO13485 medical quality system certification, CE EU certification, and FDA registration.'}
                            </p>
                            <p id="aboutParagraph2">
                                We specialize in surgical instruments, disposable medical supplies, hospital rehabilitation equipment, emergency medical devices and more. Our products are exported to hospitals, medical distributors, and healthcare institutions in over 80 countries worldwide.
                            </p>
                            <p id="aboutParagraph3">
                                We provide comprehensive OEM & ODM customization services, ensuring strict quality control throughout production while offering one-stop procurement solutions for global customers.
                            </p>
                        </div>

                        <div className="cert-list">
                            <span className="cert-item" data-i18n="about.ce">{t('about.ce')}</span>
                            <span className="cert-item" data-i18n="about.fda">{t('about.fda')}</span>
                            <span className="cert-item" data-i18n="about.iso">{t('about.iso')}</span>
                            <span className="cert-item" data-i18n="about.sgs">{t('about.sgs')}</span>
                        </div>
                    </div>

                    <div className="about-img">
                        <img 
                            data-src={aboutImage} 
                            alt="Company Team" 
                            className="lazy"
                            style={{
                                width: '100%',
                                height: 'auto',
                                objectFit: 'cover',
                                borderRadius: '10px',
                                minHeight: '200px',
                                background: '#f5f5f5'
                            }}
                        />
                        <div className="about-badge">
                            <span className="num">80+</span>
                            <span className="text">Countries</span>
                        </div>
                    </div>
                </div>
            </div>
        </section>
    );
};

export default About;