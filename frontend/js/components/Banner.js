import { getCMSData, getLangSpecificField } from '../services/cmsService.js';

export function loadBannerContent() {
    const { config, modules, images } = getCMSData();
    const moduleMap = modules;
    
    const bannerModule = moduleMap['banner'] || config.banner || {};
    if (bannerModule.enabled === false) return;

    const banner = document.querySelector('.banner');
    const bannerText = document.querySelector('.banner-text');
    if (bannerText) {
        const h1 = bannerText.querySelector('h1');
        const p = bannerText.querySelector('p');
        const title = getLangSpecificField(bannerModule, 'title') || config.banner?.title;
        const subtitle = getLangSpecificField(bannerModule, 'subtitle') || config.banner?.subtitle;

        if (h1 && title) {
            h1.innerHTML = title;
        }
        if (p && subtitle) {
            p.textContent = subtitle;
        }
    }

    if (banner) {
        if (bannerModule.imagePath) {
            banner.style.background = `linear-gradient(135deg,rgba(10,92,173,0.9) 0%,rgba(6,58,117,0.85) 100%),url('/uploads/${bannerModule.imagePath}') center/cover no-repeat`;
        } else {
            const bannerImg = images.find(i => i.category === 'banner');
            if (bannerImg) {
                banner.style.background = `linear-gradient(135deg,rgba(10,92,173,0.9) 0%,rgba(6,58,117,0.85) 100%),url('/uploads/${bannerImg.filename}') center/cover no-repeat`;
            }
        }
    }
}