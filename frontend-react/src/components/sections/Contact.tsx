/**
 * Contact 组件
 * 联系方式和表单
 */

import React, { useState, FormEvent } from 'react';
import { useLanguage } from '../../contexts/LanguageContext';
import { useCMSData } from '../../contexts/CMSContext';
import { submitContactForm } from '../../services/api';

interface FormData {
    name: string;
    email: string;
    company: string;
    inquiry: string;
}

interface FormErrors {
    name?: string;
    email?: string;
    inquiry?: string;
}

const Contact: React.FC = () => {
    const { t } = useLanguage();
    const { getSiteSettings } = useCMSData();
    const siteSettings = getSiteSettings();

    const [formData, setFormData] = useState<FormData>({
        name: '',
        email: '',
        company: '',
        inquiry: ''
    });

    const [errors, setErrors] = useState<FormErrors>({});
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [submitSuccess, setSubmitSuccess] = useState(false);

    const validateForm = (): boolean => {
        const newErrors: FormErrors = {};

        if (!formData.name.trim()) {
            newErrors.name = 'Name is required';
        }

        if (!formData.email.trim()) {
            newErrors.email = 'Email is required';
        } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
            newErrors.email = 'Invalid email format';
        }

        if (!formData.inquiry.trim()) {
            newErrors.inquiry = 'Inquiry is required';
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        
        if (!validateForm()) {
            return;
        }

        setIsSubmitting(true);
        try {
            await submitContactForm({
                name: formData.name,
                email: formData.email,
                company: formData.company,
                inquiry: formData.inquiry
            });
            setSubmitSuccess(true);
            setFormData({ name: '', email: '', company: '', inquiry: '' });
        } catch (error) {
            console.error('Failed to submit contact form:', error);
            alert('Failed to submit form. Please try again later.');
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const fieldName = e.target.name as keyof FormData;
        const value = e.target.value;
        setFormData(prev => ({ ...prev, [fieldName]: value }));
        if (errors[fieldName as keyof FormErrors]) {
            setErrors(prev => ({ ...prev, [fieldName]: undefined }));
        }
    };

    return (
        <section id="contact" className="contact-section">
            <div className="container">
                <div className="section-title">
                    <h2 data-i18n="contact.title">{t('contact.title')}</h2>
                    <p data-i18n="contact.subtitle">{t('contact.subtitle')}</p>
                </div>

                <div className="contact-wrap">
                    <div className="contact-info">
                        <h3 data-i18n="contact.getInTouch">{t('contact.getInTouch')}</h3>
                        
                        <p>
                            <span className="contact-icon">📧</span>
                            <span id="contactEmail">
                                <strong data-i18n="contact.email">{t('contact.email')}</strong><br />
                                {siteSettings?.contactEmail || 'contact@medicaldevice.com'}
                            </span>
                        </p>
                        
                        <p>
                            <span className="contact-icon">💬</span>
                            <span id="contactWhatsApp">
                                <strong data-i18n="contact.whatsapp">{t('contact.whatsapp')}</strong><br />
                                {siteSettings?.contactWhatsapp || '+86 138 0000 0000'}
                            </span>
                        </p>
                        
                        <p>
                            <span className="contact-icon">📞</span>
                            <span id="contactPhone">
                                <strong data-i18n="contact.tel">{t('contact.tel')}</strong><br />
                                {siteSettings?.contactPhone || '+86 138 0000 0000'}
                            </span>
                        </p>
                        
                        <p>
                            <span className="contact-icon">📍</span>
                            <span id="contactAddress">
                                <strong data-i18n="contact.address">{t('contact.address')}</strong><br />
                                {siteSettings?.contactAddress || 'Industrial Park, Guangdong, China'}
                            </span>
                        </p>
                        
                        <p>
                            <span className="contact-icon">🕐</span>
                            <span>
                                <strong data-i18n="contact.hours">{t('contact.hours')}</strong><br />
                                <span data-i18n="contact.online">{t('contact.online')}</span>
                            </span>
                        </p>
                    </div>

                    <div className="contact-form">
                        {submitSuccess ? (
                            <div className="success-message">
                                <h4>Thank you!</h4>
                                <p>Your inquiry has been sent successfully. We will get back to you soon.</p>
                                <button 
                                    className="submit-btn" 
                                    onClick={() => setSubmitSuccess(false)}
                                >
                                    Send Another Inquiry
                                </button>
                            </div>
                        ) : (
                            <form onSubmit={handleSubmit}>
                                <input
                                    type="text"
                                    name="name"
                                    className="form-input"
                                    placeholder={t('contact.form.name')}
                                    value={formData.name}
                                    onChange={handleChange}
                                    required
                                />
                                {errors.name && <span className="error-message">{errors.name}</span>}

                                <input
                                    type="email"
                                    name="email"
                                    className="form-input"
                                    placeholder={t('contact.form.email')}
                                    value={formData.email}
                                    onChange={handleChange}
                                    required
                                />
                                {errors.email && <span className="error-message">{errors.email}</span>}

                                <input
                                    type="text"
                                    name="company"
                                    className="form-input"
                                    placeholder={t('contact.form.company')}
                                    value={formData.company}
                                    onChange={handleChange}
                                />

                                <textarea
                                    name="inquiry"
                                    className="form-input"
                                    placeholder={t('contact.form.inquiry')}
                                    value={formData.inquiry}
                                    onChange={handleChange}
                                    rows={5}
                                    required
                                />
                                {errors.inquiry && <span className="error-message">{errors.inquiry}</span>}

                                <button 
                                    type="submit" 
                                    className="submit-btn"
                                    disabled={isSubmitting}
                                >
                                    {isSubmitting ? 'Sending...' : t('contact.form.submit')}
                                </button>
                            </form>
                        )}
                    </div>
                </div>
            </div>
        </section>
    );
};

export default Contact;