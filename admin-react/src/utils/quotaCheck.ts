/**
 * 配额检查工具
 * 提供配额检查和提示功能
 */

import { api } from '../services/api';

export interface QuotaStatus {
  resourceType: string;
  used: number;
  limit: number;
  available: number;
  isUnlimited: boolean;
  percentage: number;
}

export interface QuotaCheckResult {
  allowed: boolean;
  reason?: string;
  status?: QuotaStatus;
}

/**
 * 检查指定资源的配额状态
 */
export const checkQuotaStatus = async (resourceType: string): Promise<QuotaCheckResult> => {
  try {
    const quotaData = await api.getTenantQuota();
    
    if (!quotaData || !quotaData.quota || !quotaData.usage) {
      return { allowed: true }; // 无配额限制
    }

    const { quota, usage } = quotaData;
    
    const quotaKey = `max_${resourceType}`;
    const usageKey = `used_${resourceType}`;
    
    const limit = quota[quotaKey];
    const used = usage[usageKey] || 0;
    
    // 无限制
    if (limit === -1 || limit === undefined) {
      return {
        allowed: true,
        status: {
          resourceType,
          used,
          limit: -1,
          available: -1,
          isUnlimited: true,
          percentage: 0,
        },
      };
    }
    
    const available = limit - used;
    const percentage = (used / limit) * 100;
    
    if (used >= limit) {
      return {
        allowed: false,
        reason: `配额已用完：${resourceType} 已使用 ${used}/${limit}`,
        status: {
          resourceType,
          used,
          limit,
          available: 0,
          isUnlimited: false,
          percentage,
        },
      };
    }
    
    return {
      allowed: true,
      status: {
        resourceType,
        used,
        limit,
        available,
        isUnlimited: false,
        percentage,
      },
    };
  } catch (error) {
    console.error('配额检查失败:', error);
    // 检查失败时默认允许，避免阻断用户操作
    return { allowed: true };
  }
};

/**
 * 检查并提示配额警告（当使用量超过 80% 时）
 */
export const checkQuotaWarning = async (resourceType: string): Promise<string | null> => {
  try {
    const result = await checkQuotaStatus(resourceType);
    
    if (result.status && !result.status.isUnlimited) {
      if (result.status.percentage >= 90) {
        return `⚠️ 警告：${getResourceName(resourceType)}已使用 ${Math.round(result.status.percentage)}%，接近上限！`;
      } else if (result.status.percentage >= 80) {
        return `⚠️ 注意：${getResourceName(resourceType)}已使用 ${Math.round(result.status.percentage)}%，请注意使用量。`;
      }
    }
    
    return null;
  } catch (error) {
    console.error('配额警告检查失败:', error);
    return null;
  }
};

/**
 * 获取资源的中文名称
 */
const getResourceName = (resourceType: string): string => {
  const names: Record<string, string> = {
    images: '图片数量',
    content_items: '内容项数量',
    users: '用户数量',
    storage_mb: '存储空间 (MB)',
    modules: '模块数量',
    languages: '语言数量',
  };
  return names[resourceType] || resourceType;
};

/**
 * 上传图片前的配额检查
 */
export const checkImageUploadQuota = async (): Promise<QuotaCheckResult> => {
  return checkQuotaStatus('images');
};

/**
 * 创建内容前的配额检查
 */
export const checkContentCreateQuota = async (): Promise<QuotaCheckResult> => {
  return checkQuotaStatus('content_items');
};

/**
 * 创建用户前的配额检查
 */
export const checkUserCreateQuota = async (): Promise<QuotaCheckResult> => {
  return checkQuotaStatus('users');
};

/**
 * 显示配额警告提示
 */
export const showQuotaWarning = async (): Promise<void> => {
  const resourceTypes = ['images', 'content_items', 'storage_mb', 'users'];
  
  for (const resourceType of resourceTypes) {
    const warning = await checkQuotaWarning(resourceType);
    if (warning) {
      // 使用 alert 或自定义提示组件
      console.warn(warning);
      // 在实际项目中，这里可以调用 UI 组件显示提示
      // 例如：toast.warning(warning);
    }
  }
};