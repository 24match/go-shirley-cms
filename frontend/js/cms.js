let cmsData = { config: {}, modules: {}, images: [], contentItems: [] };

document.addEventListener('DOMContentLoaded', async () => {
    await loadCMSData();
    applyCMSContent();
    
    // 当语言切换时，重新应用内容
    const originalSetLang = window.setLang;
    window.setLang = function(lang) {
        originalSetLang(lang);
        applyCMSContent();
    };
});

async function loadCMSData() {
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

function getLangSpecificField(obj, baseField) {
    const lang = getCurrentLang();
    const langField = (lang === 'zh' ? 'zh' : 'en') + baseField.charAt(0).toUpperCase() + baseField.slice(1);
    if (obj[langField]) return obj[langField];
    // 如果没有对应语言的内容，返回非多语言字段作为后备
    return obj[baseField] || '';
}

function applyCMSContent() {
    const { config, modules, images, contentItems } = cmsData;
    const moduleMap = modules;

    const factoryItems = contentItems.filter(item => item.section === 'factory');
    const advantageItems = contentItems.filter(item => item.section === 'advantage');

    loadBannerContent(moduleMap, config, images);
    loadAboutContent(moduleMap, config, images);
    loadStatsContent(config);
    loadEventContent(config, moduleMap);
    loadContactContent(config);
    loadBrandContent(config);
    loadProductsContent(images);
    loadFactoryContent(factoryItems);
    loadAdvantageContent(advantageItems);
}

function loadBannerContent(moduleMap, config, images) {
    const bannerModule = moduleMap['banner'] || config.banner || {};
    if (bannerModule.enabled === false) return;

    const banner = document.querySelector('.banner');
    const bannerText = document.querySelector('.banner-text');
    if (bannerText) {
        const h1 = bannerText.querySelector('h1');
        const p = bannerText.querySelector('p');
        const title = getLangSpecificField(bannerModule, 'title') || config.banner?.title;
        const subtitle = getLangSpecificField(bannerModule, 'subtitle') || config.banner?.subtitle;
        const content = getLangSpecificField(bannerModule, 'content') || config.banner?.content;
        
        if (h1 && title) {
            h1.innerHTML = title;
        }
        if (p && subtitle) {
            p.textContent = subtitle;
        }
    }

    if (bannerModule.imagePath && banner) {
        banner.style.background = `linear-gradient(135deg,rgba(10,92,173,0.9) 0%,rgba(6,58,117,0.85) 100%),url('/uploads/${bannerModule.imagePath}') center/cover no-repeat`;
    } else {
        const bannerImg = images.find(i => i.category === 'banner');
        if (bannerImg && banner) {
            banner.style.background = `linear-gradient(135deg,rgba(10,92,173,0.9) 0%,rgba(6,58,117,0.85) 100%),url('/uploads/${bannerImg.filename}') center/cover no-repeat`;
        }
    }
}

function loadAboutContent(moduleMap, config, images) {
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
        aboutLeftTitle.textContent = getCurrentLang() === 'zh' ? (aboutModule.zh_left_title || aboutModule.content) : (aboutModule.en_left_title || aboutModule.content);
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

    if (aboutModule.imagePath) {
        const imgEl = document.querySelector('#about .about-img img');
        if (imgEl) imgEl.src = `/uploads/${aboutModule.imagePath}`;
    } else {
        const aboutImg = images.find(i => i.category === 'about');
        if (aboutImg) {
            const imgEl = document.querySelector('#about .about-img img');
            if (imgEl) imgEl.src = `/uploads/${aboutImg.filename}`;
        }
    }
}

function loadStatsContent(config) {
    if (config.stats) {
        const stats = document.querySelectorAll('.stat-item');
        if (stats.length >= 4) {
            stats[0].querySelector('.number').textContent = config.stats.years || '15+';
            stats[1].querySelector('.number').textContent = config.stats.products || '200+';
            stats[2].querySelector('.number').textContent = config.stats.countries || '80+';
            stats[3].querySelector('.number').textContent = config.stats.clients || '1000+';
        }
    }
}

function loadEventContent(config, moduleMap) {
    const eventsModule = moduleMap['events'] || config.event || {};
    if (eventsModule.enabled === false) return;

    const eventCard = document.querySelector('#events .event-card');
    if (eventCard) {
        const eventName = getCurrentLang() === 'zh' ? (eventsModule.zhName || eventsModule.name) : (eventsModule.enName || eventsModule.name);
        if (eventName) {
            eventCard.querySelector('.event-content h3').textContent = eventName;
        }
        
        const eventDescription = getCurrentLang() === 'zh' ? (eventsModule.zhDescription || eventsModule.description) : (eventsModule.enDescription || eventsModule.description);
        if (eventDescription) {
            eventCard.querySelector('.event-content p').textContent = eventDescription;
        }
        
        const eventLocation = getCurrentLang() === 'zh' ? (eventsModule.zhLocation || eventsModule.location) : (eventsModule.enLocation || eventsModule.location);
        eventCard.querySelector('.event-booth').textContent = 'Booth: ' + (eventsModule.booth || 'F11');
        eventCard.querySelector('.event-date').textContent = '📅 ' + (eventsModule.start_date || 'June 17-19, 2026') + ' | ' + (eventLocation || 'Miami Beach Convention Center');
    }

    if (eventsModule.left_icon || eventsModule.icon) {
        const leftIcon = document.getElementById('eventLeftIcon');
        if (leftIcon) leftIcon.textContent = eventsModule.left_icon || eventsModule.icon || '🏥';
    }
    
    const leftTitle = getCurrentLang() === 'zh' ? (eventsModule.zhLeftTitle || eventsModule.left_title) : (eventsModule.enLeftTitle || eventsModule.left_title);
    if (leftTitle) {
        const leftTitleEl = document.getElementById('eventLeftTitle');
        if (leftTitleEl) leftTitleEl.textContent = leftTitle;
    }
    
    const leftSubtitle = getCurrentLang() === 'zh' ? (eventsModule.zhLeftSubtitle || eventsModule.left_subtitle) : (eventsModule.enLeftSubtitle || eventsModule.left_subtitle);
    if (leftSubtitle) {
        const leftSubtitleEl = document.getElementById('eventLeftSubtitle');
        if (leftSubtitleEl) leftSubtitleEl.textContent = leftSubtitle;
    }
}

function loadContactContent(config) {
    if (config.contact) {
        const emailEl = document.getElementById('contactEmail');
        const phoneEl = document.getElementById('contactPhone');
        const whatsappEl = document.getElementById('contactWhatsApp');
        const addressEl = document.getElementById('contactAddress');

        if (emailEl && config.contact.email) {
            emailEl.innerHTML = '<strong>Email:</strong><br>' + config.contact.email;
        }
        if (phoneEl && config.contact.phone) {
            phoneEl.innerHTML = '<strong>Tel:</strong><br>' + config.contact.phone;
        }
        if (whatsappEl && config.contact.whatsapp) {
            whatsappEl.innerHTML = '<strong>WhatsApp:</strong><br>' + config.contact.whatsapp;
        }
        if (addressEl && config.contact.address) {
            addressEl.innerHTML = '<strong>Address:</strong><br>' + config.contact.address;
        }
    }
}

function loadBrandContent(config) {
    if (config.brand) {
        const siteLogo = document.getElementById('siteLogo');
        const siteLogoSuffix = document.getElementById('siteLogoSuffix');
        const footerLogo = document.getElementById('footerLogo');
        const footerLogoSuffix = document.getElementById('footerLogoSuffix');

        if (siteLogo && config.brand.name) {
            const parts = config.brand.name.split(' ');
            siteLogo.childNodes[0].textContent = parts[0] || 'MEDICAL';
            if (siteLogoSuffix) {
                siteLogoSuffix.textContent = parts.slice(1).join(' ') || 'PRO';
                if (config.brand.suffix_color) {
                    siteLogoSuffix.style.color = config.brand.suffix_color;
                }
            }
        }

        if (footerLogo && config.brand.name) {
            const parts = config.brand.name.split(' ');
            footerLogo.childNodes[0].textContent = parts[0] || 'MEDICAL';
            if (footerLogoSuffix) {
                footerLogoSuffix.textContent = parts.slice(1).join(' ') || 'PRO';
                if (config.brand.suffix_color) {
                    footerLogoSuffix.style.color = config.brand.suffix_color;
                }
            }
        }
    }
}

function loadProductsContent(images) {
    const products = images.filter(i => i.category === 'products').sort((a, b) => a.sort_order - b.sort_order);
    if (products.length > 0) {
        const productGrid = document.querySelector('.product-grid');
        if (productGrid) {
            const defaultProducts = [
                {title:'Disposable Medical Supplies',desc:'Disposable masks, gloves, syringes, infusion sets, sterile dressings, safe and hygienic.'},
                {title:'Surgical Instruments',desc:'Surgical forceps, scissors, scalpels, surgical kits, high precision, corrosion resistance.'},
                {title:'Hospital Furniture',desc:'Hospital beds, nursing trolleys, medical cabinets, emergency equipment.'},
                {title:'Rehabilitation Equipment',desc:'Physical therapy equipment, rehabilitation aids for hospital and home care.'},
                {title:'Medical Protective Series',desc:'Medical protective clothing, isolation gowns, protective masks.'},
                {title:'OEM & ODM Custom Devices',desc:'Custom medical device size, logo, packaging customization.'}
            ];

            productGrid.innerHTML = products.map((img, i) => {
                const dp = defaultProducts[i] || {title:'Product',desc:'Product description'};
                return `
                    <div class="product-card">
                        <div class="product-img" style="background-image:url('/uploads/${img.filename}');background-size:cover;background-position:center;"></div>
                        <div class="product-info">
                            <h4>${img.description || dp.title}</h4>
                            <p>${img.long_description || dp.desc}</p>
                            <a href="#" class="learn-more" data-i18n="products.learnMore">Learn More →</a>
                        </div>
                    </div>
                `;
            }).join('');
        }
    }
}

function loadFactoryContent(factoryItems) {
    if (factoryItems.length > 0) {
        const factoryGrid = document.querySelector('.factory-grid');
        if (factoryGrid) {
            factoryGrid.innerHTML = factoryItems.map(item => {
                const imageHtml = item.image_path
                    ? `<img src="/uploads/${item.image_path}" style="width:100%;height:150px;object-fit:cover;border-radius:10px;" />`
                    : `<div class="factory-icon">${item.icon || '🏭'}</div>`;
                
                const title = getCurrentLang() === 'zh' ? (item.zhTitle || item.title) : (item.enTitle || item.title);
                const description = getCurrentLang() === 'zh' ? (item.zhDescription || item.description) : (item.enDescription || item.description);

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

function loadAdvantageContent(advantageItems) {
    if (advantageItems.length > 0) {
        const advGrid = document.querySelector('.adv-grid');
        if (advGrid) {
            advGrid.innerHTML = advantageItems.map(item => {
                const imageHtml = item.image_path
                    ? `<img src="/uploads/${item.image_path}" style="width:80px;height:80px;object-fit:cover;border-radius:50%;" />`
                    : `<div class="adv-icon">${item.icon || '⭐'}</div>`;
                
                const title = getCurrentLang() === 'zh' ? (item.zhTitle || item.title) : (item.enTitle || item.title);
                const description = getCurrentLang() === 'zh' ? (item.zhDescription || item.description) : (item.enDescription || item.description);

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