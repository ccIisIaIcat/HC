import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation, Link } from 'react-router-dom';
import '../styles/Auth.css';
import api from '../api';
import logo from '../assets/images/logo.png';

function VerifyEmail() {
  const [verificationCode, setVerificationCode] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [countdown, setCountdown] = useState(60);
  const [canResend, setCanResend] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  
  // 从location获取邮箱，如果没有则从localStorage获取
  const email = location.state?.email || JSON.parse(localStorage.getItem('pendingUser') || '{}').email;
  
  // 如果没有邮箱信息，重定向到注册页面
  useEffect(() => {
    if (!email) {
      navigate('/register');
    }
  }, [email, navigate]);
  
  // 倒计时逻辑
  useEffect(() => {
    let timer;
    if (countdown > 0 && !canResend) {
      timer = setTimeout(() => setCountdown(countdown - 1), 1000);
    } else {
      setCanResend(true);
    }
    
    return () => clearTimeout(timer);
  }, [countdown, canResend]);
  
  const handleChange = (e) => {
    // 只接受数字输入
    const input = e.target.value.replace(/\D/g, '');
    
    // 限制验证码长度为6位
    if (input.length <= 6) {
      setVerificationCode(input);
    }
  };
  
  const handleResend = async () => {
    if (!canResend) return;
    
    setLoading(true);
    setError('');
    
    try {
      // 调用重新发送验证码API
      const response = await api.auth.resendVerificationCode(email);
      
      if (response && response.message) {
        // 重置倒计时
        setCountdown(60);
        setCanResend(false);
        setError(''); // 清除之前的错误
      } else {
        setError('重新发送验证码失败，请稍后重试');
      }
    } catch (error) {
      console.error('重新发送验证码出错:', error);
      if (error.response && error.response.data.message) {
        setError(error.response.data.message);
      } else {
        setError('网络错误，请检查网络连接');
      }
    } finally {
      setLoading(false);
    }
  };
  
  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    
    // 验证码长度验证
    if (verificationCode.length !== 6) {
      setError('请输入6位验证码');
      setLoading(false);
      return;
    }
    
    try {
      // 调用验证邮箱API
      const response = await api.auth.verifyEmail(email, verificationCode);
      
      if (response && response.message) {
        // 验证成功，清除本地存储的待验证用户
        localStorage.removeItem('pendingUser');
        
        // 导航到登录页面，带上验证成功的消息
        navigate('/login', { 
          state: { message: '邮箱验证成功，请登录您的账号' }
        });
      } else {
        setError('验证失败，请检查验证码是否正确');
      }
    } catch (error) {
      console.error('验证邮箱出错:', error);
      
      // 处理特定错误情况
      if (error.response) {
        if (error.response.status === 400) {
          setError('验证码无效或已过期');
        } else if (error.response.status === 404) {
          setError('用户不存在');
        } else if (error.response.status === 410) {
          setError('验证码已过期，请重新发送');
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
  
  // 如果没有邮箱信息，显示加载中
  if (!email) {
    return <div className="auth-container">正在加载...</div>;
  }
  
  return (
    <div className="auth-container auth-page">
      <div className="auth-box">
        <Link to="/register" className="back-link">← 返回注册</Link>
        <div className="auth-logo">
          <img src={logo} alt="HC" className="auth-logo-image" />
        </div>
        <h2 className="auth-title">验证您的邮箱</h2>
        <div className="info-message">
          验证码已发送至邮箱<br />
          <strong>{email}</strong>
        </div>
        {error && <div className="error-message">{error}</div>}
        <form className="auth-form" onSubmit={handleSubmit}>
          <div className="form-item">
            <label className="form-label" htmlFor="verificationCode">验证码</label>
            <input 
              type="text" 
              id="verificationCode" 
              className="form-input verification-code-input" 
              placeholder="请输入6位验证码"
              value={verificationCode}
              onChange={handleChange}
              required
              maxLength="6"
            />
          </div>
          <button 
            type="submit" 
            className="auth-button"
            disabled={loading || verificationCode.length !== 6}
          >
            {loading ? '验证中...' : '验证邮箱'}
          </button>
        </form>
        <div className="auth-footer">
          没有收到验证码？
          <button 
            className="resend-button"
            onClick={handleResend}
            disabled={!canResend || loading}
          >
            {canResend ? '重新发送' : `${countdown}秒后可重新发送`}
          </button>
        </div>
      </div>
    </div>
  );
}

export default VerifyEmail; 