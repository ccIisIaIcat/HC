import React from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import '../styles/Header.css';
import logo from '../assets/images/logo.png';

function Header() {
  const navigate = useNavigate();
  const location = useLocation();
  const isLoggedIn = localStorage.getItem('token');
  
  // 从本地存储中获取用户信息
  const getUserInfo = () => {
    const userStr = localStorage.getItem('user');
    if (!userStr) return null;
    
    try {
      return JSON.parse(userStr);
    } catch (error) {
      console.error('解析用户信息失败', error);
      return null;
    }
  };
  
  const user = getUserInfo();
  const isAdmin = user && user.Role === 'admin';
  
  // 检查是否在登录或注册页面
  const isAuthPage = location.pathname === '/login' || location.pathname === '/register' || location.pathname === '/verify-email';
  
  // 处理登出
  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    navigate('/login');
  };
  
  return (
    <header className={`app-header ${isAuthPage ? 'auth-header' : ''}`}>
      <div className="header-content">
        <div className="logo-container">
          <Link to={isLoggedIn ? (isAdmin ? '/admin/dashboard' : '/dashboard') : '/login'}>
            <div className="logo">
              <img src={logo} alt="健康检查系统" className="header-logo-image" />
              <span className="logo-text">健康检查系统</span>
            </div>
          </Link>
        </div>
        
        {isLoggedIn && !isAuthPage && (
          <nav className="header-nav">
            <ul>
              <li>
                <Link to="/dashboard" className={location.pathname === '/dashboard' ? 'active' : ''}>
                  仪表盘
                </Link>
              </li>
              <li>
                <Link to="/food-analysis" className={location.pathname === '/food-analysis' ? 'active' : ''}>
                  食物分析
                </Link>
              </li>
              <li>
                <Link to="/food-records" className={location.pathname === '/food-records' ? 'active' : ''}>
                  食物记录
                </Link>
              </li>
              <li>
                <Link to="/health-analysis" className={location.pathname === '/health-analysis' ? 'active' : ''}>
                  健康分析
                </Link>
              </li>
            </ul>
          </nav>
        )}
        
        {isLoggedIn && !isAuthPage && (
          <div className="user-actions">
            <span className="user-greeting">您好，{user?.Name || '用户'}</span>
            <button className="logout-button" onClick={handleLogout}>退出</button>
          </div>
        )}
      </div>
    </header>
  );
}

export default Header; 