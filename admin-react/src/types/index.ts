/**
 * 用户类型定义
 */
export interface User {
  id: number;
  email: string;
  name: string;
  role: 'superadmin' | 'tenant' | 'user';
  tenantId?: number;
}

/**
 * 登录请求
 */
export interface LoginRequest {
  email: string;
  password: string;
}

/**
 * 登录响应
 */
export interface LoginResponse {
  token: string;
  user: User;
}

/**
 * 租户类型
 */
export interface Tenant {
  id: number;
  name: string;
  domain: string;
  email: string;
  status: 'active' | 'inactive' | 'suspended';
  createdAt: string;
  updatedAt: string;
}

/**
 * 模块配置类型
 */
export interface ModuleConfig {
  id: number;
  tenantId?: number;
  moduleName: string;
  enabled: boolean;
  title: string;
  subtitle: string;
  content: string;
  description: string;
  imagePath: string;
  sortOrder: number;
  extraData: string;
  zhTitle: string;
  enTitle: string;
  zhSubtitle: string;
  enSubtitle: string;
  zhContent: string;
  enContent: string;
  zhDescription: string;
  enDescription: string;
  createdAt: string;
  updatedAt: string;
}

/**
 * 图片类型
 */
export interface ImageItem {
  id: number;
  tenantId?: number;
  filename: string;
  originalName: string;
  description: string;
  altText: string;
  createdAt: string;
  updatedAt: string;
}

/**
 * 内容项类型
 */
export interface ContentItem {
  id: number;
  tenantId?: number;
  module: string;
  title: string;
  content: string;
  imagePath: string;
  sortOrder: number;
  enabled: boolean;
  createdAt: string;
  updatedAt: string;
}

/**
 * 语言文本类型
 */
export interface LanguageText {
  id: number;
  tenantId?: number;
  key: string;
  language: string;
  value: string;
  createdAt: string;
  updatedAt: string;
}

/**
 * 站点设置类型
 */
export interface SiteSettings {
  id: number;
  tenantId?: number;
  siteName: string;
  siteDescription: string;
  logoPath: string;
  faviconPath: string;
  contactEmail: string;
  contactPhone: string;
  contactWhatsapp: string;
  contactAddress: string;
  socialLinks: string;
  seoKeywords: string;
  seoDescription: string;
  createdAt: string;
  updatedAt: string;
}

/**
 * 联系表单提交类型
 */
export interface ContactSubmission {
  id: number;
  tenantId?: number;
  name: string;
  email: string;
  company: string;
  inquiry: string;
  status: 'new' | 'read' | 'replied' | 'archived';
  createdAt: string;
}

/**
 * API 响应类型
 */
export interface APIResponse<T = unknown> {
  code: string;
  message: string;
  data?: T;
}

/**
 * 系统统计类型
 */
export interface SystemStats {
  totalTenants: number;
  activeTenants: number;
  totalUsers: number;
  totalModules: number;
  totalImages: number;
  totalContentItems: number;
  totalContactSubmissions: number;
}

/**
 * 租户配置类型
 */
export interface TenantConfig {
  id: number;
  tenant_id: number;
  feature_flags: Record<string, boolean>;
  resource_quota: Record<string, number>;
  resource_usage: Record<string, number>;
  subscription_plan: string;
  subscription_expires_at?: string;
}

/**
 * 功能模块配置
 */
export interface FeatureFlags {
  image_management?: boolean;
  page_config?: boolean;
  multi_language?: boolean;
  contact_form?: boolean;
  content_management?: boolean;
  [key: string]: boolean | undefined;
}

/**
 * 资源配额配置
 */
export interface ResourceQuota {
  max_images?: number;
  max_storage_mb?: number;
  max_content_items?: number;
  max_users?: number;
  [key: string]: number | undefined;
}

/**
 * 资源使用情况
 */
export interface ResourceUsage {
  used_images?: number;
  used_storage_mb?: number;
  used_content_items?: number;
  used_users?: number;
  [key: string]: number | undefined;
}

/**
 * 更新租户配置请求
 */
export interface UpdateTenantConfigRequest {
  feature_flags?: Record<string, boolean>;
  resource_quota?: Record<string, number>;
  subscription_plan?: string;
  subscription_expires_at?: string;
}

/**
 * 配额信息
 */
export interface QuotaInfo {
  resource_type: string;
  used: number;
  limit: number;
  available: number;
  is_unlimited: boolean;
}