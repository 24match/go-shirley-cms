/**
 * 租户管理页面（超级管理员专用）
 */

import React, { useState, useEffect } from 'react';
import { api } from '../../services/api';
import { Tenant } from '../../types';
import TenantConfigPanel from './TenantConfigPanel';

const TenantList: React.FC = () => {
  const [tenants, setTenants] = useState<Tenant[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [showConfigModal, setShowConfigModal] = useState(false);
  const [configTenant, setConfigTenant] = useState<{ id: number; name: string } | null>(null);
  const [editingTenant, setEditingTenant] = useState<Tenant | null>(null);
  const [formData, setFormData] = useState({
    name: '',
    domain: '',
    email: '',
    password: '',
  });

  useEffect(() => {
    loadTenants();
  }, []);

  const loadTenants = async () => {
    setLoading(true);
    try {
      const data = await api.getTenants();
      setTenants(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error('Failed to load tenants:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setEditingTenant(null);
    setFormData({ name: '', domain: '', email: '', password: '' });
    setShowModal(true);
  };

  const handleConfig = (tenant: Tenant) => {
    setConfigTenant({ id: tenant.id, name: tenant.name });
    setShowConfigModal(true);
  };

  const handleEdit = (tenant: Tenant) => {
    setEditingTenant(tenant);
    setFormData({
      name: tenant.name,
      domain: tenant.domain,
      email: tenant.email,
      password: '',
    });
    setShowModal(true);
  };

  const handleActivate = async (id: number) => {
    try {
      await api.activateTenant(id);
      loadTenants();
    } catch (error) {
      console.error('Failed to activate tenant:', error);
    }
  };

  const handleDisable = async (id: number) => {
    try {
      await api.disableTenant(id);
      loadTenants();
    } catch (error) {
      console.error('Failed to disable tenant:', error);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('确定要删除此租户吗？此操作不可恢复。')) {
      return;
    }
    try {
      await api.deleteTenant(id);
      loadTenants();
    } catch (error) {
      console.error('Failed to delete tenant:', error);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editingTenant) {
        await api.updateTenant(editingTenant.id, {
          name: formData.name,
          domain: formData.domain,
          email: formData.email,
        });
      } else {
        await api.createTenant(formData);
      }
      setShowModal(false);
      loadTenants();
    } catch (error) {
      console.error('Failed to save tenant:', error);
    }
  };

  if (loading) {
    return <div className="loading">加载中...</div>;
  }

  return (
    <div className="tenant-list-page">
      <div className="page-header">
        <h1>租户管理</h1>
        <button className="btn btn-primary" onClick={handleCreate}>
          + 创建租户
        </button>
      </div>

      <table className="data-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>名称</th>
            <th>域名</th>
            <th>邮箱</th>
            <th>状态</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          {tenants.map((tenant) => (
            <tr key={tenant.id}>
              <td>{tenant.id}</td>
              <td>{tenant.name}</td>
              <td>{tenant.domain}</td>
              <td>{tenant.email}</td>
              <td>
                <span className={`status-badge ${tenant.status}`}>
                  {tenant.status === 'active' ? '活跃' : tenant.status === 'inactive' ? '未激活' : '已暂停'}
                </span>
              </td>
              <td>{new Date(tenant.createdAt).toLocaleDateString()}</td>
              <td className="actions">
                <button className="btn btn-sm btn-secondary" onClick={() => handleConfig(tenant)}>
                  配置
                </button>
                <button className="btn btn-sm btn-primary" onClick={() => handleEdit(tenant)}>
                  编辑
                </button>
                {tenant.status === 'active' ? (
                  <button className="btn btn-sm btn-warning" onClick={() => handleDisable(tenant.id)}>
                    禁用
                  </button>
                ) : (
                  <button className="btn btn-sm btn-success" onClick={() => handleActivate(tenant.id)}>
                    激活
                  </button>
                )}
                <button className="btn btn-sm btn-danger" onClick={() => handleDelete(tenant.id)}>
                  删除
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {showModal && (
        <div className="modal-overlay" onClick={() => setShowModal(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h2>{editingTenant ? '编辑租户' : '创建租户'}</h2>
              <button className="close-btn" onClick={() => setShowModal(false)}>×</button>
            </div>

            <form onSubmit={handleSubmit}>
              <div className="form-group">
                <label>租户名称</label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  required
                />
              </div>

              <div className="form-group">
                <label>域名</label>
                <input
                  type="text"
                  value={formData.domain}
                  onChange={(e) => setFormData({ ...formData, domain: e.target.value })}
                  required
                />
              </div>

              <div className="form-group">
                <label>邮箱</label>
                <input
                  type="email"
                  value={formData.email}
                  onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                  required
                />
              </div>

              {!editingTenant && (
                <div className="form-group">
                  <label>密码</label>
                  <input
                    type="password"
                    value={formData.password}
                    onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                    required
                    minLength={6}
                  />
                </div>
              )}

              <div className="modal-footer">
                <button type="button" className="btn btn-secondary" onClick={() => setShowModal(false)}>
                  取消
                </button>
                <button type="submit" className="btn btn-primary">
                  {editingTenant ? '更新' : '创建'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {showConfigModal && configTenant && (
        <TenantConfigPanel
          tenantId={configTenant.id}
          tenantName={configTenant.name}
          onClose={() => setShowConfigModal(false)}
          onSave={loadTenants}
        />
      )}
    </div>
  );
};

export default TenantList;