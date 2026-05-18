const defaultTranslations = {
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
        'banner.subtitle': 'CE、FDA、ISO认证产品。专注于高品质医院设备、手术器械和一次性医疗用品，提供全方位OEM/ODM服务。',
        'banner.viewProducts': '查看产品',
        'banner.getQuote': '获取报价',
        'stats.years': '年行业经验',
        'stats.products': '产品品类',
        'stats.countries': '服务国家',
        'stats.clients': '客户数量',
        'about.title': '关于我们',
        'about.subtitle': '您值得信赖的医疗器械制造合作伙伴',
        'about.leftTitle': '中国领先制造商与出口商',
        'about.ce': 'CE认证',
        'about.fda': 'FDA注册',
        'about.iso': 'ISO13485认证',
        'about.sgs': 'SGS审核',
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
        'contact.online': '24小时在线服务',
        'contact.form.name': '您的姓名 *',
        'contact.form.email': '您的邮箱 *',
        'contact.form.company': '国家/公司名称',
        'contact.form.inquiry': '您的咨询（产品、数量、定制需求...）',
        'contact.form.submit': '提交咨询',
        'footer.quickLinks': '快速链接',
        'footer.products': '产品',
        'footer.services': '服务',
        'footer.oem': 'OEM/ODM定制',
        'footer.wholesale': '全球批发',
        'footer.shipping': '国际运输',
        'footer.support': '24小时支持',
        'footer.copyright': '© 2025 医疗器械制造商 版权所有。',
        'lang.chinese': '中文',
        'lang.english': 'English'
    }
};

let currentLang = 'en';
let dynamicTranslations = { en: {}, zh: {} };
let translationsLoaded = false;

function getCurrentLang() {
    return currentLang;
}

function setLang(lang) {
    if (defaultTranslations[lang]) {
        currentLang = lang;
        localStorage.setItem('preferredLang', lang);
        document.documentElement.lang = lang;
        applyTranslations();
    }
}

function translate(key) {
    const dynamic = dynamicTranslations[currentLang][key];
    if (dynamic !== undefined && dynamic !== null && dynamic !== '') {
        return dynamic;
    }
    return defaultTranslations[currentLang][key] || key;
}

function applyTranslations() {
    document.querySelectorAll('[data-i18n]').forEach(element => {
        const key = element.getAttribute('data-i18n');
        const translation = translate(key);
        
        if (element.tagName === 'INPUT' || element.tagName === 'TEXTAREA') {
            element.placeholder = translation;
        } else {
            if (element.getAttribute('data-i18n-html')) {
                element.innerHTML = translation;
            } else {
                element.textContent = translation;
            }
        }
    });
    
    const langBtn = document.getElementById('langBtn');
    if (langBtn) {
        langBtn.textContent = currentLang === 'zh' ? '中文' : 'EN';
    }
}

async function loadTranslationsFromAPI() {
    try {
        const res = await fetch('/api/public/lang');
        if (res.ok) {
            const data = await res.json();
            if (data.en) dynamicTranslations.en = { ...data.en };
            if (data.zh) dynamicTranslations.zh = { ...data.zh };
            translationsLoaded = true;
            console.log('✅ [i18n] Dynamic translations loaded from API');
            applyTranslations();
        }
    } catch (err) {
        console.warn('⚠️ [i18n] Failed to load translations from API, using defaults');
    }
}

function initI18n() {
    const savedLang = localStorage.getItem('preferredLang') || 'en';
    setLang(savedLang);
    loadTranslationsFromAPI();
}

document.addEventListener('DOMContentLoaded', initI18n);
