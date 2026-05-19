import { initI18n } from './utils/i18n.js';
import { initHeader } from './components/Header.js';
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

document.addEventListener('DOMContentLoaded', () => {
    initI18n();
    initHeader();
    
    loadCMSData().then(() => {
        applyCMSContent();
        
        const originalSetLang = window.setLang;
        window.setLang = function(lang) {
            originalSetLang(lang);
            applyCMSContent();
        };
    }).catch(err => {
        console.log('CMS data loading failed, using defaults');
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