/**
 * 租户配置管理组件（超级管理员专用）
 * 用于管理租户的功能模块开关、资源配额和订阅计划
 */

import React, { useState, useEffect } from 'react';
import { api } from '../../services/api';
import { TenantConfig, FeatureFlags, ResourceQuota } from '../../types';

interface TenantConfigPanelProps {
  tenantId: number;
  tenantName: string;
  onClose: () => void;
  onSave: () => void;
}

interface FormData {
  feature_flags: Record<string, boolean>;
  resource_quota: Record<string, number>;
  subscription_plan: string;
  subscription_expires_at: string;
}

const DEFAULT_FEATURE_FLAGS: FeatureFlags = {
  image_management: true,
  page_config: true,
  multi_language: true,
  contact_form: true,
  content_management: true,
};

const DEFAULT_RESOURCE_QUOTA: ResourceQuota = {
  max_images: 50,
  max_storage_mb: 512,
  max_content_items: 20,
  max_users: 3,
};

const SUBSCRIPTION_PLANS = [
  { value: 'free', label: '免费版', description: '基础功能，适合小型网站' },
  { value: 'pro', label: '专业版', description: '完整功能，适合中型企业' },
  { value: 'enterprise', label: '企业版', description: '无限制，适合大型企业' },
];

const TenantConfigPanel: React.FC<TenantConfigPanelProps> = ({
  tenantId,
  tenantName,
  onClose,
  onSave,
}) => {
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [formData, setFormData] = useState<FormData>({
    feature_flags: { ...DEFAULT_FEATURE_FLAGS },
    resource_quota: { ...DEFAULT_RESOURCE_QUOTA },
    subscription_plan: 'free',
    subscription_expires_at: '',
  });

  useEffect(() => {
    loadTenantConfig();
  }, [tenantId]);

  const loadTenantConfig = async () => {
    setLoading(true);
    try {
      const response = await api.getTenantConfig(tenantId);
      if (response.code === 'SUCCESS' && response.data) {
        const config: TenantConfig = response.data;
        setFormData({
          feature_flags: config.feature_flags || { ...DEFAULT_FEATURE_FLAGS },
          resource_quota: config.resource_quota || { ...DEFAULT_RESOURCE_QUOTA },
          subscription_plan: config.subscription_plan || 'free',
          subscription_expires_at: config.subscription_expires_at || '',
        });
      }
    } catch (error) {
      console.error('Failed to load tenant config:', error);
      alert('加载租户配置失败');
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      const response = await api.updateTenantConfig(tenantId, formData);
      if (response.code === 'SUCCESS') {
        alert('配置保存成功');
        onSave();
        onClose();
      } else {
        alert('配置保存失败：' + response.message);
      }
    } catch (error) {
      console.error('Failed to save tenant config:', error);
      alert('配置保存失败');
    } finally {
      setSaving(false);
    }
  };

  const handleFeatureToggle = (featureKey: string) => {
    setFormData((prev) => ({
      ...prev,
      feature_flags: {
        ...prev.feature_flags,
        [featureKey]: !prev.feature_flags[featureKey],
      },
    }));
  };

  const handleQuotaChange = (quotaKey: string, value: number) => {
    setFormData((prev) => ({
      ...prev,
      resource_quota: {
        ...prev.resource_quota,
        [quotaKey]: value,
      },
    }));
  };

  const handlePlanChange = (plan: string) => {
    // 根据订阅计划调整默认配额
    let newQuota: ResourceQuota;
    switch (plan) {
      case 'pro':
        newQuota = {
          max_images: 500,
          max_storage_mb: 5120,
          max_content_items: 100,
          max_users: 10,
        };
        break;
      case 'enterprise':
        newQuota = {
          max_images: -1,
          max_storage_mb: -1,
          max_content_items: -1,
          max_users: -1,
        };
        break;
      default:
        newQuota = { ...DEFAULT_RESOURCE_QUOTA };
    }
    setFormData((prev) => ({
      ...prev,
      subscription_plan: plan,
      resource_quota: newQuota,
    }));
  };

  if (loading) {
    return (
      <div className="modal-overlay" onClick={onClose}>
        <div className="modal">
          <div className="modal-body">
            <div className="loading">加载配置中...</div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal tenant-config-modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>租户配置 - {tenantName}</h2>
          <button className="close-btn" onClick={onClose}>×</button>
        </div>

        <div className="modal-body">
          {/* 订阅计划选择 */}
          <div className="form-section">
            <h3>订阅计划</h3>
            <div className="subscription-plans">
              {SUBSCRIPTION_PLANS.map((plan) => (
                <label
                  key={plan.value}
                  className={`plan-option ${formData.subscription_plan === plan.value ? 'selected' : ''}`}
                >
                  <input
                    type="radio"
                    name="subscription_plan"
                    value={plan.value}
                    checked={formData.subscription_plan === plan.value}
                    onChange={(e) => handlePlanChange(e.target.value)}
                  />
                  <div className="plan-info">
                    <span className="plan-name">{plan.label}</span>
                    <span className="plan-description">{plan.description}</span>
                  </div>
                </label>
              ))}
            </div>
          </div>

          {/* 订阅过期时间 */}
          <div className="form-section">
            <h3>订阅过期时间</h3>
            <div className="form-group">
              <input
                type="date"
                value={formData.subscription_expires_at ? formData.subscription_expires_at.split('T')[0] : ''}
                onChange={(e) =>
                  setFormData((prev) => ({
                    ...prev,
                    subscription_expires_at: e.target.value ? new Date(e.target.value).toISOString() : '',
                  }))
                }
                className="form-control"
              />
              <small className="form-hint">留空表示永不过期</small>
            </div>
          </div>

          {/* 功能模块开关 */}
          <div className="form-section">
            <h3>功能模块开关</h3>
            <div className="feature-toggles">
              <label className="toggle-item">
                <input
                  type="checkbox"
                  checked={formData.feature_flags.image_management ?? false}
                  onChange={() => handleFeatureToggle('image_management')}
                />
                <span>图片管理</span>
              </label>
              <label className="toggle-item">
                <input
                  type="checkbox"
                  checked={formData.feature_flags.page_config ?? false}
                  onChange={() => handleFeatureToggle('page_config')}
                />
                <span>页面配置</span>
              </label>
              <label className="toggle-item">
                <input
                  type="checkbox"
                  checked={formData.feature_flags.multi_language ?? false}
                  onChange={() => handleFeatureToggle('multi_language')}
                />
                <span>多语言支持</span>
              </label>
              <label className="toggle-item">
                <input
                  type="checkbox"
                  checked={formData.feature_flags.contact_form ?? false}
                  onChange={() => handleFeatureToggle('contact_form')}
                />
                <span>联系表单</span>
              </label>
              <label className="toggle-item">
                <input
                  type="checkbox"
                  checked={formData.feature_flags.content_management ?? false}
                  onChange={() => handleFeatureToggle('content_management')}
                />
                <span>内容管理</span>
              </label>
            </div>
          </div>

          {/* 资源配额配置 */}
          <div className="form-section">
            <h3>资源配额配置</h3>
            <div className="quota-inputs">
              <div className="form-group">
                <label>最大图片数量</label>
                <input
                  type="number"
                  value={formData.resource_quota.max_images ?? 50}
                  onChange={(e) => handleQuotaChange('max_images', parseInt(e.target.value) || 0)}
                  disabled={formData.subscription_plan === 'enterprise'}
                  className="form-control"
                />
                {formData.subscription_plan === 'enterprise' && (
                  <small className="form-hint">企业版无限制</small>
                )}
              </div>
              <div className="form-group">
                <label>最大存储空间 (MB)</label>
                <input
                  type="number"
                  value={formData.resource_quota.max_storage_mb ?? 512}
                  onChange={(e) => handleQuotaChange('max_storage_mb', parseInt(e.target.value) || 0)}
                  disabled={formData.subscription_plan === 'enterprise'}
                  className="form-control"
                />
              </div>
              <div className="form-group">
                <label>最大内容项数量</label>
                <input
                  type="number"
                  value={formData.resource_quota.max_content_items ?? 20}
                  onChange={(e) => handleQuotaChange('max_content_items', parseInt(e.target.value) || 0)}
                  disabled={formData.subscription_plan === 'enterprise'}
                  className="form-control"
                />
              </div>
              <div className="form-group">
                <label>最大用户数量</label>
                <input
                  type="number"
                  value={formData.resource_quota.max_users ?? 3}
                  onChange={(e) => handleQuotaChange('max_users', parseInt(e.target.value) || 0)}
                  disabled={formData.subscription_plan === 'enterprise'}
                  className="form-control"
                />
              </div>
            </div>
          </div>
        </div>

        <div className="modal-footer">
          <button type="button" className="btn btn-secondary" onClick={onClose} disabled={saving}>
            取消
          </button>
          <button type="button" className="btn btn-primary" onClick={handleSave} disabled={saving}>
            {saving ? '保存中...' : '保存配置'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default TenantConfigPanel;