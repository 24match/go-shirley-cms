import { initI18n } from './utils/i18n.js';
import { initHeader, loadSiteLogo } from './components/Header.js';
import { loadBannerContent } from './components/Banner.js';
import { loadAboutContent } from './components/About.js';
import { loadStatsContent } from './components/Stats.js';
import { loadProductsContent } from './components/Products.js';
import { loadFactoryContent } from './components/Factory.js';
import { loadAdvantageContent } from './components/Advantage.js';
import { loadEventContent } from './components/Events.js';
import { loadContactContent } from './components/Contact.js';
import { loadBrandContent } from './components/Brand.js';
import { loadCMSData, getCMSData } from './services/cmsService.js';
import { waitForContentReady } from './utils/loadingManager.js';

document.addEventListener('DOMContentLoaded', () => {
    initI18n();
    initHeader(); // 初始化语言切换器和菜单
    
    loadCMSData().then(() => {
        // 数据加载完成后，加载 Logo 和其他内容
        loadSiteLogo();
        applyCMSContent();
        
        const originalSetLang = window.setLang;
        window.setLang = function(lang) {
            originalSetLang(lang);
            applyCMSContent();
            loadSiteLogo(); // 切换语言时更新 Logo
        };
        
        return waitForContentReady();
    }).catch(err => {
        console.log('CMS data loading failed, using defaults');
        waitForContentReady();
    });
    
    initializeLazyLoading();
});

function initializeLazyLoading() {
    const lazyImages = document.querySelectorAll('img[data-src]');
    const imageObserver = new IntersectionObserver((entries, observer) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                const img = entry.target;
                img.src = img.dataset.src;
                img.classList.remove('lazy');
                observer.unobserve(img);
            }
        });
    }, { rootMargin: '50px' });
    
    lazyImages.forEach(img => imageObserver.observe(img));
}

function applyCMSContent() {
    loadBannerContent();
    loadAboutContent();
    loadStatsContent();
    loadEventContent();
    loadContactContent();
    loadBrandContent();
    loadProductsContent();
    loadFactoryContent();
    loadAdvantageContent();
}