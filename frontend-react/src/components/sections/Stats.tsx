/**
 * Stats 组件
 * 统计数据展示（15+ 年经验、200+ 产品等）
 */

import React from 'react';
import { useLanguage } from '../../contexts/LanguageContext';

interface StatItem {
    number: string;
    labelKey: string;
}

const stats: StatItem[] = [
    { number: '15+', labelKey: 'stats.years' },
    { number: '200+', labelKey: 'stats.products' },
    { number: '80+', labelKey: 'stats.countries' },
    { number: '1000+', labelKey: 'stats.clients' }
];

const Stats: React.FC = () => {
    const { t } = useLanguage();

    return (
        <section className="stats">
            <div className="container">
                <div className="stats-grid">
                    {stats.map((stat) => (
                        <div key={stat.labelKey} className="stat-item">
                            <span className="number">{stat.number}</span>
                            <span className="label" data-i18n={stat.labelKey}>
                                {t(stat.labelKey)}
                            </span>
                        </div>
                    ))}
                </div>
            </div>
        </section>
    );
};

export default Stats;