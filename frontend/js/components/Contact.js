import { getCMSData } from '../services/cmsService.js';

export function loadContactContent() {
    const { config } = getCMSData();
    
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