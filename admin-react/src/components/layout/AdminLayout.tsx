/**
 * 管理后台布局组件
 * 包含侧边栏导航和顶部栏
 */

import React, { useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';

interface AdminLayoutProps {
  children: React.ReactNode;
}

const AdminLayout: React.FC<AdminLayoutProps> = ({ children }) => {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const { user, logout, isSuperAdmin } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  // 租户管理员菜单
  const tenantMenuItems = [
    { path: '/dashboard', label: '仪表盘', icon: '📊' },
    { path: '/modules', label: '模块管理', icon: '🔧' },
  ];

  // 超级管理员菜单
  const superAdminMenuItems = [
    { path: '/dashboard', label: '仪表盘', icon: '📊' },
    { path: '/tenants', label: '租户管理', icon: '🏢' },
  ];

  const menuItems = isSuperAdmin ? superAdminMenuItems : tenantMenuItems;

  return (
    <div className="admin-layout">
      {/* 侧边栏 */}
      <aside className={`sidebar ${sidebarOpen ? 'open' : 'closed'}`}>
        <div className="sidebar-header">
          <h2>
            <Link to="/dashboard">Medical CMS</Link>
          </h2>
          {isSuperAdmin && <span className="badge">超级管理员</span>}
        </div>
        
        <nav className="sidebar-nav">
          {menuItems.map((item) => (
            <Link
              key={item.path}
              to={item.path}
              className={`nav-item ${location.pathname === item.path ? 'active' : ''}`}
            >
              <span className="nav-icon">{item.icon}</span>
              <span className="nav-label">{item.label}</span>
            </Link>
          ))}
        </nav>
      </aside>

      {/* 主内容区 */}
      <div className="main-content">
        {/* 顶部栏 */}
        <header className="top-bar">
          <button
            className="sidebar-toggle"
            onClick={() => setSidebarOpen(!sidebarOpen)}
          >
            ☰
          </button>
          
          <div className="top-bar-right">
            <div className="user-info">
              <span className="user-name">{user?.name || user?.email}</span>
              <span className="user-role">{user?.role}</span>
            </div>
            <button className="logout-btn" onClick={handleLogout}>
              退出登录
            </button>
          </div>
        </header>

        {/* 页面内容 */}
        <main className="page-content">
          {children}
        </main>
      </div>
    </div>
  );
};

export default AdminLayout;