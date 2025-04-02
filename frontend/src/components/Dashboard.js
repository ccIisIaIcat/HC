import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import '../styles/Dashboard.css';
import api from '../api';

function Dashboard() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    // 通过本地存储获取用户数据
    const userStr = localStorage.getItem('user');
    if (userStr) {
      try {
        const userData = JSON.parse(userStr);
        setUser(userData);
      } catch (error) {
        console.error('解析用户数据失败:', error);
      }
    }
    
    // 尝试刷新用户信息
    const refreshUserData = async () => {
      try {
        const data = await api.auth.getCurrentUser();
        if (data && data.user) {
          // 更新本地存储的用户信息
          localStorage.setItem('user', JSON.stringify(data.user));
          setUser(data.user);
        }
      } catch (error) {
        console.error('刷新用户信息失败:', error);
        // 发生错误但不重定向，因为重定向在ProtectedRoute中处理
      } finally {
        setLoading(false);
      }
    };

    refreshUserData();
  }, []);

  // 格式化日期
  const formatDate = (dateString) => {
    if (!dateString) return '未知';
    const date = new Date(dateString);
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  if (loading && !user) {
    return <div className="dashboard-container loading">加载中...</div>;
  }

  return (
    <div className="dashboard-container">
      <main className="dashboard-content">
        <h2>个人中心</h2>
        <div className="user-profile-card">
          <div className="user-profile-header">
            <div className="user-avatar">
              <span>{user?.Name?.charAt(0) || '用'}</span>
            </div>
            <div className="user-info">
              <h3>{user?.Name || '用户'}</h3>
              <p className="user-email">{user?.Email || '未知邮箱'}</p>
            </div>
          </div>
          <div className="user-profile-details">
            <div className="profile-item">
              <span className="profile-label">账户类型</span>
              <span className="profile-value">{user?.Role === 'admin' ? '管理员' : '普通用户'}</span>
            </div>
            <div className="profile-item">
              <span className="profile-label">账户状态</span>
              <span className="profile-value">{user?.Verified ? '已验证' : '未验证'}</span>
            </div>
            <div className="profile-item">
              <span className="profile-label">注册时间</span>
              <span className="profile-value">{formatDate(user?.CreatedAt)}</span>
            </div>
          </div>
        </div>

        <div className="dashboard-welcome">
          <h3>欢迎使用健康检查系统</h3>
          <p>您可以通过上方导航栏访问系统的各项功能：</p>
          <ul className="feature-list">
            <li><span className="feature-icon">🍎</span> <strong>食物分析</strong> - 上传食物图片进行营养成分分析</li>
            <li><span className="feature-icon">📋</span> <strong>食物记录</strong> - 查看您的饮食历史记录</li>
            <li><span className="feature-icon">📊</span> <strong>健康分析</strong> - 获取个性化的健康分析报告</li>
          </ul>
        </div>
      </main>
    </div>
  );
}

export default Dashboard; 