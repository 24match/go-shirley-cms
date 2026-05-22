import { getCMSData, getCurrentLangFromStorage } from '../services/cmsService.js';

function parseExtraData(module) {
    if (!module) {
        console.warn('[Events] Module is null/undefined');
        return {};
    }
    
    if (!module.extraData) {
        console.log('[Events] No extraData found, checking if fields are already merged...');
        return {};
    }
    
    try {
        let parsed;
        if (typeof module.extraData === 'string') {
            parsed = JSON.parse(module.extraData);
        } else {
            parsed = module.extraData;
        }
        
        console.log('[Events] ✅ Successfully parsed extraData:', parsed);
        return parsed;
    } catch (e) {
        console.error('[Events] ❌ Failed to parse extraData:', e);
        console.error('[Events] Raw extraData value:', module.extraData);
        return {};
    }
}

function getEventField(module, extra, fieldName, fallbackFieldName) {
    // 策略1: 从模块顶层读取（新接口格式，mergeExtraData 已合并）
    if (module[fieldName] !== undefined && module[fieldName] !== '') {
        console.log(`[Events] 📖 Found "${fieldName}" in module top-level:`, module[fieldName]);
        return module[fieldName];
    }
    
    // 策略2: 从 extraData 解析后读取（旧接口格式）
    if (extra[fieldName] !== undefined && extra[fieldName] !== '') {
        console.log(`[Events] 📖 Found "${fieldName}" in extraData:`, extra[fieldName]);
        return extra[fieldName];
    }
    
    // 策略3: 使用备选字段名
    if (fallbackFieldName) {
        if (module[fallbackFieldName] !== undefined && module[fallbackFieldName] !== '') {
            console.log(`[Events] 📖 Found fallback "${fallbackFieldName}" in module:`, module[fallbackFieldName]);
            return module[fallbackFieldName];
        }
        if (extra[fallbackFieldName] !== undefined && extra[fallbackFieldName] !== '') {
            console.log(`[Events] 📖 Found fallback "${fallbackFieldName}" in extra:`, extra[fallbackFieldName]);
            return extra[fallbackFieldName];
        }
    }
    
    return '';
}

