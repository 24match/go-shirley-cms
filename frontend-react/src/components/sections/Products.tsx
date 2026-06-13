/**
 * Products 组件
 * 产品展示网格，包含 6 个产品类别
 */

import React from 'react';
import { useLanguage } from '../../contexts/LanguageContext';
import { useCMSData } from '../../contexts/CMSContext';

interface Product {
    image: string;
    title: string;
    description: string;
}

const Products: React.FC = () => {
    const { t } = useLanguage();
    const { getModule, getImageByCategory } = useCMSData();

    const productsModule = getModule('products');

    const products: Product[] = [
        {
            image: '/uploads/1778807177495650000_favicon-32x32.png',
            title: 'Disposable Medical Supplies',
            description: 'Disposable masks, gloves, syringes, infusion sets, sterile dressings. Safe, hygienic, meeting international medical standards.'
        },
        {
            image: '/uploads/1778807177498171000_icon.png',
            title: 'Surgical Instruments',
            description: 'Stainless steel surgical forceps, scissors, scalpels, surgical kits. High precision, excellent corrosion resistance.'
        },
        {
            image: '/uploads/1778807177499079000_pdf.png',
            title: 'Hospital Furniture',
            description: 'Hospital beds, nursing trolleys, medical cabinets, emergency equipment. Durable and hospital-grade compliant.'
        },
        {
            image: '/uploads/1778807177501456000_qrcode-for-doocs.jpg',
            title: 'Rehabilitation Equipment',
            description: 'Physical therapy equipment, rehabilitation aids, nursing devices for hospital and home care settings.'
        },
        {
            image: '/uploads/1778806978916330000_IMG_4329.png',
            title: 'Medical Protective Series',
            description: 'Medical protective clothing, isolation gowns, protective masks, eye shields. Complete protection solutions.'
        },
        {
            image: '/uploads/1778806950150892000_IMG_4329.png',
            title: 'OEM & ODM Custom Devices',
            description: 'Custom medical device size, logo, packaging customization. Meeting diverse customer market demands.'
        }
    ];

    return (
        <section id="products" className="products-section">
            <div className="container">
                <div className="section-title">
                    <h2 data-i18n="products.title">{t('products.title')}</h2>
                    <p data-i18n="products.subtitle">{t('products.subtitle')}</p>
                </div>

                <div className="product-grid">
                    {products.map((product, index) => (
                        <div key={index} className="product-card">
                            <div 
                                className="product-img"
                                style={{
                                    backgroundImage: `url('${product.image}')`,
                                    backgroundSize: 'cover',
                                    backgroundPosition: 'center'
                                }}
                            />
                            <div className="product-info">
                                <h4>{product.title}</h4>
                                <p>{product.description}</p>
                                <a href="#" className="learn-more" data-i18n="products.learnMore">
                                    {t('products.learnMore')} →
                                </a>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </section>
    );
};

export default Products;