document.addEventListener('DOMContentLoaded', function() {
    initLanguageSwitcher();
    initMobileMenu();
    initNavLinks();
});

function initLanguageSwitcher() {
    const langBtn = document.getElementById('langBtn');
    if (langBtn) {
        const savedLang = localStorage.getItem('preferredLang') || 'en';
        langBtn.textContent = savedLang === 'zh' ? '中文' : 'EN';

        langBtn.addEventListener('click', function() {
            const currentLang = localStorage.getItem('preferredLang') || 'en';
            const newLang = currentLang === 'en' ? 'zh' : 'en';

            setLang(newLang);
            langBtn.textContent = newLang === 'zh' ? '中文' : 'EN';

            document.querySelectorAll('.nav-menu a').forEach(link => {
                link.classList.remove('active');
            });
        });
    }
}

function initMobileMenu() {
    const hamburger = document.getElementById('hamburger');
    if (hamburger) {
        hamburger.addEventListener('click', function() {
            document.getElementById('navMenu').classList.toggle('active');
        });
    }
}

function initNavLinks() {
    document.querySelectorAll('.nav-menu a').forEach(link => {
        link.addEventListener('click', () => {
            document.getElementById('navMenu').classList.remove('active');
        });
    });
}