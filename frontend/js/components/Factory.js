import { getCMSData, getLangSpecificField, getCurrentLangFromStorage } from '../services/cmsService.js';

export function loadFactoryContent() {
    const { contentItems } = getCMSData();
    
    const factoryItems = contentItems.filter(item => item.section === 'factory');
    
    if (factoryItems.length > 0) {
        const factoryGrid = document.querySelector('.factory-grid');
        if (factoryGrid) {
            factoryGrid.innerHTML = factoryItems.map(item => {
                const imageHtml = item.image_path
                    ? `<img src="/uploads/${item.image_path}" style="width:100%;height:150px;object-fit:cover;border-radius:10px;" />`
                    : `<div class="factory-icon">${item.icon || '🏭'}</div>`;
                
                const title = getCurrentLangFromStorage() === 'zh' ? (item.zhTitle || item.title) : (item.enTitle || item.title);
                const description = getCurrentLangFromStorage() === 'zh' ? (item.zhDescription || item.description) : (item.enDescription || item.description);

                return `
                    <div class="factory-item">
                        ${imageHtml}
                        <h4>${title}</h4>
                        <p>${description}</p>
                    </div>
                `;
            }).join('');
        }
    }
}