/**
 * Events 组件
 * 展会活动信息（WHX Miami 2026）
 */

import React from 'react';
import { useLanguage } from '../../contexts/LanguageContext';
import { useCMSData } from '../../contexts/CMSContext';

const Events: React.FC = () => {
    const { t } = useLanguage();
    const { getModule } = useCMSData();

    const eventsModule = getModule('events');

    return (
        <section id="events">
            <div className="container">
                <div className="section-title">
                    <h2 data-i18n="events.title">{t('events.title')}</h2>
                    <p data-i18n="events.subtitle">{t('events.subtitle')}</p>
                </div>

                <div className="event-card">
                    <div className="event-img">
                        <div>
                            <div style={{ fontSize: '48px', marginBottom: '15px' }} id="eventLeftIcon">
                                🏥
                            </div>
                            <h3 style={{ fontSize: '24px', fontWeight: 700 }} id="eventLeftTitle">
                                WHX Miami 2026
                            </h3>
                            <p style={{ opacity: 0.9 }} id="eventLeftSubtitle">
                                World Health Expo
                            </p>
                        </div>
                    </div>
                    <div className="event-content">
                        <span className="event-tag" data-i18n="events.tag">{t('events.tag')}</span>
                        <h3>World Health Exhibition Miami 2026</h3>
                        <p>
                            Join us at the premier global healthcare exhibition. Discover our latest medical innovations and connect with industry leaders.
                        </p>
                        <div className="event-booth" data-i18n="events.booth">
                            {t('events.booth')} A123
                        </div>
                        <div className="event-date">
                            <span data-i18n="events.date">{t('events.date')}</span> June 17-19, 2026 | Miami Beach Convention Center, Florida, USA
                        </div>
                    </div>
                </div>
            </div>
        </section>
    );
};

export default Events;