export function loadEventContent() {
    const { modules } = getCMSData();
    const eventsModule = modules['events'] || {};
    
    console.log('🔍 [Events] Loading event content...');
    console.log('🔍 [Events] Raw eventsModule keys:', Object.keys(eventsModule));
    
    if (eventsModule.enabled === false) {
        console.log('[Events] ⚠️ Module is disabled, skipping');
        return;
    }

    const extra = parseExtraData(eventsModule);
    
    const eventCard = document.querySelector('#events .event-card');
    if (!eventCard) {
        console.warn('[Events] Event card element not found');
        return;
    }
    
    const lang = getCurrentLangFromStorage();
    const isZh = lang === 'zh';
    
    // ========== 1. 展会名称 ==========
    const nameField = isZh ? 'zh_name' : 'en_name';
    const nameFallback = isZh ? 'zhName' : 'enName';
    const eventName = getEventField(eventsModule, extra, nameField, nameFallback);
    
    console.log(`[Events] 📝 Event Name (${lang}):`, eventName || '(empty)');
    
    if (eventName) {
        const h3 = eventCard.querySelector('.event-content h3');
        if (h3) h3.textContent = eventName;
    } else {
        console.warn('[Events] ⚠️ No event name found, keeping default');
    }
    
    // ========== 2. 展会描述 ==========
    const descField = isZh ? 'zhDescription' : 'enDescription';
    const eventDescription = eventsModule[descField] || '';
    
    console.log(`[Events] 📝 Event Description (${lang}):`, eventDescription || '(empty)');
    
    if (eventDescription) {
        const p = eventCard.querySelector('.event-content p');
        if (p) p.textContent = eventDescription;
    }
    
    // ========== 3. 展位号 ==========
    const booth = getEventField(eventsModule, extra, 'booth', null);
    console.log('[Events] 📝 Booth:', booth || '(empty)');
    
    const boothEl = eventCard.querySelector('.event-booth');
    if (boothEl) {
        boothEl.textContent = booth ? `Booth: ${booth}` : '';
    }
    
    // ========== 4. 日期信息 ==========
    const startDate = getEventField(eventsModule, extra, 'start_date', 'startDate');
    const endDate = getEventField(eventsModule, extra, 'end_date', 'endDate');
    
    console.log('[Events] 📅 Start Date:', startDate || '(empty)');
    console.log('[Events] 📅 End Date:', endDate || '(empty)');
    
    const locationField = isZh ? 'zh_location' : 'en_location';
    const locationFallback = isZh ? 'zhLocation' : 'enLocation';
    const eventLocation = getEventField(eventsModule, extra, locationField, locationFallback);
    
    console.log(`[Events] 📍 Location (${lang}):`, eventLocation || '(empty)');
    
    let dateText = '';
    let hasDateError = false;
    
    if (startDate && endDate) {
        dateText = `${startDate} - ${endDate}`;
        console.log('[Events] ✅ Showing date range:', dateText);
    } else if (startDate && !endDate) {
        dateText = startDate;
        console.log('[Events] ✅ Showing start date only:', dateText);
    } else if (!startDate && endDate) {
        hasDateError = true;
        dateText = '<span style="color:#dc3545;">⚠️ 请配置展会开始日期</span>';
        console.warn('[Events] ⚠️ End date without start date!');
    } else {
        dateText = 'June 17-19, 2026';
        console.log('[Events] ℹ️ Using default dates (none configured)');
    }
    
    const locationText = eventLocation || 'Miami Beach Convention Center';
    const dateEl = eventCard.querySelector('.event-date');
    
    if (dateEl) {
        if (hasDateError) {
            dateEl.innerHTML = `📅 ${dateText}`;
            
            const locationEl = document.createElement('div');
            locationEl.className = 'event-location-separate';
            locationEl.style.cssText = 'margin-top:8px;font-size:14px;opacity:0.9;';
            locationEl.innerHTML = `📍 ${locationText}`;
            
            const existingLocation = eventCard.querySelector('.event-location-separate');
            if (existingLocation) existingLocation.remove();
            
            dateEl.after(locationEl);
        } else {
            const existingLocation = eventCard.querySelector('.event-location-separate');
            if (existingLocation) existingLocation.remove();
            
            dateEl.innerHTML = `📅 ${dateText} | ${locationText}`;
        }
    }
    
    // ========== 5. 左侧图标 ==========
    const icon = getEventField(eventsModule, extra, 'icon', null);
    console.log('[Events] 🎨 Icon:', icon || '(empty)');
    
    if (icon) {
        const leftIcon = document.getElementById('eventLeftIcon');
        if (leftIcon) leftIcon.textContent = icon;
    }
    
    // ========== 6. 左侧标题 ==========
    const leftTitleField = isZh ? 'zhLeftTitle' : 'enLeftTitle';
    const leftTitle = getEventField(eventsModule, extra, leftTitleField, null);
    
    console.log(`[Events] 📝 Left Title (${lang}):`, leftTitle || '(empty)');
    
    if (leftTitle) {
        const leftTitleEl = document.getElementById('eventLeftTitle');
        if (leftTitleEl) leftTitleEl.textContent = leftTitle;
    }
    
    // ========== 7. 左侧副标题 ==========
    const leftSubtitleField = isZh ? 'zhLeftSubtitle' : 'enLeftSubtitle';
    const leftSubtitle = getEventField(eventsModule, extra, leftSubtitleField, null);
    
    console.log(`[Events] 📝 Left Subtitle (${lang}):`, leftSubtitle || '(empty)');
    
    if (leftSubtitle) {
        const leftSubtitleEl = document.getElementById('eventLeftSubtitle');
        if (leftSubtitleEl) leftSubtitleEl.textContent = leftSubtitle;
    }
    
    console.log('[Events] ✅ Event content loaded successfully!\n');
}
