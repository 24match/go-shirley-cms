import { getCMSData, getLangSpecificField, getCurrentLangFromStorage } from '../services/cmsService.js';

export function loadEventContent() {
    const { config, modules } = getCMSData();
    const moduleMap = modules;
    
    const eventsModule = moduleMap['events'] || config.event || {};
    if (eventsModule.enabled === false) return;

    const eventCard = document.querySelector('#events .event-card');
    if (eventCard) {
        const eventName = getCurrentLangFromStorage() === 'zh' ? (eventsModule.zhName || eventsModule.name) : (eventsModule.enName || eventsModule.name);
        if (eventName) {
            eventCard.querySelector('.event-content h3').textContent = eventName;
        }
        
        const eventDescription = getCurrentLangFromStorage() === 'zh' ? (eventsModule.zhDescription || eventsModule.description) : (eventsModule.enDescription || eventsModule.description);
        if (eventDescription) {
            eventCard.querySelector('.event-content p').textContent = eventDescription;
        }
        
        const eventLocation = getCurrentLangFromStorage() === 'zh' ? (eventsModule.zhLocation || eventsModule.location) : (eventsModule.enLocation || eventsModule.location);
        eventCard.querySelector('.event-booth').textContent = 'Booth: ' + (eventsModule.booth || 'F11');
        eventCard.querySelector('.event-date').textContent = '📅 ' + (eventsModule.start_date || 'June 17-19, 2026') + ' | ' + (eventLocation || 'Miami Beach Convention Center');
    }

    if (eventsModule.left_icon || eventsModule.icon) {
        const leftIcon = document.getElementById('eventLeftIcon');
        if (leftIcon) leftIcon.textContent = eventsModule.left_icon || eventsModule.icon || '🏥';
    }
    
    const leftTitle = getCurrentLangFromStorage() === 'zh' ? (eventsModule.zhLeftTitle || eventsModule.left_title) : (eventsModule.enLeftTitle || eventsModule.left_title);
    if (leftTitle) {
        const leftTitleEl = document.getElementById('eventLeftTitle');
        if (leftTitleEl) leftTitleEl.textContent = leftTitle;
    }
    
    const leftSubtitle = getCurrentLangFromStorage() === 'zh' ? (eventsModule.zhLeftSubtitle || eventsModule.left_subtitle) : (eventsModule.enLeftSubtitle || eventsModule.left_subtitle);
    if (leftSubtitle) {
        const leftSubtitleEl = document.getElementById('eventLeftSubtitle');
        if (leftSubtitleEl) leftSubtitleEl.textContent = leftSubtitle;
    }
}