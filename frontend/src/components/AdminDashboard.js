import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import '../styles/Dashboard.css';
import api from '../api';

function AdminDashboard() {
  const [user, setUser] = useState(null);
  const [stats, setStats] = useState(null);
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        // 从localStorage获取用户信息
        const storedUser = JSON.parse(localStorage.getItem('user'));
        console.log('当前用户信息:', storedUser);
        
        if (!storedUser || storedUser.Role !== 'admin') {
          throw new Error('无管理员权限');
        }
        setUser(storedUser);

        // 获取系统统计信息 - 添加调试日志
        console.log('调用API: api.admin.getStats()');
        const statsResponse = await api.admin.getStats();
        console.log('统计信息响应:', statsResponse);
        
        if (statsResponse) {
          // 处理不同的响应格式
          if (statsResponse.stats) {
            setStats(statsResponse.stats);
          } else {
            // 如果响应中没有stats字段，则将整个响应作为stats
            setStats(statsResponse);
          }
        }

        // 获取用户列表 - 添加调试日志
        console.log('调用API: api.admin.getUsers()');
        const usersResponse = await api.admin.getUsers();
        console.log('用户列表响应:', usersResponse);
        
        if (usersResponse) {
          // 处理不同的响应格式
          if (usersResponse.users) {
            setUsers(usersResponse.users);
          } else if (Array.isArray(usersResponse)) {
            // 如果响应是数组，则直接使用
            setUsers(usersResponse);
          }
        }
      } catch (err) {
        console.error('获取数据失败:', err);
        setError('获取数据失败，请重新登录或稍后再试');
        
        // 如果是权限问题，重定向到用户仪表盘
        if (err.message === '无管理员权限' || 
            (err.response && err.response.status === 403)) {
          navigate('/dashboard');
        }
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [navigate]);

  const handleLogout = () => {
    api.auth.logout();
    navigate('/login');
  };

  if (loading) {
    return <div className="dashboard-container loading">加载中...</div>;
  }

  if (error) {
    return (
      <div className="dashboard-container error">
        <div className="error-message">{error}</div>
        <button className="dashboard-button" onClick={() => navigate('/login')}>
          返回登录
        </button>
      </div>
    );
  }

  return (
    <div className="dashboard-container">
      <div className="dashboard-header">
        <h1>管理员控制台</h1>
        <div className="user-info">
          <span>欢迎，{user?.name || '管理员'}</span>
          <button className="logout-button" onClick={handleLogout}>
            退出登录
          </button>
        </div>
      </div>

      <div className="dashboard-content">
        <div className="stats-section">
          <h2>系统统计</h2>
          {stats ? (
            <div className="stats-grid">
              <div className="stat-card">
                <h3>总用户数</h3>
                <div className="stat-value">{stats.userCount || 0}</div>
              </div>
              <div className="stat-card">
                <h3>已验证用户</h3>
                <div className="stat-value">{stats.verifiedUserCount || 0}</div>
              </div>
              <div className="stat-card">
                <h3>管理员数量</h3>
                <div className="stat-value">{stats.adminCount || 0}</div>
              </div>
              <div className="stat-card">
                <h3>分析次数</h3>
                <div className="stat-value">{stats.analysisCount || 0}</div>
              </div>
              <div className="stat-card">
                <h3>7天活跃用户</h3>
                <div className="stat-value">{stats.activeUsersLast7Days || 0}</div>
              </div>
            </div>
          ) : (
            <p>无统计数据可用</p>
          )}
        </div>

        <div className="users-section">
          <h2>用户管理</h2>
          {users.length > 0 ? (
            <div className="users-table-container">
              <table className="users-table">
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>用户名</th>
                    <th>邮箱</th>
                    <th>角色</th>
                    <th>验证状态</th>
                    <th>注册时间</th>
                  </tr>
                </thead>
                <tbody>
                  {users.map(user => (
                    <tr key={user.ID}>
                      <td>{user.ID}</td>
                      <td>{user.Name}</td>
                      <td>{user.Email}</td>
                      <td>{user.Role === 'admin' ? '管理员' : '用户'}</td>
                      <td>{user.Verified ? '已验证' : '未验证'}</td>
                      <td>{new Date(user.CreatedAt).toLocaleString()}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <p>暂无用户数据</p>
          )}
        </div>
      </div>
    </div>
  );
}

export default AdminDashboard;