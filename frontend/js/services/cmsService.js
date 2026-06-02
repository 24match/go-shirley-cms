let cmsData = { config: {}, modules: {}, images: [], contentItems: [], siteSettings: {} };

export async function loadCMSData() {
    try {
        const [configRes, modulesRes, imagesRes, contentRes, settingsRes] = await Promise.all([
            fetch('/api/public/config'),
            fetch('/api/public/modules'),
            fetch('/api/public/images'),
            fetch('/api/public/content'),
            fetch('/api/public/site-settings')
        ]);
        const configData = await configRes.json();
        const modulesData = await modulesRes.json();
        const imagesData = await imagesRes.json();
        const contentData = await contentRes.json();
        const settingsData = await settingsRes.json();

        // 将 PageConfig 数组转换为以 pageName 为键的对象映射
        const configArray = configData.data || configData;
        if (Array.isArray(configArray)) {
            cmsData.config = {};
            configArray.forEach(cfg => {
                if (cfg.pageName) {
                    try {
                        cmsData.config[cfg.pageName] = JSON.parse(cfg.configData);
                    } catch (e) {
                        cmsData.config[cfg.pageName] = {};
                    }
                }
            });
        } else {
            cmsData.config = configArray;
        }
        cmsData.images = imagesData.data || imagesData;
        cmsData.contentItems = contentData.data || contentData;
        cmsData.siteSettings = settingsData.data || settingsData;

        cmsData.modules = {};
        const modulesArray = Array.isArray(modulesData.data) ? modulesData.data : (Array.isArray(modulesData) ? modulesData : []);
        modulesArray.forEach(m => { cmsData.modules[m.moduleName] = m; });
    } catch(e) {
        console.log('CMS data loading failed, using defaults');
        // 使用默认 Logo 配置
        cmsData.siteSettings = {
            zhSiteLogo: '医疗',
            enSiteLogo: 'MEDICAL',
            siteLogoColor: '#06a499'
        };
    }
}

export function getCMSData() {
    return cmsData;
}

export function getLangSpecificField(obj, baseField) {
    const lang = localStorage.getItem('preferredLang') || 'en';
    const langField = (lang === 'zh' ? 'zh' : 'en') + baseField.charAt(0).toUpperCase() + baseField.slice(1);
    if (obj[langField]) return obj[langField];
    return obj[baseField] || '';
}

export function getCurrentLangFromStorage() {
    return localStorage.getItem('preferredLang') || 'en';
}