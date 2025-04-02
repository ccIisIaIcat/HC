import React, { useEffect, useState } from 'react';
import { 
  BrowserRouter as Router, 
  Routes, 
  Route, 
  Navigate, 
  useNavigate, 
  useLocation,
  createRoutesFromElements,
  createBrowserRouter
} from 'react-router-dom';
import './App.css';
import Login from './components/Login';
import Register from './components/Register';
import VerifyEmail from './components/VerifyEmail';
import Dashboard from './components/Dashboard';
import AdminDashboard from './components/AdminDashboard';
import FoodAnalysis from './components/FoodAnalysis';
import FoodRecords from './components/FoodRecords';
import FoodRecordEdit from './components/FoodRecordEdit';
import HealthAnalysis from './components/HealthAnalysis';
import Header from './components/Header';

// 受保护的路由组件，用于检查认证状态
function ProtectedRoute({ children, adminOnly = false }) {
  const navigate = useNavigate();
  const location = useLocation();
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isAdmin, setIsAdmin] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // 检查用户认证状态
    const token = localStorage.getItem('token');
    const userStr = localStorage.getItem('user');
    
    if (!token || !userStr) {
      // 未登录，重定向到登录页面
      navigate('/login', { state: { from: location.pathname } });
    } else {
      try {
        // 解析用户信息
        const user = JSON.parse(userStr);
        setIsAuthenticated(true);
        setIsAdmin(user.Role === 'admin');
        
        // 如果是管理员专属路由，但用户不是管理员
        if (adminOnly && user.Role !== 'admin') {
          navigate('/dashboard'); // 重定向到普通用户仪表盘
        }
      } catch (error) {
        console.error('无效的用户信息', error);
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        navigate('/login');
      }
    }
    
    setIsLoading(false);
  }, [navigate, location.pathname, adminOnly]);

  if (isLoading) {
    return <div className="loading">加载中...</div>;
  }

  return isAuthenticated && (!adminOnly || isAdmin) ? children : null;
}

function App() {
  return (
    <Router>
      <div className="App">
        {/* 在所有页面上显示Header */}
        <Header />
        <div className="app-content">
          <Routes>
            {/* 公共路由 */}
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route path="/verify-email" element={<VerifyEmail />} />
            
            {/* 受保护的普通用户路由 */}
            <Route 
              path="/dashboard" 
              element={
                <ProtectedRoute>
                  <Dashboard />
                </ProtectedRoute>
              } 
            />
            
            <Route 
              path="/food-analysis" 
              element={
                <ProtectedRoute>
                  <FoodAnalysis />
                </ProtectedRoute>
              } 
            />
            
            <Route 
              path="/food-records" 
              element={
                <ProtectedRoute>
                  <FoodRecords />
                </ProtectedRoute>
              } 
            />
            
            <Route 
              path="/food-records/edit/:id" 
              element={
                <ProtectedRoute>
                  <FoodRecordEdit />
                </ProtectedRoute>
              } 
            />
            
            {/* 健康分析页面 */}
            <Route 
              path="/health-analysis" 
              element={
                <ProtectedRoute>
                  <HealthAnalysis />
                </ProtectedRoute>
              } 
            />
            
            {/* 受保护的管理员路由 */}
            <Route 
              path="/admin/dashboard" 
              element={
                <ProtectedRoute adminOnly={true}>
                  <AdminDashboard />
                </ProtectedRoute>
              } 
            />
            
            {/* 默认重定向到登录页面 */}
            <Route path="/" element={<Navigate to="/login" />} />
            
            {/* 捕获所有未匹配的路由 */}
            <Route path="*" element={<Navigate to="/login" />} />
          </Routes>
        </div>
      </div>
    </Router>
  );
}

export default App;
