let cmsData = { config: {}, modules: {}, images: [], contentItems: [] };

export async function loadCMSData() {
    try {
        const [configRes, modulesRes, imagesRes, contentRes] = await Promise.all([
            fetch('/api/public/config'),
            fetch('/api/public/modules'),
            fetch('/api/public/images'),
            fetch('/api/public/content')
        ]);
        const configData = await configRes.json();
        const modulesData = await modulesRes.json();
        const imagesData = await imagesRes.json();
        const contentData = await contentRes.json();

        cmsData.config = configData.data || configData;
        cmsData.images = imagesData.data || imagesData;
        cmsData.contentItems = contentData.data || contentData;

        cmsData.modules = {};
        const modulesArray = Array.isArray(modulesData.data) ? modulesData.data : (Array.isArray(modulesData) ? modulesData : []);
        modulesArray.forEach(m => { cmsData.modules[m.moduleName] = m; });
    } catch(e) {
        console.log('CMS data loading failed, using defaults');
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