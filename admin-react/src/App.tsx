/**
 * 管理后台 App 根组件
 */

import React from 'react';
import { BrowserRouter, Routes, Route, Navigate, useLocation } from 'react-router-dom';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import AdminLayout from './components/layout/AdminLayout';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import ModuleList from './pages/modules/ModuleList';
import TenantList from './pages/tenants/TenantList';
import './styles/admin.css';

// 需要认证的路由包装器
function ProtectedRoute({ children, superAdminOnly = false }: { children: React.ReactNode; superAdminOnly?: boolean }) {
  const { isAuthenticated, isLoading, isSuperAdmin } = useAuth();
  const location = useLocation();

  if (isLoading) {
    return <div className="loading-fullscreen">加载中...</div>;
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  if (superAdminOnly && !isSuperAdmin) {
    return <Navigate to="/dashboard" replace />;
  }

  return <>{children}</>;
}

// 带布局的路由包装器
function LayoutRoute({ children }: { children: React.ReactNode }) {
  return (
    <AdminLayout>
      {children}
    </AdminLayout>
  );
}

function AppRoutes() {
  return (
    <Routes>
      {/* 公开路由 */}
      <Route path="/login" element={<Login />} />

      {/* 需要认证的路由 */}
      <Route
        path="/dashboard"
        element={
          <ProtectedRoute>
            <LayoutRoute>
              <Dashboard />
            </LayoutRoute>
          </ProtectedRoute>
        }
      />

      <Route
        path="/modules"
        element={
          <ProtectedRoute>
            <LayoutRoute>
              <ModuleList />
            </LayoutRoute>
          </ProtectedRoute>
        }
      />

      <Route
        path="/tenants"
        element={
          <ProtectedRoute superAdminOnly>
            <LayoutRoute>
              <TenantList />
            </LayoutRoute>
          </ProtectedRoute>
        }
      />

      {/* 默认重定向 */}
      <Route path="/" element={<Navigate to="/dashboard" replace />} />
      <Route path="*" element={<Navigate to="/dashboard" replace />} />
    </Routes>
  );
}

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <AppRoutes />
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;