/**
 * 租户详情页面（超级管理员专用）
 * 展示租户基本信息、功能模块状态、资源配额使用情况和订阅计划信息
 */

import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { api } from '../../services/api';
import { Tenant, TenantConfig, QuotaInfo } from '../../types';
import QuotaProgress from '../../components/common/QuotaProgress';

const TenantDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  
  const [loading, setLoading] = useState(true);
  const [tenant, setTenant] = useState<Tenant | null>(null);
  const [config, setConfig] = useState<TenantConfig | null>(null);
  const [quotaInfo, setQuotaInfo] = useState<QuotaInfo[]>([]);

  useEffect(() => {
    if (id) {
      loadTenantDetail(id);
    }
  }, [id]);

  const loadTenantDetail = async (tenantId: string) => {
    setLoading(true);
    try {
      // 加载租户基本信息
      const tenantResponse = await api.getTenant(parseInt(tenantId));
      if (tenantResponse.code === 'SUCCESS') {
        setTenant(tenantResponse.data);
      }

      // 加载租户配置
      const configResponse = await api.getTenantConfig(parseInt(tenantId));
      if (configResponse.code === 'SUCCESS') {
        setConfig(configResponse.data);
      }

      // 加载配额使用情况
      const quotaResponse = await api.getQuotaUsage(parseInt(tenantId));
      if (quotaResponse.code === 'SUCCESS' && quotaResponse.data) {
        const { quota, usage } = quotaResponse.data;
        const quotaList: QuotaInfo[] = Object.entries(quota).map(([key, limit]) => {
          const resourceType = key.replace('max_', '');
          const usedKey = `used_${resourceType}`;
          const used = usage[usedKey] || 0;
          return {
            resource_type: resourceType,
            used,
            limit,
            available: limit === -1 ? -1 : limit - used,
            is_unlimited: limit === -1,
          };
        });
        setQuotaInfo(quotaList);
      }
    } catch (error) {
      console.error('加载租户详情失败:', error);
    } finally {
      setLoading(false);
    }
  };

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

  const getFeatureName = (featureKey: string): string => {
    const names: Record<string, string> = {
      image_management: '图片管理',
      page_config: '页面配置',
      multi_language: '多语言支持',
      contact_form: '联系表单',
      analytics: '数据分析',
      seo_tools: 'SEO 工具',
      custom_domain: '自定义域名',
      api_access: 'API 访问',
    };
    return names[featureKey] || featureKey;
  };

  if (loading) {
    return (
      <div className="tenant-detail loading">
        <div className="loading-spinner">加载中...</div>
      </div>
    );
  }

  if (!tenant) {
    return (
      <div className="tenant-detail error">
        <h2>租户不存在</h2>
        <button onClick={() => navigate('/admin/tenants')}>返回列表</button>
      </div>
    );
  }

  return (
    <div className="tenant-detail">
      <div className="tenant-detail-header">
        <button className="back-btn" onClick={() => navigate('/admin/tenants')}>
          ← 返回
        </button>
        <h1>租户详情：{tenant.name}</h1>
      </div>

      {/* 基本信息 */}
      <section className="tenant-detail-section">
        <h2>基本信息</h2>
        <div className="info-grid">
          <div className="info-item">
            <label>租户 ID</label>
            <span>{tenant.id}</span>
          </div>
          <div className="info-item">
            <label>租户名称</label>
            <span>{tenant.name}</span>
          </div>
          <div className="info-item">
            <label>域名</label>
            <span>{tenant.domain || '-'}</span>
          </div>
          <div className="info-item">
            <label>状态</label>
            <span className={`status-badge ${tenant.status}`}>
              {tenant.status === 'active' ? '正常' : tenant.status === 'suspended' ? '已暂停' : '未知'}
            </span>
          </div>
          <div className="info-item">
            <label>创建时间</label>
            <span>{tenant.created_at ? new Date(tenant.created_at).toLocaleString('zh-CN') : '-'}</span>
          </div>
          <div className="info-item">
            <label>更新时间</label>
            <span>{tenant.updated_at ? new Date(tenant.updated_at).toLocaleString('zh-CN') : '-'}</span>
          </div>
        </div>
      </section>

      {/* 订阅计划 */}
      {config && (
        <section className="tenant-detail-section">
          <h2>订阅计划</h2>
          <div className="info-grid">
            <div className="info-item">
              <label>计划类型</label>
              <span className="plan-badge">{config.subscription_plan || '基础版'}</span>
            </div>
            <div className="info-item">
              <label>到期时间</label>
              <span>
                {config.subscription_expires_at 
                  ? new Date(config.subscription_expires_at).toLocaleString('zh-CN')
                  : '永久有效'}
              </span>
            </div>
          </div>
        </section>
      )}

      {/* 功能模块状态 */}
      {config && config.feature_flags && (
        <section className="tenant-detail-section">
          <h2>功能模块状态</h2>
          <div className="feature-grid">
            {Object.entries(config.feature_flags).map(([key, enabled]) => (
              <div key={key} className={`feature-item ${enabled ? 'enabled' : 'disabled'}`}>
                <span className="feature-icon">{enabled ? '✓' : '✗'}</span>
                <span className="feature-name">{getFeatureName(key)}</span>
                <span className="feature-status">{enabled ? '已启用' : '已禁用'}</span>
              </div>
            ))}
          </div>
        </section>
      )}

      {/* 资源配额使用情况 */}
      {quotaInfo.length > 0 && (
        <section className="tenant-detail-section">
          <h2>资源配额使用情况</h2>
          <div className="quota-grid">
            {quotaInfo.map((info) => (
              <div key={info.resource_type} className="quota-item">
                <QuotaProgress
                  used={info.used}
                  limit={info.limit}
                  resourceName={getResourceName(info.resource_type)}
                />
              </div>
            ))}
          </div>
        </section>
      )}
    </div>
  );
};

export default TenantDetail;