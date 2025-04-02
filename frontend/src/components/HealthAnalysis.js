import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../api';
import '../styles/HealthAnalysis.css';

function HealthAnalysis() {
  const navigate = useNavigate();
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [description, setDescription] = useState('');
  const [analysisType, setAnalysisType] = useState('comprehensive'); // 默认全面分析
  const [analysisResult, setAnalysisResult] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [analysisDate, setAnalysisDate] = useState('');

  // 返回仪表盘
  const handleBack = () => {
    navigate('/dashboard');
  };

  // 处理分析请求
  const handleAnalyze = async () => {
    // 重置状态
    setAnalysisResult('');
    setError('');
    setAnalysisDate('');
    
    // 表单验证
    if (!startDate || !endDate) {
      setError('请选择起始和结束日期');
      return;
    }

    if (new Date(startDate) > new Date(endDate)) {
      setError('起始日期不能晚于结束日期');
      return;
    }

    // 准备请求数据
    const analysisData = {
      start_date: startDate,
      end_date: endDate,
      analysis_type: analysisType,
      description: description
    };

    // 发送请求
    setLoading(true);
    try {
      const result = await api.health.analyzeHealth(analysisData);
      if (result && result.analysis) {
        setAnalysisResult(result.analysis);
        // 设置分析日期
        setAnalysisDate(new Date().toLocaleString('zh-CN'));
      } else if (result && result.message) {
        setAnalysisResult(result.message);
        setAnalysisDate(new Date().toLocaleString('zh-CN'));
      } else {
        setAnalysisResult("分析完成，但未返回详细结果。");
        setAnalysisDate(new Date().toLocaleString('zh-CN'));
      }
    } catch (error) {
      console.error('健康分析失败:', error);
      setError(error.response?.data?.error || error.response?.data?.message || '分析请求失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  // 获取分析类型名称
  const getAnalysisTypeName = () => {
    switch (analysisType) {
      case 'comprehensive': return '全面分析';
      case 'nutrition': return '营养均衡分析';
      case 'calories': return '热量摄入分析';
      case 'macros': return '宏量营养素分析';
      default: return '全面分析';
    }
  };

  return (
    <div className="health-analysis-container">
      <div className="health-analysis-content">
        <h2>健康分析</h2>
        <button className="back-button" onClick={handleBack}>返回仪表盘</button>
        
        <div className="health-analysis-section">
          <h3>健康数据分析</h3>
          <p className="section-description">
            根据您的饮食记录，我们可以为您提供个性化的健康分析报告。请选择您希望分析的时间范围和分析类型。
          </p>
          
          <div className="analysis-form">
            <div className="form-group">
              <label htmlFor="start-date">起始日期</label>
              <input 
                type="date" 
                id="start-date" 
                value={startDate}
                onChange={(e) => setStartDate(e.target.value)}
                className="date-input"
              />
            </div>
            
            <div className="form-group">
              <label htmlFor="end-date">结束日期</label>
              <input 
                type="date" 
                id="end-date" 
                value={endDate}
                onChange={(e) => setEndDate(e.target.value)}
                className="date-input"
              />
            </div>
            
            <div className="form-group">
              <label htmlFor="analysis-type">分析类型</label>
              <select 
                id="analysis-type"
                value={analysisType}
                onChange={(e) => setAnalysisType(e.target.value)}
                className="analysis-type-select"
              >
                <option value="comprehensive">全面分析</option>
                <option value="nutrition">营养均衡分析</option>
                <option value="calories">热量摄入分析</option>
                <option value="macros">宏量营养素分析</option>
              </select>
            </div>
            
            <div className="form-group">
              <label htmlFor="description">附加描述 (可选)</label>
              <textarea 
                id="description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="请输入您的健康关注点或饮食习惯，例如：我正在减肥，每天运动30分钟..."
                className="description-textarea"
                rows="4"
              />
            </div>
            
            {error && <div className="error-message">{error}</div>}
            
            <button 
              className="analyze-button" 
              onClick={handleAnalyze} 
              disabled={loading}
            >
              {loading ? '分析中...' : '开始分析'}
            </button>
          </div>

          {/* 加载状态 */}
          {loading && (
            <div className="loading-container">
              <div className="loading-spinner"></div>
              <p>正在分析您的健康数据，请稍候...</p>
            </div>
          )}

          {/* 分析结果显示 */}
          {analysisResult && (
            <div className="analysis-result">
              <h3>您的健康分析报告</h3>
              {analysisDate && <p className="analysis-date">分析时间: {analysisDate}</p>}
              <div className="result-content">
                {typeof analysisResult === 'string' ? 
                  (() => {
                    // 先处理Markdown语法
                    const processedText = analysisResult
                      .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>') // 处理粗体
                      .replace(/\*(.*?)\*/g, '<em>$1</em>') // 处理斜体
                      .replace(/###\s+(.*?)(?:\n|$)/g, '<h4>$1</h4>') // 处理三级标题
                      .replace(/##\s+(.*?)(?:\n|$)/g, '<h3>$1</h3>') // 处理二级标题
                      .replace(/#\s+(.*?)(?:\n|$)/g, '<h2>$1</h2>'); // 处理一级标题
                    
                    // 按空行分割成段落
                    const paragraphs = processedText.split(/\n\s*\n/);
                    
                    return paragraphs.map((paragraph, index) => {
                      // 如果段落是标题，直接渲染HTML
                      if (paragraph.trim().startsWith('<h')) {
                        return <div key={index} dangerouslySetInnerHTML={{ __html: paragraph }} />;
                      }
                      
                      // 处理普通段落，保留单个换行符
                      const processedParagraph = paragraph
                        .replace(/\n/g, '<br />') // 将换行符转换为<br>标签
                        .trim();
                      
                      return processedParagraph ? (
                        <p key={index} dangerouslySetInnerHTML={{ __html: processedParagraph }} />
                      ) : <p key={index}>&nbsp;</p>; // 空行
                    });
                  })()
                 : (
                  <p>分析完成，但无法显示详细结果</p>
                )}
              </div>
            </div>
          )}
        </div>
        
        <div className="health-analysis-section">
          <h2>分析说明</h2>
          <div className="analysis-info-card">
            <div className="info-item">
              <div className="info-icon">📊</div>
              <div className="info-content">
                <h3>全面分析</h3>
                <p>对所选时间范围内的所有饮食数据进行全面评估，包括营养平衡、热量摄入、宏量营养素等各个方面。</p>
              </div>
            </div>
            
            <div className="info-item">
              <div className="info-icon">🍎</div>
              <div className="info-content">
                <h3>营养均衡分析</h3>
                <p>重点分析您的饮食中各种营养素的摄入是否均衡，包括维生素、矿物质等微量元素的摄入情况。</p>
              </div>
            </div>
            
            <div className="info-item">
              <div className="info-icon">🔥</div>
              <div className="info-content">
                <h3>热量摄入分析</h3>
                <p>专注于分析您的热量摄入情况，帮助您了解自己的能量平衡状态，适合关注体重管理的用户。</p>
              </div>
            </div>
            
            <div className="info-item">
              <div className="info-icon">🥩</div>
              <div className="info-content">
                <h3>宏量营养素分析</h3>
                <p>分析您的蛋白质、脂肪和碳水化合物的摄入比例，帮助您调整饮食结构，达到理想的宏量营养素分配。</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default HealthAnalysis; 