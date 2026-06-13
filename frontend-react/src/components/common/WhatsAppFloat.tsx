/**
 * WhatsApp 浮动按钮组件
 */

import React from 'react';
import { useCMSData } from '../../contexts/CMSContext';

const WhatsAppFloat: React.FC = () => {
    const { getSiteSettings } = useCMSData();
    const siteSettings = getSiteSettings();
    
    // 移除 WhatsApp 号码中的非数字字符
    const whatsappNumber = siteSettings?.contactWhatsapp?.replace(/[^\d]/g, '') || '8613800000000';
    const whatsappUrl = `https://wa.me/${whatsappNumber}`;

    return (
        <a 
            href={whatsappUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="whatsapp-float"
            aria-label="Contact us on WhatsApp"
        >
            W
        </a>
    );
};

export default WhatsAppFloat;