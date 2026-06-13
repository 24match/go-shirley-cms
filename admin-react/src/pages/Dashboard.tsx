/**
 * 仪表盘页面
 */

import React, { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { api } from '../services/api';
import { SystemStats, ModuleConfig, ContactSubmission, TenantConfig, QuotaInfo } from '../types';
import QuotaProgress from '../components/common/QuotaProgress';

const Dashboard: React.FC = () => {
  const { isSuperAdmin, user } = useAuth();
  const [stats, setStats] = useState<SystemStats | null>(null);
  const [modules, setModules] = useState<ModuleConfig[]>([]);
  const [recentSubmissions, setRecentSubmissions] = useState<ContactSubmission[]>([]);
  const [loading, setLoading] = useState(true);
  
  // 租户管理员专用状态
  const [tenantConfig, setTenantConfig] = useState<TenantConfig | null>(null);
  const [quotaInfo, setQuotaInfo] = useState<QuotaInfo[]>([]);
  const [features, setFeatures] = useState<Record<string, boolean>>({});

  useEffect(() => {
    loadDashboardData();
  }, [isSuperAdmin]);

  const loadDashboardData = async () => {
    setLoading(true);
    try {
      if (isSuperAdmin) {
        // 超级管理员：加载系统统计
        const statsData = await api.getSystemStats();
        setStats(statsData);
      } else {
        // 租户管理员：加载模块和联系表单
        const modulesData = await api.getModules();
        setModules(Array.isArray(modulesData) ? modulesData : []);
        
        const submissionsData = await api.getContactSubmissions();
        setRecentSubmissions(Array.isArray(submissionsData) ? submissionsData.slice(0, 5) : []);
        
        // 加载功能模块状态
        try {
          const featuresData = await api.getTenantFeatures();
          if (featuresData && featuresData.features) {
            setFeatures(featuresData.features);
          }
        } catch (e) {
          console.warn('加载功能状态失败:', e);
        }
        
        // 加载配额使用情况
        try {
          const quotaData = await api.getTenantQuota();
          if (quotaData && quotaData.quota && quotaData.usage) {
            const { quota, usage } = quotaData;
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
        } catch (e) {
          console.warn('加载配额信息失败:', e);
        }
        
        // 加载租户配置
        try {
          const configData = await api.getTenantConfigSelf();
          if (configData) {
            setTenantConfig(configData);
          }
        } catch (e) {
          console.warn('加载租户配置失败:', e);
        }
      }
    } catch (error) {
      console.error('Failed to load dashboard data:', error);
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
      content_management: '内容管理',
      analytics: '数据分析',
      seo_tools: 'SEO 工具',
    };
    return names[featureKey] || featureKey;
  };

  // 检查是否有配额接近上限
  const getQuotaWarning = (): string | null => {
    for (const info of quotaInfo) {
      if (info.is_unlimited) continue;
      const percentage = (info.used / info.limit) * 100;
      if (percentage >= 90) {
        return `${getResourceName(info.resource_type)}已使用 ${Math.round(percentage)}%，接近上限！`;
      }
    }
    return null;
  };

  const quotaWarning = getQuotaWarning();

  if (loading) {
    return <div className="dashboard-loading">加载中...</div>;
  }

  return (
    <div className="dashboard">
      <h1>仪表盘</h1>
      
      {/* 配额警告提示 */}
      {quotaWarning && (
        <div className="quota-warning">
          ⚠️ {quotaWarning}
        </div>
      )}

      {isSuperAdmin ? (
        /* 超级管理员视图 */
        <div className="stats-grid">
          <div className="stat-card">
            <h3>总租户数</h3>
            <p className="stat-value">{stats?.totalTenants || 0}</p>
          </div>
          <div className="stat-card">
            <h3>活跃租户</h3>
            <p className="stat-value">{stats?.activeTenants || 0}</p>
          </div>
          <div className="stat-card">
            <h3>总用户数</h3>
            <p className="stat-value">{stats?.totalUsers || 0}</p>
          </div>
          <div className="stat-card">
            <h3>总模块数</h3>
            <p className="stat-value">{stats?.totalModules || 0}</p>
          </div>
          <div className="stat-card">
            <h3>总图片数</h3>
            <p className="stat-value">{stats?.totalImages || 0}</p>
          </div>
          <div className="stat-card">
            <h3>联系表单</h3>
            <p className="stat-value">{stats?.totalContactSubmissions || 0}</p>
          </div>
        </div>
      ) : (
        /* 租户管理员视图 */
        <div className="tenant-dashboard">
          {/* 功能模块状态 */}
          {Object.keys(features).length > 0 && (
            <section className="dashboard-section">
              <h2>功能模块状态</h2>
              <div className="feature-grid">
                {Object.entries(features).map(([key, enabled]) => (
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
            <section className="dashboard-section">
              <h2>资源配额使用情况</h2>
              <div className="quota-overview-grid">
                {quotaInfo.map((info) => (
                  <div key={info.resource_type} className="quota-item">
                    <QuotaProgress
                      used={info.used}
                      limit={info.limit}
                      resourceName={getResourceName(info.resource_type)}
                      showLabels={true}
                    />
                  </div>
                ))}
              </div>
            </section>
          )}

          {/* 订阅计划信息 */}
          {tenantConfig && (
            <section className="dashboard-section">
              <h2>订阅计划</h2>
              <div className="subscription-info">
                <span className="plan-name">{tenantConfig.subscription_plan || '基础版'}</span>
                {tenantConfig.subscription_expires_at && (
                  <span className="plan-expires">
                    到期：{new Date(tenantConfig.subscription_expires_at).toLocaleDateString('zh-CN')}
                  </span>
                )}
              </div>
            </section>
          )}

          {/* 模块概览 */}
          <section className="dashboard-section">
            <h2>模块概览</h2>
            <div className="modules-overview">
              {modules.map((module) => (
                <div key={module.id} className="module-item">
                  <span className="module-name">{module.moduleName}</span>
                  <span className={`module-status ${module.enabled ? 'enabled' : 'disabled'}`}>
                    {module.enabled ? '已启用' : '已禁用'}
                  </span>
                </div>
              ))}
            </div>
          </section>

          {/* 最近联系表单 */}
          <section className="dashboard-section">
            <h2>最近联系表单</h2>
            {recentSubmissions.length > 0 ? (
              <table className="submissions-table">
                <thead>
                  <tr>
                    <th>姓名</th>
                    <th>邮箱</th>
                    <th>公司</th>
                    <th>日期</th>
                  </tr>
                </thead>
                <tbody>
                  {recentSubmissions.map((submission) => (
                    <tr key={submission.id}>
                      <td>{submission.name}</td>
                      <td>{submission.email}</td>
                      <td>{submission.company || '-'}</td>
                      <td>{new Date(submission.createdAt).toLocaleDateString()}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            ) : (
              <p className="empty-message">暂无联系表单</p>
            )}
          </section>
        </div>
      )}
    </div>
  );
};

export default Dashboard;