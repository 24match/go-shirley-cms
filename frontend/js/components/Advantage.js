import { getCMSData, getLangSpecificField, getCurrentLangFromStorage } from '../services/cmsService.js';

export function loadAdvantageContent() {
    const { contentItems } = getCMSData();
    
    const advantageItems = contentItems.filter(item => item.section === 'advantage');
    
    if (advantageItems.length > 0) {
        const advGrid = document.querySelector('.adv-grid');
        if (advGrid) {
            advGrid.innerHTML = advantageItems.map(item => {
                const imageHtml = item.imagePath
                    ? `<img src="/uploads/${item.imagePath}" style="width:80px;height:80px;object-fit:cover;border-radius:50%;" />`
                    : `<div class="adv-icon">${item.icon || '⭐'}</div>`;
                
                const title = getCurrentLangFromStorage() === 'zh' ? (item.zhTitle || item.title) : (item.enTitle || item.title);
                const description = getCurrentLangFromStorage() === 'zh' ? (item.zhDescription || item.description) : (item.enDescription || item.description);

                return `
                    <div class="adv-item">
                        ${imageHtml}
                        <h4>${title}</h4>
                        <p>${description}</p>
                    </div>
                `;
            }).join('');
        }
    }
}