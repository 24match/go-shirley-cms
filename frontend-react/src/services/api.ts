/**
 * API 服务层
 * 负责与后端 CMS API 进行通信
 */

import type { 
    SiteSettings, 
    ModuleConfig, 
    ImageItem, 
    ContentItem, 
    PageConfig,
    Translations,
    CMSData 
} from '../types/cms';
import type { ApiResponse } from '../types/api';
import { API_ENDPOINTS } from '../types/api';

/**
 * 安全 fetch 包装器，处理错误并返回默认值
 */
export async function safeFetch<T>(
    url: string, 
    defaultValue: T,
    errorMessage?: string
): Promise<T> {
    try {
        const res = await fetch(url);
        if (!res.ok) {
            throw new Error(`HTTP ${res.status}: ${res.statusText}`);
        }
        const data: ApiResponse<T> = await res.json();
        return data.data ?? defaultValue;
    } catch (error) {
        console.warn(errorMessage ?? `API fetch failed: ${url}`, error);
        return defaultValue;
    }
}

/**
 * CMS API 服务
 */
export const cmsApi = {
    /**
     * 获取页面配置
     */
    async getConfig(): Promise<Record<string, any>> {
        const configArray: PageConfig[] = await safeFetch(
            API_ENDPOINTS.CONFIG, 
            [],
            'Failed to fetch config'
        );
        
        // 将 PageConfig 数组转换为以 pageName 为键的对象映射
        const config: Record<string, any> = {};
        if (Array.isArray(configArray)) {
            configArray.forEach(cfg => {
                if (cfg.pageName) {
                    try {
                        config[cfg.pageName] = JSON.parse(cfg.configData);
                    } catch (e) {
                        config[cfg.pageName] = {};
                    }
                }
            });
        }
        return config;
    },

    /**
     * 获取模块配置
     */
    async getModules(): Promise<Record<string, ModuleConfig>> {
        const modulesArray: ModuleConfig[] = await safeFetch(
            API_ENDPOINTS.MODULES, 
            [],
            'Failed to fetch modules'
        );
        
        const modules: Record<string, ModuleConfig> = {};
        if (Array.isArray(modulesArray)) {
            modulesArray.forEach(m => {
                modules[m.moduleName] = m;
            });
        }
        return modules;
    },

    /**
     * 获取图片列表
     */
    async getImages(): Promise<ImageItem[]> {
        return safeFetch(
            API_ENDPOINTS.IMAGES, 
            [],
            'Failed to fetch images'
        );
    },

    /**
     * 获取内容项
     */
    async getContent(): Promise<ContentItem[]> {
        return safeFetch(
            API_ENDPOINTS.CONTENT, 
            [],
            'Failed to fetch content'
        );
    },

    /**
     * 获取站点设置
     */
    async getSiteSettings(): Promise<SiteSettings> {
        return safeFetch(
            API_ENDPOINTS.SITE_SETTINGS, 
            {
                zhSiteLogo: '医疗',
                enSiteLogo: 'MEDICAL',
                siteLogoColor: '#06a499'
            },
            'Failed to fetch site settings'
        );
    },

    /**
     * 获取翻译数据
     */
    async getTranslations(): Promise<Translations> {
        return safeFetch(
            API_ENDPOINTS.LANG, 
            { en: {}, zh: {} },
            'Failed to fetch translations'
        );
    }
};

/**
 * 统一加载所有 CMS 数据
 */
export async function loadAllCMSData(): Promise<CMSData> {
    const [config, modules, images, content, siteSettings, translations] = await Promise.all([
        cmsApi.getConfig(),
        cmsApi.getModules(),
        cmsApi.getImages(),
        cmsApi.getContent(),
        cmsApi.getSiteSettings(),
        cmsApi.getTranslations()
    ]);

    return {
        config,
        modules,
        images,
        contentItems: content,
        siteSettings,
        translations
    };
}

/**
 * 提交联系表单
 */
export async function submitContactForm(data: {
    name: string;
    email: string;
    company?: string;
    inquiry: string;
}): Promise<void> {
    const res = await fetch(API_ENDPOINTS.CONTACT_SUBMIT, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    });
    
    if (!res.ok) {
        throw new Error(`Failed to submit contact form: ${res.status}`);
    }
}