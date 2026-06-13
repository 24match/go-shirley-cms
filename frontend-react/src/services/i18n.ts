/**
 * 国际化服务
 * 管理多语言翻译和语言切换
 */

import { cmsApi } from './api';
import type { Translations } from '../types/cms';

/**
 * 默认翻译数据（离线备用）
 */
export const defaultTranslations: Translations = {
    en: {
        'nav.home': 'Home',
        'nav.about': 'About Us',
        'nav.products': 'Products',
        'nav.factory': 'Factory',
        'nav.advantage': 'Advantages',
        'nav.contact': 'Contact',
        'banner.title': 'Professional Medical Device <br/>Manufacturer & Exporter',
        'banner.subtitle': 'CE, FDA, ISO Certified Products. We specialize in high-quality hospital equipment, surgical instruments, and disposable medical supplies with comprehensive OEM/ODM services.',
        'banner.viewProducts': 'View Products',
        'banner.getQuote': 'Get Free Quote',
        'stats.years': 'Years of Excellence',
        'stats.products': 'Products Categories',
        'stats.countries': 'Countries Served',
        'stats.clients': 'Happy Clients',
        'about.title': 'About Our Company',
        'about.subtitle': 'Your Trusted Partner in Medical Device Manufacturing',
        'about.leftTitle': 'Leading Manufacturer & Exporter from China',
        'about.ce': 'CE Certified',
        'about.fda': 'FDA Registered',
        'about.iso': 'ISO13485',
        'about.sgs': 'SGS Audited',
        'products.title': 'Our Main Products',
        'products.subtitle': 'High-quality medical devices for global healthcare',
        'products.learnMore': 'Learn More',
        'factory.title': 'Factory Strength',
        'factory.subtitle': 'Advanced manufacturing capabilities',
        'advantage.title': 'Our Advantages',
        'advantage.subtitle': 'Why choose us as your medical device partner',
        'events.title': 'Upcoming Events',
        'events.subtitle': 'Join us at international medical exhibitions',
        'events.tag': 'Upcoming Event',
        'events.booth': 'Booth:',
        'events.date': '📅',
        'contact.title': 'Contact Us',
        'contact.subtitle': 'Get the latest medical device wholesale prices and customized solutions',
        'contact.getInTouch': 'Get In Touch',
        'contact.email': 'Email:',
        'contact.whatsapp': 'WhatsApp:',
        'contact.tel': 'Tel:',
        'contact.address': 'Address:',
        'contact.hours': 'Business Hours:',
        'contact.online': '24h Online Service Available',
        'contact.form.name': 'Your Name *',
        'contact.form.email': 'Your Email *',
        'contact.form.company': 'Country / Company Name',
        'contact.form.inquiry': 'Your Inquiry (Products, Quantity, Custom Requirements...)',
        'contact.form.submit': 'Send Inquiry',
        'footer.quickLinks': 'Quick Links',
        'footer.products': 'Products',
        'footer.services': 'Services',
        'footer.oem': 'OEM / ODM Custom',
        'footer.wholesale': 'Global Wholesale',
        'footer.shipping': 'International Shipping',
        'footer.support': '24h Support',
        'footer.copyright': '© 2025 Medical Device Manufacturer. All Rights Reserved.',
        'lang.chinese': '中文',
        'lang.english': 'English'
    },
    zh: {
        'nav.home': '首页',
        'nav.about': '关于我们',
        'nav.products': '产品中心',
        'nav.factory': '工厂实力',
        'nav.advantage': '核心优势',
        'nav.contact': '联系我们',
        'banner.title': '专业医疗器械 <br/>制造商与出口商',
        'banner.subtitle': 'CE、FDA、ISO 认证产品。专注于高品质医院设备、手术器械和一次性医疗用品，提供全方位 OEM/ODM 服务。',
        'banner.viewProducts': '查看产品',
        'banner.getQuote': '获取报价',
        'stats.years': '年行业经验',
        'stats.products': '产品品类',
        'stats.countries': '服务国家',
        'stats.clients': '客户数量',
        'about.title': '关于我们',
        'about.subtitle': '您值得信赖的医疗器械制造合作伙伴',
        'about.leftTitle': '中国领先制造商与出口商',
        'about.ce': 'CE 认证',
        'about.fda': 'FDA 注册',
        'about.iso': 'ISO13485 认证',
        'about.sgs': 'SGS 审核',
        'products.title': '主要产品',
        'products.subtitle': '面向全球医疗健康的高品质医疗器械',
        'products.learnMore': '了解更多',
        'factory.title': '工厂实力',
        'factory.subtitle': '先进的制造能力',
        'advantage.title': '核心优势',
        'advantage.subtitle': '为什么选择我们作为您的医疗器械合作伙伴',
        'events.title': '即将举办的展会',
        'events.subtitle': '诚邀您参加国际医疗展会',
        'events.tag': '即将开展',
        'events.booth': '展位:',
        'events.date': '📅',
        'contact.title': '联系我们',
        'contact.subtitle': '获取最新医疗器械批发价格和定制解决方案',
        'contact.getInTouch': '联系我们',
        'contact.email': '邮箱:',
        'contact.whatsapp': 'WhatsApp:',
        'contact.tel': '电话:',
        'contact.address': '地址:',
        'contact.hours': '营业时间:',
        'contact.online': '24 小时在线服务',
        'contact.form.name': '您的姓名 *',
        'contact.form.email': '您的邮箱 *',
        'contact.form.company': '国家/公司名称',
        'contact.form.inquiry': '您的咨询（产品、数量、定制需求...）',
        'contact.form.submit': '提交咨询',
        'footer.quickLinks': '快速链接',
        'footer.products': '产品',
        'footer.services': '服务',
        'footer.oem': 'OEM/ODM 定制',
        'footer.wholesale': '全球批发',
        'footer.shipping': '国际运输',
        'footer.support': '24 小时支持',
        'footer.copyright': '© 2025 医疗器械制造商 版权所有。',
        'lang.chinese': '中文',
        'lang.english': 'English'
    }
};

