/**
 * Factory 组件
 * 工厂实力展示（无尘车间、QC 系统、先进设备、大容量）
 */

import React from 'react';
import { useLanguage } from '../../contexts/LanguageContext';
import { useCMSData } from '../../contexts/CMSContext';
import LazyImage from '../common/LazyImage';

interface FactoryItem {
    image: string;
    title: string;
    description: string;
}

const Factory: React.FC = () => {
    const { t } = useLanguage();
    const { getModule } = useCMSData();

    const factoryModule = getModule('factory');

    const factoryItems: FactoryItem[] = [
        {
            image: '/uploads/1779074837_pdf.png',
            title: 'Dust-free Workshop',
            description: '100,000 grade sterile dust-free production workshop, fully compliant with international medical standards.'
        },
        {
            image: '/uploads/1779075086_favicon-32x32.png',
            title: 'Strict QC System',
            description: 'Full-process quality inspection from raw materials to finished products, 100% quality testing before delivery.'
        },
        {
            image: '/uploads/1778807729666742000_starcharts.svg',
            title: 'Advanced Equipment',
            description: 'Imported automated production equipment ensuring high efficiency and consistent product quality.'
        },
        {
            image: '/uploads/1778806978916330000_IMG_4329.png',
            title: 'Large Capacity',
            description: 'Large daily output with sufficient inventory, supporting bulk orders and fast delivery worldwide.'
        }
    ];

    return (
        <section id="factory">
            <div className="container">
                <div className="section-title">
                    <h2 data-i18n="factory.title">{t('factory.title')}</h2>
                    <p data-i18n="factory.subtitle">{t('factory.subtitle')}</p>
                </div>

                <div className="factory-grid">
                    {factoryItems.map((item, index) => (
                        <div key={index} className="factory-item">
                            <LazyImage
                                dataSrc={item.image}
                                alt={item.title}
                                style={{
                                    width: '100%',
                                    height: '150px',
                                    objectFit: 'cover',
                                    borderRadius: '10px',
                                    minHeight: '150px',
                                    background: '#f5f5f5'
                                }}
                            />
                            <h4>{item.title}</h4>
                            <p>{item.description}</p>
                        </div>
                    ))}
                </div>
            </div>
        </section>
    );
};

export default Factory;