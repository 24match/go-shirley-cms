/**
 * 模块管理页面
 */

import React, { useState, useEffect } from 'react';
import { api } from '../../services/api';
import { ModuleConfig } from '../../types';

const ModuleList: React.FC = () => {
  const [modules, setModules] = useState<ModuleConfig[]>([]);
  const [loading, setLoading] = useState(true);
  const [editingModule, setEditingModule] = useState<ModuleConfig | null>(null);
  const [showModal, setShowModal] = useState(false);

  useEffect(() => {
    loadModules();
  }, []);

  const loadModules = async () => {
    setLoading(true);
    try {
      const data = await api.getModules();
      setModules(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error('Failed to load modules:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleEdit = (module: ModuleConfig) => {
    setEditingModule(module);
    setShowModal(true);
  };

  const handleToggleEnabled = async (moduleName: string, currentEnabled: boolean) => {
    try {
      await api.saveModule(moduleName, { enabled: !currentEnabled });
      loadModules();
    } catch (error) {
      console.error('Failed to update module:', error);
    }
  };

  const handleDelete = async (moduleName: string) => {
    if (!confirm(`确定要删除模块 "${moduleName}" 吗？`)) {
      return;
    }
    try {
      await api.deleteModule(moduleName);
      loadModules();
    } catch (error) {
      console.error('Failed to delete module:', error);
    }
  };

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingModule) return;

    try {
      await api.saveModule(editingModule.moduleName, editingModule);
      setShowModal(false);
      setEditingModule(null);
      loadModules();
    } catch (error) {
      console.error('Failed to save module:', error);
    }
  };

  if (loading) {
    return <div className="loading">加载中...</div>;
  }

  return (
    <div className="module-list-page">
      <div className="page-header">
        <h1>模块管理</h1>
      </div>

      <div className="module-grid">
        {modules.map((module) => (
          <div key={module.id} className="module-card">
            <div className="module-card-header">
              <h3>{module.moduleName}</h3>
              <span className={`status-badge ${module.enabled ? 'enabled' : 'disabled'}`}>
                {module.enabled ? '已启用' : '已禁用'}
              </span>
            </div>
            
            <div className="module-card-body">
              <p className="module-title">{module.title || module.zhTitle}</p>
              <p className="module-subtitle">{module.subtitle || module.zhSubtitle}</p>
            </div>

            <div className="module-card-footer">
              <button
                className="btn btn-secondary"
                onClick={() => handleToggleEnabled(module.moduleName, module.enabled)}
              >
                {module.enabled ? '禁用' : '启用'}
              </button>
              <button
                className="btn btn-primary"
                onClick={() => handleEdit(module)}
              >
                编辑
              </button>
              <button
                className="btn btn-danger"
                onClick={() => handleDelete(module.moduleName)}
              >
                删除
              </button>
            </div>
          </div>
        ))}
      </div>

      {showModal && editingModule && (
        <div className="modal-overlay" onClick={() => setShowModal(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h2>编辑模块：{editingModule.moduleName}</h2>
              <button className="close-btn" onClick={() => setShowModal(false)}>×</button>
            </div>
            
            <form onSubmit={handleSave}>
              <div className="form-group">
                <label>中文标题</label>
                <input
                  type="text"
                  value={editingModule.zhTitle || ''}
                  onChange={(e) => setEditingModule({ ...editingModule, zhTitle: e.target.value })}
                />
              </div>
              
              <div className="form-group">
                <label>英文标题</label>
                <input
                  type="text"
                  value={editingModule.enTitle || ''}
                  onChange={(e) => setEditingModule({ ...editingModule, enTitle: e.target.value })}
                />
              </div>
              
              <div className="form-group">
                <label>中文副标题</label>
                <input
                  type="text"
                  value={editingModule.zhSubtitle || ''}
                  onChange={(e) => setEditingModule({ ...editingModule, zhSubtitle: e.target.value })}
                />
              </div>
              
              <div className="form-group">
                <label>英文副标题</label>
                <input
                  type="text"
                  value={editingModule.enSubtitle || ''}
                  onChange={(e) => setEditingModule({ ...editingModule, enSubtitle: e.target.value })}
                />
              </div>
              
              <div className="form-group">
                <label>中文内容</label>
                <textarea
                  value={editingModule.zhContent || ''}
                  onChange={(e) => setEditingModule({ ...editingModule, zhContent: e.target.value })}
                  rows={4}
                />
              </div>
              
              <div className="form-group">
                <label>英文内容</label>
                <textarea
                  value={editingModule.enContent || ''}
                  onChange={(e) => setEditingModule({ ...editingModule, enContent: e.target.value })}
                  rows={4}
                />
              </div>

              <div className="modal-footer">
                <button type="button" className="btn btn-secondary" onClick={() => setShowModal(false)}>
                  取消
                </button>
                <button type="submit" className="btn btn-primary">
                  保存
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default ModuleList;