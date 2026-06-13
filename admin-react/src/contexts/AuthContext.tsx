/**
 * 认证上下文
 * 管理用户登录状态和权限
 */

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { User, LoginResponse } from '../types';
import { api } from '../services/api';

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
  isSuperAdmin: boolean;
  features: Record<string, boolean>;
  isFeatureEnabled: (featureName: string) => boolean;
  loadFeatures: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [features, setFeatures] = useState<Record<string, boolean>>({});

  useEffect(() => {
    // 从 localStorage 恢复用户信息
    const token = localStorage.getItem('auth_token');
    const userInfo = localStorage.getItem('user_info');
    
    if (token && userInfo) {
      try {
        setUser(JSON.parse(userInfo));
      } catch (e) {
        localStorage.removeItem('auth_token');
        localStorage.removeItem('user_info');
      }
    }
    setIsLoading(false);
  }, []);

  // 加载功能模块状态（仅租户管理员）
  useEffect(() => {
    if (user && user.role !== 'superadmin') {
      loadFeatures();
    }
  }, [user]);

  const loadFeatures = async () => {
    try {
      const featuresData = await api.getTenantFeatures();
      if (featuresData && featuresData.features) {
        setFeatures(featuresData.features);
        localStorage.setItem('tenant_features', JSON.stringify(featuresData.features));
      }
    } catch (e) {
      console.warn('加载功能状态失败:', e);
      // 从 localStorage 恢复
      const cached = localStorage.getItem('tenant_features');
      if (cached) {
        try {
          setFeatures(JSON.parse(cached));
        } catch (e) {
          // ignore
        }
      }
    }
  };

  const isFeatureEnabled = (featureName: string): boolean => {
    // 超级管理员默认所有功能都启用
    if (user?.role === 'superadmin') {
      return true;
    }
    return features[featureName] !== false;
  };

  const login = async (username: string, password: string) => {
    try {
      const response: LoginResponse = await api.login(username, password);
      
      localStorage.setItem('auth_token', response.token);
      localStorage.setItem('user_info', JSON.stringify(response.user));
      
      setUser(response.user);
    } catch (error) {
      console.error('Login failed:', error);
      throw error;
    }
  };

  const logout = () => {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('user_info');
    setUser(null);
  };

  const isSuperAdmin = user?.role === 'superadmin';

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated: !!user,
        isLoading,
        login,
        logout,
        isSuperAdmin,
        features,
        isFeatureEnabled,
        loadFeatures,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}

export default AuthContext;