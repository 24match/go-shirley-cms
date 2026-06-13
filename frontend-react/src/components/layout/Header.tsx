/**
 * Header 组件
 * 包含 Logo、导航菜单、语言切换和汉堡菜单
 */

import React, { useState } from 'react';
import { useLanguage } from '../../contexts/LanguageContext';
import { useCMSData } from '../../contexts/CMSContext';

const Header: React.FC = () => {
    const { currentLang, setLanguage, t } = useLanguage();
    const { getSiteSettings } = useCMSData();
    const [isMenuOpen, setIsMenuOpen] = useState(false);

    const siteSettings = getSiteSettings();
    const logoText = currentLang === 'zh' 
        ? siteSettings?.zhSiteLogo 
        : siteSettings?.enSiteLogo;
    const logoColor = siteSettings?.siteLogoColor || '#06a499';

    const handleLanguageToggle = async () => {
        await setLanguage(currentLang === 'en' ? 'zh' : 'en');
    };

    const handleNavClick = () => {
        setIsMenuOpen(false);
    };

    const navLinks = [
        { href: '#home', label: t('nav.home') },
        { href: '#about', label: t('nav.about') },
        { href: '#products', label: t('nav.products') },
        { href: '#factory', label: t('nav.factory') },
        { href: '#advantage', label: t('nav.advantage') },
        { href: '#contact', label: t('nav.contact') }
    ];

    return (
        <header className="header">
            <div className="container nav-wrap">
                <a href="#home" className="logo" id="siteLogo">
                    {logoText}
                    {currentLang === 'en' && logoText?.toUpperCase() === logoText && logoText.includes('MEDICAL') && (
                        <span id="siteLogoSuffix" style={{ color: logoColor }}>PRO</span>
                    )}
                </a>

                <nav className={`nav-menu ${isMenuOpen ? 'active' : ''}`} id="navMenu">
                    {navLinks.map((link) => (
                        <a
                            key={link.href}
                            href={link.href}
                            onClick={handleNavClick}
                        >
                            {link.label}
                        </a>
                    ))}
                </nav>

                <div className="lang-switcher" id="langSwitcher">
                    <button 
                        id="langBtn" 
                        className="lang-btn"
                        onClick={handleLanguageToggle}
                    >
                        {currentLang === 'zh' ? '中文' : 'EN'}
                    </button>
                </div>

                <div 
                    className={`hamburger ${isMenuOpen ? 'active' : ''}`}
                    id="hamburger"
                    onClick={() => setIsMenuOpen(!isMenuOpen)}
                >
                    <span></span>
                    <span></span>
                    <span></span>
                </div>
            </div>
        </header>
    );
};

export default Header;