import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import '../styles/Auth.css';
import api from '../api';
import logo from '../assets/images/logo.png';

function Register() {
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    confirmPassword: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleChange = (e) => {
    const { id, value } = e.target;
    setFormData(prevData => ({
      ...prevData,
      [id]: value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    
    // 验证密码匹配
    if (formData.password !== formData.confirmPassword) {
      setError('两次输入的密码不一致');
      setLoading(false);
      return;
    }
    
    // 验证密码强度
    if (formData.password.length < 6) {
      setError('密码长度必须至少为6个字符');
      setLoading(false);
      return;
    }
    
    try {
      // 调用注册API
      const response = await api.auth.register(
        formData.username,
        formData.email,
        formData.password
      );
      
      // 注册成功，重定向到验证页面
      if (response && response.email) {
        // 存储注册邮箱到localStorage，以便验证页面使用
        localStorage.setItem('pendingUser', JSON.stringify({ email: response.email }));
        
        // 导航到验证页面，带上用户邮箱信息
        navigate('/verify-email', { state: { email: response.email } });
      } else {
        setError('注册失败，请稍后重试');
      }
    } catch (error) {
      console.error('注册出错:', error);
      
      // 处理特定错误情况
      if (error.response) {
        if (error.response.status === 409) {
          setError('该邮箱已被注册');
        } else if (error.response.status === 400) {
          setError(error.response.data.message || '请求数据格式错误');
        } else {
          setError(error.response.data.message || '服务器错误，请稍后再试');
        }
      } else {
        setError('网络错误，请检查网络连接');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-container auth-page">
      <div className="auth-box">
        <div className="auth-logo">
          <img src={logo} alt="HC" className="auth-logo-image" />
        </div>
        <h2 className="auth-title">注册健康检查系统</h2>
        {error && <div className="error-message">{error}</div>}
        <form className="auth-form" onSubmit={handleSubmit}>
          <div className="form-item">
            <label className="form-label" htmlFor="username">用户名</label>
            <input 
              type="text" 
              id="username" 
              className="form-input" 
              placeholder="请输入用户名"
              value={formData.username}
              onChange={handleChange}
              required
            />
          </div>
          <div className="form-item">
            <label className="form-label" htmlFor="email">邮箱</label>
            <input 
              type="email" 
              id="email" 
              className="form-input" 
              placeholder="请输入邮箱"
              value={formData.email}
              onChange={handleChange}
              required
            />
          </div>
          <div className="form-item">
            <label className="form-label" htmlFor="password">密码</label>
            <input 
              type="password" 
              id="password" 
              className="form-input" 
              placeholder="请输入密码（至少6个字符）"
              value={formData.password}
              onChange={handleChange}
              required
              minLength="6"
            />
          </div>
          <div className="form-item">
            <label className="form-label" htmlFor="confirmPassword">确认密码</label>
            <input 
              type="password" 
              id="confirmPassword" 
              className="form-input" 
              placeholder="请再次输入密码"
              value={formData.confirmPassword}
              onChange={handleChange}
              required
            />
          </div>
          <button 
            type="submit" 
            className="auth-button"
            disabled={loading}
          >
            {loading ? '注册中...' : '注册'}
          </button>
        </form>
        <div className="auth-footer">
          已有账号？<Link to="/login">立即登录</Link>
        </div>
      </div>
    </div>
  );
}

export default Register; 