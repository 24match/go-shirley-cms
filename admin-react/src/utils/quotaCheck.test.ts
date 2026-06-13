/**
 * 配额检查工具单元测试
 */

import {
  checkQuotaStatus,
  checkQuotaWarning,
  checkImageUploadQuota,
  checkContentCreateQuota,
  checkUserCreateQuota,
  getResourceName,
} from './quotaCheck';

// Mock API
const mockGetTenantQuota = jest.fn();
jest.mock('../services/api', () => ({
  api: {
    getTenantQuota: () => mockGetTenantQuota(),
  },
}));

describe('quotaCheck', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('checkQuotaStatus', () => {
    it('应该允许无配额限制的情况', async () => {
      mockGetTenantQuota.mockResolvedValue({
        quota: { max_images: -1 },
        usage: { used_images: 1000 },
      });

      const result = await checkQuotaStatus('images');

      expect(result.allowed).toBe(true);
      expect(result.status?.isUnlimited).toBe(true);
    });

    it('应该允许配额充足的情况', async () => {
      mockGetTenantQuota.mockResolvedValue({
        quota: { max_images: 100 },
        usage: { used_images: 50 },
      });

      const result = await checkQuotaStatus('images');

      expect(result.allowed).toBe(true);
      expect(result.status?.used).toBe(50);
      expect(result.status?.limit).toBe(100);
      expect(result.status?.available).toBe(50);
      expect(result.status?.percentage).toBe(50);
    });

    it('应该拒绝配额超限的情况', async () => {
      mockGetTenantQuota.mockResolvedValue({
        quota: { max_images: 100 },
        usage: { used_images: 100 },
      });

      const result = await checkQuotaStatus('images');

      expect(result.allowed).toBe(false);
      expect(result.reason).toContain('配额已用完');
    });

    it('应该处理 API 错误情况', async () => {
      mockGetTenantQuota.mockRejectedValue(new Error('API Error'));

      const result = await checkQuotaStatus('images');

      // 错误时默认允许，避免阻断用户操作
      expect(result.allowed).toBe(true);
    });

    it('应该处理空响应', async () => {
      mockGetTenantQuota.mockResolvedValue(null);

      const result = await checkQuotaStatus('images');

      expect(result.allowed).toBe(true);
    });
  });

  describe('checkQuotaWarning', () => {
    it('当使用量超过 90% 时返回警告', async () => {
      mockGetTenantQuota.mockResolvedValue({
        quota: { max_images: 100 },
        usage: { used_images: 95 },
      });

      const warning = await checkQuotaWarning('images');

      expect(warning).toContain('⚠️ 警告');
      expect(warning).toContain('95%');
    });

    it('当使用量超过 80% 但低于 90% 时返回注意', async () => {
      mockGetTenantQuota.mockResolvedValue({
        quota: { max_images: 100 },
        usage: { used_images: 85 },
      });

      const warning = await checkQuotaWarning('images');

      expect(warning).toContain('⚠️ 注意');
      expect(warning).toContain('85%');
    });

    it('当使用量低于 80% 时不返回警告', async () => {
      mockGetTenantQuota.mockResolvedValue({
        quota: { max_images: 100 },
        usage: { used_images: 50 },
      });

      const warning = await checkQuotaWarning('images');

      expect(warning).toBeNull();
    });

    it('无限制时不返回警告', async () => {
      mockGetTenantQuota.mockResolvedValue({
        quota: { max_images: -1 },
        usage: { used_images: 1000 },
      });

      const warning = await checkQuotaWarning('images');

      expect(warning).toBeNull();
    });
  });

  describe('便捷检查函数', () => {
    it('checkImageUploadQuota 应该检查图片配额', async () => {
      mockGetTenantQuota.mockResolvedValue({
        quota: { max_images: 100 },
        usage: { used_images: 50 },
      });

      const result = await checkImageUploadQuota();

      expect(result.status?.resourceType).toBe('images');
    });

    it('checkContentCreateQuota 应该检查内容项配额', async () => {
      mockGetTenantQuota.mockResolvedValue({
        quota: { max_content_items: 50 },
        usage: { used_content_items: 25 },
      });

      const result = await checkContentCreateQuota();

      expect(result.status?.resourceType).toBe('content_items');
    });

    it('checkUserCreateQuota 应该检查用户配额', async () => {
      mockGetTenantQuota.mockResolvedValue({
        quota: { max_users: 10 },
        usage: { used_users: 5 },
      });

      const result = await checkUserCreateQuota();

      expect(result.status?.resourceType).toBe('users');
    });
  });

  describe('getResourceName', () => {
    it('应该返回正确的资源名称', () => {
      expect(getResourceName('images')).toBe('图片数量');
      expect(getResourceName('content_items')).toBe('内容项数量');
      expect(getResourceName('users')).toBe('用户数量');
      expect(getResourceName('storage_mb')).toBe('存储空间 (MB)');
      expect(getResourceName('modules')).toBe('模块数量');
      expect(getResourceName('languages')).toBe('语言数量');
    });

    it('未知资源类型应该返回原始类型名', () => {
      expect(getResourceName('unknown')).toBe('unknown');
    });
  });
});