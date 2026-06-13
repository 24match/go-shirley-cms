/**
 * CMS Context
 * 管理全局 CMS 数据状态
 */

import React, { createContext, useContext, useState, useEffect, useCallback, ReactNode } from 'react';
import { loadAllCMSData, cmsApi } from '../services/api';
import { getLangField } from '../services/i18n';
import type { CMSData, SiteSettings, ModuleConfig, ImageItem, ContentItem } from '../types/cms';
import { getCurrentLang } from '../services/i18n';

export interface CMSContextType {
    data: CMSData | null;
    isLoading: boolean;
    error: Error | null;
    refreshData: () => Promise<void>;
    getLangField: <T extends Record<string, any>>(obj: T, baseField: string) => string;
    getSiteSettings: () => SiteSettings | null;
    getModule: (moduleName: string) => ModuleConfig | undefined;
    getImageByCategory: (category: string) => ImageItem | undefined;
}

const CMSContext = createContext<CMSContextType | undefined>(undefined);

// 默认 CMS 数据
const defaultCMSData: CMSData = {
    config: {},
    modules: {},
    images: [],
    contentItems: [],
    siteSettings: {
        zhSiteLogo: '医疗',
        enSiteLogo: 'MEDICAL',
        siteLogoColor: '#06a499'
    }
};

export function CMSProvider({ children }: { children: ReactNode }) {
    const [data, setData] = useState<CMSData | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<Error | null>(null);

    const loadData = useCallback(async () => {
        setIsLoading(true);
        setError(null);
        try {
            const cmsData = await loadAllCMSData();
            setData(cmsData);
        } catch (err) {
            console.error('Failed to load CMS data:', err);
            setError(err instanceof Error ? err : new Error('Unknown error'));
            // 使用默认数据降级
            setData(defaultCMSData);
        } finally {
            setIsLoading(false);
        }
    }, []);

    useEffect(() => {
        loadData();
    }, [loadData]);

    const refreshData = useCallback(async () => {
        await loadData();
    }, [loadData]);

    const getSiteSettings = useCallback((): SiteSettings | null => {
        return data?.siteSettings ?? null;
    }, [data]);

    const getModule = useCallback((moduleName: string): ModuleConfig | undefined => {
        return data?.modules[moduleName];
    }, [data]);

    const getImageByCategory = useCallback((category: string): ImageItem | undefined => {
        return data?.images.find(img => img.category === category);
    }, [data]);

    const value: CMSContextType = {
        data,
        isLoading,
        error,
        refreshData,
        getLangField: (obj, baseField) => getLangField(obj, baseField),
        getSiteSettings,
        getModule,
        getImageByCategory
    };

    return (
        <CMSContext.Provider value={value}>
            {children}
        </CMSContext.Provider>
    );
}

/**
 * 使用 CMS Context 的 Hook
 */
export function useCMSData(): CMSContextType {
    const context = useContext(CMSContext);
    if (context === undefined) {
        throw new Error('useCMSData must be used within a CMSProvider');
    }
    return context;
}

export default CMSContext;