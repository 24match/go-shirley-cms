import { getCMSData } from '../services/cmsService.js';

export function loadContactContent() {
    const { modules } = getCMSData();
    
    const contactModule = modules['contact'];
    console.log('[Contact] Module config:', contactModule);
    
    if (contactModule && contactModule.enabled !== false) {
        // 从 extraData 中解析联系信息（支持两种格式）
        let contactData = {};
        
        // 方式 1：后端 mergeExtraData 已经解析出来的字段
        if (contactModule.email) contactData.email = contactModule.email;
        if (contactModule.phone) contactData.phone = contactModule.phone;
        if (contactModule.whatsapp) contactData.whatsapp = contactModule.whatsapp;
        if (contactModule.address) contactData.address = contactModule.address;
        
        console.log('[Contact] Parsed data from top-level fields:', contactData);
        
        // 方式 2：从 extraData JSON 字符串中解析
        if (contactModule.extraData) {
            try {
                const extraData = JSON.parse(contactModule.extraData);
                if (extraData.email && !contactData.email) contactData.email = extraData.email;
                if (extraData.phone && !contactData.phone) contactData.phone = extraData.phone;
                if (extraData.whatsapp && !contactData.whatsapp) contactData.whatsapp = extraData.whatsapp;
                if (extraData.address && !contactData.address) contactData.address = extraData.address;
                console.log('[Contact] Parsed data from extraData:', extraData);
            } catch (e) {
                console.warn('[Contact] Failed to parse extraData:', e);
            }
        }
        
        console.log('[Contact] Final contact data:', contactData);
        
        const emailEl = document.getElementById('contactEmail');
        const phoneEl = document.getElementById('contactPhone');
        const whatsappEl = document.getElementById('contactWhatsApp');
        const addressEl = document.getElementById('contactAddress');

        // 更新 Email（有数据则显示配置的值，否则显示默认值）
        if (emailEl) {
            if (contactData.email) {
                emailEl.innerHTML = '<strong data-i18n="contact.email">Email:</strong><br>' + contactData.email;
            }
        }
        
        // 更新 Phone（有数据则显示配置的值，否则显示默认值）
        if (phoneEl) {
            if (contactData.phone) {
                phoneEl.innerHTML = '<strong data-i18n="contact.tel">Tel:</strong><br>' + contactData.phone;
            }
        }
        
        // 更新 WhatsApp（有数据则显示配置的值，否则显示默认值）
        if (whatsappEl) {
            if (contactData.whatsapp) {
                whatsappEl.innerHTML = '<strong data-i18n="contact.whatsapp">WhatsApp:</strong><br>' + contactData.whatsapp;
            }
        }
        
        // 更新 Address（有数据则显示配置的值，否则显示默认值）
        if (addressEl) {
            if (contactData.address) {
                addressEl.innerHTML = '<strong data-i18n="contact.address">Address:</strong><br>' + contactData.address;
            }
        }
    }
    
    // 初始化表单提交处理
    initContactForm();
}

// 初始化联系表单提交处理
function initContactForm() {
    const contactSection = document.getElementById('contact');
    if (!contactSection) return;
    
    const form = contactSection.querySelector('form');
    if (!form) return;
    
    form.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const inputs = form.querySelectorAll('input, textarea');
        const nameInput = inputs[0];
        const emailInput = inputs[1];
        const companyInput = inputs[2];
        const inquiryInput = inputs[3];
        
        const submitBtn = form.querySelector('button[type="submit"]');
        const originalText = submitBtn.textContent;
        
        // 验证必填字段
        if (!nameInput.value.trim() || !emailInput.value.trim()) {
            alert('Please fill in required fields (Name and Email)');
            return;
        }
        
        // 禁用按钮防止重复提交
        submitBtn.disabled = true;
        submitBtn.textContent = 'Submitting...';
        
        try {
            const response = await fetch('/api/contact/submit', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    name: nameInput.value.trim(),
                    email: emailInput.value.trim(),
                    company: companyInput.value.trim(),
                    inquiry: inquiryInput.value.trim()
                })
            });
            
            const result = await response.json();
            
            if (response.ok && result.code === 200) {
                alert('Thank you for your inquiry! We will contact you soon.');
                form.reset();
            } else {
                alert(result.message || 'Submission failed. Please try again.');
            }
        } catch (error) {
            console.error('Form submission error:', error);
            alert('Submission failed. Please check your network connection.');
        } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = originalText;
        }
    });
}