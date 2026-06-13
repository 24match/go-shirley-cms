/**
 * 语言 Context
 * 管理全局语言状态和翻译功能
 */

import React, { createContext, useContext, useState, useEffect, useCallback, ReactNode } from 'react';
import { 
    initLanguage, 
    setLanguage as setLang,
    translate, 
    loadTranslationsFromAPI,
    updateSiteTitle,
    getCurrentLang
} from '../services/i18n';

export type Language = 'en' | 'zh';

export interface LanguageContextType {
    currentLang: Language;
    setLanguage: (lang: Language) => Promise<void>;
    t: (key: string) => string;
    isLoading: boolean;
}

const LanguageContext = createContext<LanguageContextType | undefined>(undefined);

export function LanguageProvider({ children }: { children: ReactNode }) {
    const [currentLang, setCurrentLang] = useState<Language>(() => initLanguage());
    const [isLoading, setIsLoading] = useState(false);

    // 初始化时加载翻译
    useEffect(() => {
        setIsLoading(true);
        loadTranslationsFromAPI()
            .finally(() => setIsLoading(false));
    }, []);

    const setLanguage = useCallback(async (lang: Language) => {
        setLang(lang);
        setCurrentLang(lang);
        await loadTranslationsFromAPI();
        await updateSiteTitle();
    }, []);

    const t = useCallback((key: string): string => {
        return translate(key);
    }, []);

    const value: LanguageContextType = {
        currentLang,
        setLanguage,
        t,
        isLoading
    };

    return (
        <LanguageContext.Provider value={value}>
            {children}
        </LanguageContext.Provider>
    );
}

/**
 * 使用语言 Context 的 Hook
 */
export function useLanguage(): LanguageContextType {
    const context = useContext(LanguageContext);
    if (context === undefined) {
        throw new Error('useLanguage must be used within a LanguageProvider');
    }
    return context;
}

export default LanguageContext;