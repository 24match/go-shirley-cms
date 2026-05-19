import { getCMSData } from '../services/cmsService.js';

export function loadBrandContent() {
    const { config } = getCMSData();
    
    if (config.brand) {
        const siteLogo = document.getElementById('siteLogo');
        const siteLogoSuffix = document.getElementById('siteLogoSuffix');
        const footerLogo = document.getElementById('footerLogo');
        const footerLogoSuffix = document.getElementById('footerLogoSuffix');

        if (siteLogo && config.brand.name) {
            const parts = config.brand.name.split(' ');
            siteLogo.childNodes[0].textContent = parts[0] || 'MEDICAL';
            if (siteLogoSuffix) {
                siteLogoSuffix.textContent = parts.slice(1).join(' ') || 'PRO';
                if (config.brand.suffix_color) {
                    siteLogoSuffix.style.color = config.brand.suffix_color;
                }
            }
        }

        if (footerLogo && config.brand.name) {
            const parts = config.brand.name.split(' ');
            footerLogo.childNodes[0].textContent = parts[0] || 'MEDICAL';
            if (footerLogoSuffix) {
                footerLogoSuffix.textContent = parts.slice(1).join(' ') || 'PRO';
                if (config.brand.suffix_color) {
                    footerLogoSuffix.style.color = config.brand.suffix_color;
                }
            }
        }
    }
}