/**
 * 动态翻译数据（从 API 加载）
 */
let dynamicTranslations: Translations = { en: {}, zh: {} };

/**
 * 当前语言
 */
let currentLang: 'en' | 'zh' = 'en';

/**
 * 获取当前语言
 */
export function getCurrentLang(): 'en' | 'zh' {
    return currentLang;
}

/**
 * 设置语言
 */
export function setLanguage(lang: 'en' | 'zh'): void {
    currentLang = lang;
    localStorage.setItem('preferredLang', lang);
    document.documentElement.lang = lang;
}

/**
 * 初始化语言设置
 */
export function initLanguage(): 'en' | 'zh' {
    const savedLang = localStorage.getItem('preferredLang') as 'en' | 'zh' || 'en';
    currentLang = savedLang;
    document.documentElement.lang = savedLang;
    return savedLang;
}

/**
 * 翻译单个键
 */
export function translate(key: string): string {
    // 优先使用动态翻译
    const dynamicValue = dynamicTranslations[currentLang][key];
    if (dynamicValue !== undefined && dynamicValue !== null && dynamicValue !== '') {
        return dynamicValue;
    }
    // 回退到默认翻译
    return defaultTranslations[currentLang][key] || key;
}

/**
 * 批量翻译
 */
export function translateAll(keys: string[]): Record<string, string> {
    const result: Record<string, string> = {};
    keys.forEach(key => {
        result[key] = translate(key);
    });
    return result;
}

/**
 * 从 API 加载翻译
 */
export async function loadTranslationsFromAPI(): Promise<void> {
    try {
        const translations = await cmsApi.getTranslations();
        if (translations.en) {
            dynamicTranslations.en = { ...translations.en };
        }
        if (translations.zh) {
            dynamicTranslations.zh = { ...translations.zh };
        }
        console.log('✅ [i18n] Dynamic translations loaded from API');
    } catch (error) {
        console.warn('⚠️ [i18n] Failed to load translations from API, using defaults', error);
    }
}

/**
 * 根据语言获取字段值（用于 CMS 数据中的多语言字段）
 */
export function getLangField<T extends Record<string, any>>(
    obj: T,
    baseField: string
): string {
    const langPrefix = currentLang === 'zh' ? 'zh' : 'en';
    const langField = `${langPrefix}${baseField.charAt(0).toUpperCase()}${baseField.slice(1)}`;
    
    if (obj[langField]) {
        return String(obj[langField]);
    }
    if (obj[baseField]) {
        return String(obj[baseField]);
    }
    return '';
}

/**
 * 更新站点标题
 */
export async function updateSiteTitle(siteSettings?: { zhSiteTitle?: string; enSiteTitle?: string }): Promise<void> {
    if (!siteSettings) {
        try {
            const settings = await cmsApi.getSiteSettings();
            siteSettings = settings;
        } catch (error) {
            console.warn('⚠️ [Site Title] Failed to load site settings');
            return;
        }
    }
    
    const siteTitle = currentLang === 'zh' ? siteSettings.zhSiteTitle : siteSettings.enSiteTitle;
    if (siteTitle) {
        document.title = siteTitle;
    }
}