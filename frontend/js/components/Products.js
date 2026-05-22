import { getCMSData } from '../services/cmsService.js';

export function loadProductsContent() {
    const { images } = getCMSData();
    
    const products = images.filter(i => i.category === 'products').sort((a, b) => a.sortOrder - b.sortOrder);
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
                        <div class="product-img lazy-bg" data-src="/uploads/${img.filename}"></div>
                        <div class="product-info">
                            <h4>${img.description || dp.title}</h4>
                            <p>${img.longDescription || dp.desc}</p>
                            <a href="#" class="learn-more" data-i18n="products.learnMore">Learn More →</a>
                        </div>
                    </div>
                `;
            }).join('');
            
            initProductLazyLoading();
        }
    }
}

function initProductLazyLoading() {
    const lazyImages = document.querySelectorAll('.product-img.lazy-bg');
    const imageObserver = new IntersectionObserver((entries, observer) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                const img = entry.target;
                img.style.background = `url('${img.dataset.src}') center/cover no-repeat`;
                img.classList.remove('lazy-bg');
                observer.unobserve(img);
            }
        });
    }, { rootMargin: '50px' });
    
    lazyImages.forEach(img => imageObserver.observe(img));
}