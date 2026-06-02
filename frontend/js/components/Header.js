import { setLang, applyTranslations } from '../utils/i18n.js';
import { getCMSData, getCurrentLangFromStorage } from '../services/cmsService.js';

export function initHeader() {
    initLanguageSwitcher();
    initMobileMenu();
    initNavLinks();
    loadSiteLogo();
}

export function loadSiteLogo() {
    const { siteSettings } = getCMSData();
    if (!siteSettings) return;
    
    const lang = getCurrentLangFromStorage();
    const logoText = lang === 'zh' ? siteSettings.zhSiteLogo : siteSettings.enSiteLogo;
    const logoColor = siteSettings.siteLogoColor || '#06a499';
    
    // 更新顶部导航 Logo
    const siteLogo = document.getElementById('siteLogo');
    const siteLogoSuffix = document.getElementById('siteLogoSuffix');
    if (siteLogo && siteLogoSuffix && logoText) {
        // 检查是否包含 PRO 后缀（英文 Logo 的特殊处理）
        if (lang === 'en' && logoText && logoText.toUpperCase() === logoText) {
            // 英文大写 Logo，可能包含 PRO 后缀
            const parts = logoText.split('PRO');
            if (parts.length === 2) {
                siteLogo.textContent = parts[0] + ' ';
                siteLogoSuffix.textContent = 'PRO';
            } else {
                siteLogo.textContent = logoText;
                siteLogoSuffix.textContent = '';
            }
        } else {
            siteLogo.textContent = logoText;
            siteLogoSuffix.textContent = '';
        }
        siteLogoSuffix.style.color = logoColor;
    }
    
    // 更新页脚 Logo
    const footerLogo = document.getElementById('footerLogo');
    const footerLogoSuffix = document.getElementById('footerLogoSuffix');
    if (footerLogo && footerLogoSuffix && logoText) {
        if (lang === 'en' && logoText && logoText.toUpperCase() === logoText) {
            const parts = logoText.split('PRO');
            if (parts.length === 2) {
                footerLogo.textContent = parts[0] + ' ';
                footerLogoSuffix.textContent = 'PRO';
            } else {
                footerLogo.textContent = logoText;
                footerLogoSuffix.textContent = '';
            }
        } else {
            footerLogo.textContent = logoText;
            footerLogoSuffix.textContent = '';
        }
        footerLogoSuffix.style.color = logoColor;
    }
}

function initLanguageSwitcher() {
    const langBtn = document.getElementById('langBtn');
    if (langBtn) {
        const savedLang = localStorage.getItem('preferredLang') || 'en';
        langBtn.textContent = savedLang === 'zh' ? '中文' : 'EN';

        langBtn.addEventListener('click', function() {
            const currentLang = localStorage.getItem('preferredLang') || 'en';
            const newLang = currentLang === 'en' ? 'zh' : 'en';

            setLang(newLang);
            langBtn.textContent = newLang === 'zh' ? '中文' : 'EN';

            document.querySelectorAll('.nav-menu a').forEach(link => {
                link.classList.remove('active');
            });
        });
    }
}

function initMobileMenu() {
    const hamburger = document.getElementById('hamburger');
    if (hamburger) {
        hamburger.addEventListener('click', function() {
            document.getElementById('navMenu').classList.toggle('active');
        });
    }
}

function initNavLinks() {
    document.querySelectorAll('.nav-menu a').forEach(link => {
        link.addEventListener('click', () => {
            document.getElementById('navMenu').classList.remove('active');
        });
    });
}