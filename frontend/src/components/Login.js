import React, { useState, useEffect } from 'react';
import { useNavigate, Link, useLocation } from 'react-router-dom';
import '../styles/Auth.css';
import api from '../api';
import logo from '../assets/images/logo.png';

function Login() {
  const [formData, setFormData] = useState({
    email: '',
    password: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');
  const navigate = useNavigate();
  const location = useLocation();

  // 检查是否有来自注册页面的成功消息
  useEffect(() => {
    if (location.state?.message) {
      setMessage(location.state.message);
      // 清除location state，防止刷新页面后消息依然存在
      window.history.replaceState({}, document.title);
    }
  }, [location.state]);

  // 处理输入变化
  const handleChange = (e) => {
    const { id, value } = e.target;
    setFormData(prevData => ({
      ...prevData,
      [id]: value
    }));
  };

  // 处理表单提交的函数
  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    
    try {
      // 调用登录API
      const response = await api.auth.login(formData.email, formData.password);
      console.log('登录响应:', response);
      
      if (response && response.token) {
        // 登录成功，保存token和用户信息
        localStorage.setItem('token', response.token);
        localStorage.setItem('user', JSON.stringify(response.user));
        console.log('保存的用户信息:', response.user);
        console.log('用户角色:', response.user.Role);
        
        // 根据用户角色导航到不同页面
        if (response.user && response.user.Role === 'admin') {
          console.log('是管理员，准备导航到: /admin/dashboard');
          navigate('/admin/dashboard');
        } else {
          console.log('是普通用户，准备导航到: /dashboard');
          navigate('/dashboard');
        }
      } else {
        setError('登录失败，请检查邮箱和密码');
      }
    } catch (error) {
      console.error('登录出错:', error);
      
      // 处理特定错误情况
      if (error.response) {
        if (error.response.status === 401) {
          setError('邮箱或密码错误');
        } else if (error.response.status === 403 && error.response.data?.message?.includes('未验证')) {
          // 邮箱未验证，存储邮箱并跳转到验证页面
          localStorage.setItem('pendingUser', JSON.stringify({ email: formData.email }));
          navigate('/verify-email', { state: { email: formData.email } });
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
        <h2 className="auth-title">健康检查系统</h2>
        {message && (
          <div className="success-message">
            <span role="img" aria-label="成功" style={{ marginRight: '8px' }}>✅</span>
            {message}
          </div>
        )}
        {error && (
          <div className="error-message">
            <span role="img" aria-label="错误" style={{ marginRight: '8px' }}>❌</span>
            {error}
          </div>
        )}
        <form className="auth-form" onSubmit={handleSubmit}>
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
              placeholder="请输入密码"
              value={formData.password}
              onChange={handleChange}
              required
            />
          </div>
          <button 
            type="submit" 
            className="auth-button"
            disabled={loading}
          >
            {loading ? '登录中...' : '登录'}
          </button>
        </form>
        <div className="auth-footer">
          还没有账号？<Link to="/register">立即注册</Link>
        </div>
      </div>
    </div>
  );
}

export default Login; 