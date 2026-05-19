import { getCMSData } from '../services/cmsService.js';

export function loadStatsContent() {
    const { config } = getCMSData();
    
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