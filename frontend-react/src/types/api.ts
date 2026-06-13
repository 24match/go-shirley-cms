/**
 * API 响应类型定义
 */

export interface ApiResponse<T> {
    code?: number;
    message?: string;
    data: T;
}

export interface ApiError {
    code: number;
    message: string;
    details?: string;
}

// API 端点路径常量
export const API_ENDPOINTS = {
    CONFIG: '/api/public/config',
    MODULES: '/api/public/modules',
    IMAGES: '/api/public/images',
    CONTENT: '/api/public/content',
    SITE_SETTINGS: '/api/public/site-settings',
    LANG: '/api/public/lang',
    CONTACT_SUBMIT: '/api/public/contact'
} as const;