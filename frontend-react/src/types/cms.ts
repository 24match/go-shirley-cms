/**
 * CMS 数据类型定义
 */

export interface SiteSettings {
    id?: number;
    zhSiteLogo: string;
    enSiteLogo: string;
    siteLogoColor: string;
    zhSiteTitle?: string;
    enSiteTitle?: string;
    contactEmail?: string;
    contactWhatsapp?: string;
    contactPhone?: string;
    contactAddress?: string;
    createdAt?: string;
    updatedAt?: string;
}

export interface ModuleConfig {
    id: number;
    moduleName: string;
    moduleType: string;
    enabled: boolean;
    zhTitle?: string;
    enTitle?: string;
    zhSubtitle?: string;
    enSubtitle?: string;
    zhContent?: string;
    enContent?: string;
    imagePath?: string;
    sortOrder: number;
    createdAt?: string;
    updatedAt?: string;
}

export interface ImageItem {
    id: number;
    filename: string;
    originalName: string;
    category: string;
    uploadedAt: string;
}

export interface ContentItem {
    id: number;
    pageName: string;
    contentType: string;
    zhContent: string;
    enContent: string;
    createdAt?: string;
    updatedAt?: string;
}

export interface PageConfig {
    id: number;
    pageName: string;
    configData: string;
    createdAt?: string;
    updatedAt?: string;
}

export interface Translations {
    en: Record<string, string>;
    zh: Record<string, string>;
}

export interface CMSData {
    config: Record<string, any>;
    modules: Record<string, ModuleConfig>;
    images: ImageItem[];
    contentItems: ContentItem[];
    siteSettings: SiteSettings;
    translations?: Translations;
}