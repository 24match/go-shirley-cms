import { getCMSData, getLangSpecificField, getCurrentLangFromStorage } from '../services/cmsService.js';

export function loadAboutContent() {
    const { config, modules, images } = getCMSData();
    const moduleMap = modules;
    
    const aboutModule = moduleMap['about'] || config.about || {};
    if (aboutModule.enabled === false) return;

    const aboutTitle = document.getElementById('aboutTitle');
    const title = getLangSpecificField(aboutModule, 'title') || config.about?.title;
    if (aboutTitle && title) {
        aboutTitle.textContent = title;
    }

    const aboutSubtitle = document.getElementById('aboutSubtitle');
    const subtitle = getLangSpecificField(aboutModule, 'subtitle') || config.about?.subtitle;
    if (aboutSubtitle && subtitle) {
        aboutSubtitle.textContent = subtitle;
    }

    const aboutLeftTitle = document.getElementById('aboutLeftTitle');
    const leftTitle = aboutModule.zh_left_title || aboutModule.en_left_title || aboutModule.content || config.about?.left_title;
    if (aboutLeftTitle) {
        aboutLeftTitle.textContent = getCurrentLangFromStorage() === 'zh' ? (aboutModule.zh_left_title || aboutModule.content) : (aboutModule.en_left_title || aboutModule.content);
    }

    const contentToShow = getLangSpecificField(aboutModule, 'content') || config.about?.content;
    if (contentToShow) {
        const container = document.getElementById('aboutContentContainer');
        if (container) {
            container.innerHTML = '';
            const paragraphs = contentToShow.split(/\n\s*\n/).filter(p => p.trim());
            paragraphs.forEach(text => {
                const p = document.createElement('p');
                p.textContent = text.trim();
                container.appendChild(p);
            });
        }
    }

    const imgEl = document.querySelector('#about .about-img img');
    if (imgEl) {
        if (aboutModule.imagePath) {
            imgEl.src = `/uploads/${aboutModule.imagePath}`;
        } else {
            const aboutImg = images.find(i => i.category === 'about');
            if (aboutImg) {
                imgEl.src = `/uploads/${aboutImg.filename}`;
            }
        }
    }
}