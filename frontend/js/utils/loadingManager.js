class LoadingManager {
    constructor() {
        this.isLoading = true;
        this.loadStartTime = Date.now();
        this.minDisplayTime = 300;
        this.maxWaitTime = 5000;
        this.timeoutId = null;
        this.init();
    }

    init() {
        this.addLoadingStyles();
        this.hideDefaultContent();
        this.showLoadingState();
        this.startTimeout();
    }

    addLoadingStyles() {
        if (document.getElementById('loading-manager-styles')) return;

        const style = document.createElement('style');
        style.id = 'loading-manager-styles';
        style.textContent = `
            .cms-loading {
                opacity: 0 !important;
                visibility: hidden !important;
                transition: opacity 0.4s ease-out, visibility 0.4s ease-out;
            }
            
            .cms-loaded {
                opacity: 1 !important;
                visibility: visible !important;
                transition: opacity 0.4s ease-out, visibility 0.4s ease-out;
            }
            
            .loading-skeleton {
                position: relative;
                overflow: hidden;
                background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
                background-size: 200% 100%;
                animation: skeleton-loading 1.5s ease-in-out infinite;
                border-radius: 4px;
            }
            
            @keyframes skeleton-loading {
                0% { background-position: 200% 0; }
                100% { background-position: -200% 0; }
            }
            
            .loading-overlay {
                position: fixed;
                top: 0;
                left: 0;
                width: 100%;
                height: 100%;
                background: #ffffff;
                display: flex;
                justify-content: center;
                align-items: center;
                z-index: 9999;
                transition: opacity 0.4s ease-out, visibility 0.4s ease-out;
            }
            
            .loading-overlay.fade-out {
                opacity: 0;
                visibility: hidden;
            }
            
            .loading-spinner {
                width: 50px;
                height: 50px;
                border: 3px solid #f3f3f3;
                border-top: 3px solid #0a5cad;
                border-radius: 50%;
                animation: spin 1s linear infinite;
            }
            
            @keyframes spin {
                0% { transform: rotate(0deg); }
                100% { transform: rotate(360deg); }
            }
            
            .banner-skeleton,
            .stats-skeleton,
            .about-skeleton,
            .products-skeleton,
            .factory-skeleton,
            .advantage-skeleton,
            .events-skeleton,
            .contact-skeleton {
                min-height: 200px;
                margin: 40px auto;
            }
            
            .banner-skeleton {
                min-height: calc(100vh - 70px);
                margin: 70px auto 0;
            }
            
            .skeleton-text {
                height: 20px;
                margin-bottom: 15px;
                width: 100%;
            }
            
            .skeleton-text.short {
                width: 60%;
            }
            
            .skeleton-text.title {
                height: 40px;
                width: 80%;
                margin-bottom: 25px;
            }
            
            .skeleton-title {
                height: 35px;
                width: 50%;
                margin: 0 auto 20px;
            }
            
            .skeleton-grid {
                display: grid;
                grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
                gap: 30px;
                padding: 20px;
            }
            
            .skeleton-card {
                background: white;
                border-radius: 10px;
                padding: 25px;
                box-shadow: 0 2px 15px rgba(0,0,0,0.08);
            }
            
            .skeleton-card-img {
                height: 180px;
                width: 100%;
                margin-bottom: 20px;
                border-radius: 8px;
            }
            
            .skeleton-circle {
                width: 80px;
                height: 80px;
                border-radius: 50%;
                margin: 0 auto 20px;
            }
        `;
        document.head.appendChild(style);
    }

    hideDefaultContent() {
        const contentSections = [
            'header',
            '.banner',
            '.stats',
            '#about',
            '.products-section',
            '#factory',
            '.advantage-section',
            '#events',
            '.contact-section',
            'footer'
        ];

        contentSections.forEach(selector => {
            const elements = document.querySelectorAll(selector);
            elements.forEach(el => {
                el.classList.add('cms-loading');
                el.setAttribute('data-was-hidden', 'true');
            });
        });
    }

    showLoadingState() {
        const overlay = document.createElement('div');
        overlay.className = 'loading-overlay';
        overlay.id = 'loadingOverlay';
        overlay.innerHTML = '<div class="loading-spinner"></div>';
        document.body.appendChild(overlay);

        document.body.style.overflow = 'hidden';
    }

    startTimeout() {
        this.timeoutId = setTimeout(() => {
            console.warn('Loading timeout - showing content with defaults');
            this.forceShowContent();
        }, this.maxWaitTime);
    }

    async markAsLoaded() {
        if (this.timeoutId) {
            clearTimeout(this.timeoutId);
        }

        const elapsed = Date.now() - this.loadStartTime;
        const remainingTime = Math.max(0, this.minDisplayTime - elapsed);

        if (remainingTime > 0) {
            await new Promise(resolve => setTimeout(resolve, remainingTime));
        }

        this.showContentWithTransition();
    }

    showContentWithTransition() {
        const overlay = document.getElementById('loadingOverlay');
        
        if (overlay) {
            overlay.classList.add('fade-out');
            
            setTimeout(() => {
                if (overlay.parentNode) {
                    overlay.parentNode.removeChild(overlay);
                }
                document.body.style.overflow = '';
                this.revealContent();
            }, 400);
        } else {
            this.revealContent();
        }
    }

    revealContent() {
        const hiddenElements = document.querySelectorAll('[data-was-hidden="true"]');
        
        hiddenElements.forEach((el, index) => {
            setTimeout(() => {
                el.classList.remove('cms-loading');
                el.classList.add('cms-loaded');
                el.removeAttribute('data-was-hidden');
            }, index * 50);
        });

        setTimeout(() => {
            const loadedElements = document.querySelectorAll('.cms-loaded');
            loadedElements.forEach(el => {
                el.classList.remove('cms-loaded');
            });
        }, hiddenElements.length * 50 + 400);

        this.isLoading = false;
    }

    forceShowContent() {
        this.showContentWithTransition();
    }
}

const loadingManager = new LoadingManager();

export function waitForContentReady() {
    return loadingManager.markAsLoaded();
}

export function isContentLoading() {
    return loadingManager.isLoading;
}
