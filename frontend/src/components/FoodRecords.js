import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../api';
import '../styles/FoodRecords.css';

function FoodRecords() {
  const [records, setRecords] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalRecords, setTotalRecords] = useState(0);
  const recordsPerPage = 10; // 每页显示10条记录
  
  const navigate = useNavigate();

  // 获取食物记录列表
  useEffect(() => {
    fetchFoodRecords();
  }, [currentPage]);

  // 获取食物记录的函数
  const fetchFoodRecords = async () => {
    try {
      setLoading(true);
      setError('');
      
      // 计算六个月前的日期
      const today = new Date();
      const sixMonthsAgo = new Date();
      sixMonthsAgo.setMonth(today.getMonth() - 6);
      
      // 格式化为YYYY-MM-DD
      const startDate = sixMonthsAgo.toISOString().split('T')[0];
      const endDate = today.toISOString().split('T')[0];
      
      // 从API获取数据
      const response = await api.food.getFoodRecords(startDate, endDate);
      
      // 判断响应数据格式并提取记录数组
      let recordsArray = [];
      if (response && response.records && Array.isArray(response.records)) {
        recordsArray = response.records;
      } else if (response && Array.isArray(response)) {
        recordsArray = response;
      }
      
      // 计算总页数
      const totalItems = (response && response.total) || recordsArray.length;
      setTotalRecords(totalItems);
      setTotalPages(Math.ceil(totalItems / recordsPerPage));
      
      // 获取当前页的数据
      const startIndex = (currentPage - 1) * recordsPerPage;
      const endIndex = startIndex + recordsPerPage;
      const currentPageData = recordsArray.slice(startIndex, endIndex);
      
      setRecords(currentPageData);
    } catch (err) {
      console.error('获取食物记录失败:', err);
      setError('获取食物记录失败，请稍后再试');
    } finally {
      setLoading(false);
    }
  };

  // 格式化日期时间
  const formatDateTime = (dateString) => {
    if (!dateString) return '未知时间';
    
    try {
      const date = new Date(dateString);
      
      // 检查日期是否有效
      if (isNaN(date.getTime())) {
        return '日期无效';
      }
      
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      });
    } catch (error) {
      console.error('日期格式化错误:', error);
      return '格式错误';
    }
  };
  
  // 处理页面导航
  const handlePageChange = (newPage) => {
    if (newPage >= 1 && newPage <= totalPages) {
      setCurrentPage(newPage);
    }
  };

  // 返回仪表盘
  const handleBack = () => {
    navigate('/dashboard');
  };
  
  // 编辑食物记录
  const handleEditRecord = (id) => {
    if (!id) {
      alert('无法编辑：记录ID不存在');
      return;
    }
    navigate(`/food-records/edit/${id}`);
  };
  
  // 删除食物记录
  const handleDeleteRecord = async (id) => {
    if (!id) {
      alert('无法删除：记录ID不存在');
      return;
    }
    
    if (window.confirm('确定要删除这条食物记录吗？此操作不可撤销。')) {
      try {
        setLoading(true);
        await api.food.deleteFoodRecord(id);
        alert('记录已成功删除');
        // 重新获取数据
        fetchFoodRecords();
      } catch (err) {
        console.error('删除记录失败:', err);
        alert('删除记录失败，请稍后再试');
      } finally {
        setLoading(false);
      }
    }
  };
  
  // 添加新记录（跳转到食物分析页面）
  const handleAddRecord = () => {
    navigate('/food-analysis');
  };

  return (
    <div className="food-records-container">
      <div className="food-records-content">
        <div className="records-header">
          <h2>我的饮食记录</h2>
          <div className="records-actions">
            <button className="back-button" onClick={handleBack}>返回仪表盘</button>
            <button className="add-record-button" onClick={handleAddRecord}>
              添加记录
            </button>
          </div>
        </div>
        
        {loading && <div className="loading-indicator">正在加载数据...</div>}
        
        {error && <div className="error-message">{error}</div>}
        
        {!loading && !error && records.length === 0 && (
          <div className="empty-records">
            <p>暂无饮食记录</p>
            <button className="add-record-button" onClick={handleAddRecord}>
              开始记录饮食
            </button>
          </div>
        )}
        
        {!loading && !error && records.length > 0 && (
          <>
            <div className="records-summary">
              共 {totalRecords} 条记录，显示 {(currentPage - 1) * recordsPerPage + 1} - {Math.min(currentPage * recordsPerPage, totalRecords)} 条
            </div>
            
            <div className="records-list">
              {records.map((record, index) => (
                <div key={record.ID || record.id || `record-${index}`} className="record-item">
                  <div className="record-time">
                    {formatDateTime(record.record_time)}
                  </div>
                  <div className="record-info">
                    <div className="record-food-name">{record.food_name}</div>
                    <div className="record-meal-type">{record.meal_type}</div>
                  </div>
                  <div className="record-calories">
                    {(record.calories || 0).toFixed(0)} 卡路里
                  </div>
                  <div className="record-actions">
                    <button 
                      className="edit-button" 
                      title="编辑"
                      onClick={() => handleEditRecord(record.ID || record.id)}
                    >
                      ✏️
                    </button>
                    <button 
                      className="delete-button" 
                      title="删除"
                      onClick={() => handleDeleteRecord(record.ID || record.id)}
                    >
                      🗑️
                    </button>
                  </div>
                </div>
              ))}
            </div>
            
            <div className="pagination">
              <button 
                className="pagination-button" 
                disabled={currentPage === 1}
                onClick={() => handlePageChange(1)}
              >
                首页
              </button>
              <button 
                className="pagination-button" 
                disabled={currentPage === 1}
                onClick={() => handlePageChange(currentPage - 1)}
              >
                上一页
              </button>
              <span className="pagination-info">
                {currentPage} / {totalPages}
              </span>
              <button 
                className="pagination-button" 
                disabled={currentPage === totalPages}
                onClick={() => handlePageChange(currentPage + 1)}
              >
                下一页
              </button>
              <button 
                className="pagination-button" 
                disabled={currentPage === totalPages}
                onClick={() => handlePageChange(totalPages)}
              >
                末页
              </button>
            </div>
          </>
        )}
      </div>
    </div>
  );
}

export default FoodRecords; 