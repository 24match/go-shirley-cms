/**
 * API 服务层
 * 提供所有与后端 API 交互的方法
 */

import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';

const API_BASE_URL = '/api';

// 创建 axios 实例
const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器 - 添加 JWT token
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器 - 处理错误
apiClient.interceptors.response.use(
  (response) => {
    return response.data;
  },
  (error) => {
    if (error.response?.status === 401) {
      // Token 过期或无效，清除并跳转登录
      localStorage.removeItem('auth_token');
      localStorage.removeItem('user_info');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const api = {
  // 认证相关
  login: async (username: string, password: string) => {
    return apiClient.post('/login', { username, password });
  },

  // 模块管理
  getModules: async (moduleName?: string) => {
    return apiClient.get('/admin/modules', { params: { module: moduleName } });
  },
  getModule: async (name: string) => {
    return apiClient.get(`/admin/modules/${name}`);
  },
  saveModule: async (moduleName: string, data: FormData | unknown) => {
    if (data instanceof FormData) {
      return apiClient.post(`/admin/modules?moduleName=${moduleName}`, data, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
    }
    return apiClient.put(`/admin/modules/${moduleName}`, data);
  },
  deleteModule: async (name: string) => {
    return apiClient.delete(`/admin/modules/${name}`);
  },
  deleteModuleImage: async (name: string) => {
    return apiClient.delete(`/admin/modules/${name}/image`);
  },

  // 图片管理
  getImages: async () => {
    return apiClient.get('/admin/images');
  },
  uploadImage: async (file: File, description?: string, altText?: string) => {
    const formData = new FormData();
    formData.append('image', file);
    if (description) formData.append('description', description);
    if (altText) formData.append('altText', altText);
    return apiClient.post('/admin/images', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
  },
  uploadMultipleImages: async (files: File[]) => {
    const formData = new FormData();
    files.forEach((file) => formData.append('images', file));
    return apiClient.post('/admin/images/batch', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
  },
  updateImage: async (id: number, data: { description?: string; altText?: string }) => {
    return apiClient.put(`/admin/images/${id}`, data);
  },
  deleteImage: async (id: number) => {
    return apiClient.delete(`/admin/images/${id}`);
  },
  deleteImageByFilename: async (filename: string) => {
    return apiClient.delete(`/admin/images/by-filename/${filename}`);
  },

  // 内容管理
  getContent: async (module?: string) => {
    return apiClient.get('/admin/content', { params: { module } });
  },
  createContent: async (data: { module: string; title: string; content: string; imagePath?: string; sortOrder?: number; enabled?: boolean }) => {
    return apiClient.post('/admin/content', data);
  },
  updateContent: async (id: number, data: { title?: string; content?: string; imagePath?: string; sortOrder?: number; enabled?: boolean }) => {
    return apiClient.put(`/admin/content/${id}`, data);
  },
  deleteContent: async (id: number) => {
    return apiClient.delete(`/admin/content/${id}`);
  },

  // 语言管理
  getLanguages: async () => {
    return apiClient.get('/admin/lang');
  },
  createLanguage: async (data: { key: string; language: string; value: string }) => {
    return apiClient.post('/admin/lang', data);
  },
  updateLanguage: async (id: number, data: { value: string }) => {
    return apiClient.put(`/admin/lang/${id}`, data);
  },
  deleteLanguage: async (id: number) => {
    return apiClient.delete(`/admin/lang/${id}`);
  },

  // 站点设置
  getSiteSettings: async () => {
    return apiClient.get('/admin/site-settings');
  },
  saveSiteSettings: async (data: {
    siteName?: string;
    siteDescription?: string;
    logoPath?: string;
    faviconPath?: string;
    contactEmail?: string;
    contactPhone?: string;
    contactWhatsapp?: string;
    contactAddress?: string;
    socialLinks?: string;
    seoKeywords?: string;
    seoDescription?: string;
  }) => {
    return apiClient.post('/admin/site-settings', data);
  },

  // 联系表单提交管理
  getContactSubmissions: async () => {
    return apiClient.get('/admin/contact-submissions');
  },
  deleteContactSubmission: async (id: number) => {
    return apiClient.delete(`/admin/contact-submissions/${id}`);
  },

  // 超级管理员接口 - 租户管理
  getTenants: async () => {
    return apiClient.get('/superadmin/tenants');
  },
  getTenant: async (id: number) => {
    return apiClient.get(`/superadmin/tenants/${id}`);
  },
  createTenant: async (data: { name: string; domain: string; email: string; password: string }) => {
    return apiClient.post('/superadmin/tenants', data);
  },
  updateTenant: async (id: number, data: { name?: string; domain?: string; email?: string; status?: string }) => {
    return apiClient.put(`/superadmin/tenants/${id}`, data);
  },
  deleteTenant: async (id: number) => {
    return apiClient.delete(`/superadmin/tenants/${id}`);
  },
  activateTenant: async (id: number) => {
    return apiClient.post(`/superadmin/tenants/${id}/activate`);
  },
  disableTenant: async (id: number) => {
    return apiClient.post(`/superadmin/tenants/${id}/disable`);
  },

  // 超级管理员接口 - 系统统计
  getSystemStats: async () => {
    return apiClient.get('/superadmin/stats');
  },

  // 超级管理员接口 - 租户配置管理
  getTenantConfig: async (tenantId: number) => {
    return apiClient.get(`/superadmin/tenants/${tenantId}/config`);
  },
  updateTenantConfig: async (tenantId: number, data: { feature_flags?: Record<string, boolean>; resource_quota?: Record<string, number>; subscription_plan?: string; subscription_expires_at?: string }) => {
    return apiClient.put(`/superadmin/tenants/${tenantId}/config`, data);
  },
  resetQuota: async (tenantId: number, resourceType?: string) => {
    return apiClient.post(`/superadmin/tenants/${tenantId}/quota/reset`, null, { params: { resourceType } });
  },
  getQuotaUsage: async (tenantId: number) => {
    return apiClient.get(`/superadmin/tenants/${tenantId}/quota/usage`);
  },

  // 租户接口 - 配置管理
  getTenantFeatures: async () => {
    return apiClient.get('/tenant/features');
  },
  getTenantQuota: async () => {
    return apiClient.get('/tenant/quota');
  },
  getTenantConfigSelf: async () => {
    return apiClient.get('/tenant/config');
  },
};

export default api;