/**
 * Advantage 组件
 * 核心优势网格（医疗认证、OEM/ODM、快速交付、售后支持）
 */

import React from 'react';
import { useLanguage } from '../../contexts/LanguageContext';
import { useCMSData } from '../../contexts/CMSContext';
import LazyImage from '../common/LazyImage';

interface AdvantageItem {
    image: string;
    title: string;
    description: string;
}

const Advantage: React.FC = () => {
    const { t } = useLanguage();
    const { getModule } = useCMSData();

    const advantageModule = getModule('advantage');

    const advantages: AdvantageItem[] = [
        {
            image: '/uploads/1778807177495650000_favicon-32x32.png',
            title: 'Medical Certifications',
            description: 'All products pass international medical certifications for safe global market distribution.'
        },
        {
            image: '/uploads/1778807177498171000_icon.png',
            title: 'OEM/ODM Service',
            description: 'Professional R&D team providing one-stop customized medical device solutions.'
        },
        {
            image: '/uploads/1778807177499079000_pdf.png',
            title: 'Fast Global Delivery',
            description: 'Professional logistics partners supporting DDP/FOB/CIF worldwide shipping.'
        },
        {
            image: '/uploads/1778807177501456000_qrcode-for-doocs.jpg',
            title: 'After-sales Support',
            description: 'Professional foreign trade team with 24-hour online service and technical support.'
        }
    ];

    return (
        <section id="advantage" className="advantage-section">
            <div className="container">
                <div className="section-title">
                    <h2 data-i18n="advantage.title">{t('advantage.title')}</h2>
                    <p data-i18n="advantage.subtitle">{t('advantage.subtitle')}</p>
                </div>

                <div className="adv-grid">
                    {advantages.map((adv, index) => (
                        <div key={index} className="adv-item">
                            <LazyImage
                                dataSrc={adv.image}
                                alt={adv.title}
                                style={{
                                    width: '80px',
                                    height: '80px',
                                    objectFit: 'cover',
                                    borderRadius: '50%',
                                    minHeight: '80px',
                                    background: '#f5f5f5'
                                }}
                            />
                            <h4>{adv.title}</h4>
                            <p>{adv.description}</p>
                        </div>
                    ))}
                </div>
            </div>
        </section>
    );
};

export default Advantage;