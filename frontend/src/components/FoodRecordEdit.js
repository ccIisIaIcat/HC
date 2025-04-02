import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import api from '../api';
import '../styles/FoodAnalysis.css';
import '../styles/FoodRecordEdit.css';

function FoodRecordEdit() {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [record, setRecord] = useState(null);
  const [editableRecord, setEditableRecord] = useState(null);
  const [originalWeight, setOriginalWeight] = useState(0);
  const [notes, setNotes] = useState('');
  const [mealType, setMealType] = useState('午餐');
  
  const navigate = useNavigate();
  const { id } = useParams(); // 从URL获取记录ID
  
  // 加载食物记录数据
  useEffect(() => {
    fetchFoodRecord();
  }, [id]);
  
  // 当记录加载完成后，复制到可编辑状态
  useEffect(() => {
    if (record) {
      setEditableRecord(JSON.parse(JSON.stringify(record)));
      setOriginalWeight(record.weight || 100);
      setNotes(record.notes || '');
      setMealType(record.meal_type || '午餐');
    }
  }, [record]);
  
  // 获取食物记录详情
  const fetchFoodRecord = async () => {
    if (!id) {
      setError('无效的记录ID');
      setLoading(false);
      return;
    }
    
    try {
      setLoading(true);
      setError('');
      
      const response = await api.food.getFoodRecord(id);
      if (response && response.record) {
        setRecord(response.record);
      } else {
        setError('未找到记录数据');
      }
    } catch (err) {
      console.error('获取食物记录失败:', err);
      setError('获取食物记录失败: ' + (err.message || '未知错误'));
    } finally {
      setLoading(false);
    }
  };
  
  // 处理返回按钮
  const handleBack = () => {
    navigate('/food-records');
  };
  
  // 处理营养成分值变化
  const handleNutrientChange = (key, value) => {
    if (!editableRecord) return;
    
    // 创建新的状态对象
    const newRecord = { ...editableRecord };
    
    // 转换为数字
    const numValue = parseFloat(value) || 0;
    
    // 更新营养成分值
    newRecord[key] = numValue;
    
    setEditableRecord(newRecord);
  };
  
  // 处理重量变化并按比例调整其他营养成分
  const handleWeightChange = (value) => {
    if (!editableRecord || !originalWeight) return;
    
    const newWeight = parseFloat(value) || 0;
    if (newWeight <= 0) return;
    
    // 计算比例因子
    const ratio = newWeight / editableRecord.weight;
    
    // 创建新的状态对象
    const newRecord = { ...editableRecord };
    newRecord.weight = newWeight;
    
    // 按比例调整基本营养素
    const nutrientFields = [
      'calories', 'protein', 'total_fat', 'saturated_fat', 
      'trans_fat', 'unsaturated_fat', 'carbohydrates', 
      'sugar', 'fiber', 'vitamin_a', 'vitamin_c', 
      'vitamin_d', 'vitamin_b1', 'vitamin_b2', 'calcium', 
      'iron', 'sodium', 'potassium'
    ];
    
    nutrientFields.forEach(field => {
      if (typeof newRecord[field] === 'number') {
        newRecord[field] = +(newRecord[field] * ratio).toFixed(2);
      }
    });
    
    setEditableRecord(newRecord);
  };
  
  // 保存更新的记录
  const handleSave = async () => {
    if (!editableRecord) {
      setError('没有记录数据可保存');
      return;
    }
    
    try {
      setLoading(true);
      
      // 准备更新的数据
      const updateData = {
        ...editableRecord,
        meal_type: mealType,
        notes: notes
      };
      
      // 调用API更新记录
      await api.food.updateFoodRecord(id, updateData);
      
      alert('食物记录更新成功！');
      
      // 返回到记录列表
      navigate('/food-records');
    } catch (err) {
      console.error('更新食物记录失败:', err);
      setError('更新食物记录失败: ' + (err.message || '未知错误'));
      alert('更新记录失败，请查看错误信息');
    } finally {
      setLoading(false);
    }
  };
  
  // 创建可编辑的营养素输入框
  const renderEditableNutrientField = (label, fieldKey, value, unit) => {
    return (
      <div className="nutrition-item editable">
        <span className="nutrition-label">{label}</span>
        <div className="nutrition-value-input">
          <input 
            type="number" 
            value={value} 
            onChange={(e) => handleNutrientChange(fieldKey, e.target.value)}
            min="0"
            step="0.1"
            className="nutrition-input"
          />
          <span className="nutrition-unit">{unit}</span>
        </div>
      </div>
    );
  };
  
  return (
    <div className="food-analysis-container">
      <div className="food-analysis-header">
        <button className="back-button" onClick={handleBack}>返回列表</button>
        <h1>编辑食物记录</h1>
      </div>
      
      <div className="food-analysis-content">
        {loading && !editableRecord && (
          <div className="loading-indicator">正在加载数据...</div>
        )}
        
        {error && <div className="error-message">{error}</div>}
        
        {editableRecord && (
          <div className="result-section">
            <h2>食物记录 <small>(值可编辑)</small></h2>
            <div className="result-card">
              <div className="result-header">
                <h3>食物: 
                  <input 
                    type="text" 
                    value={editableRecord.food_name}
                    onChange={(e) => {
                      const newRecord = { ...editableRecord };
                      newRecord.food_name = e.target.value;
                      setEditableRecord(newRecord);
                    }}
                    className="food-name-input"
                    placeholder="输入食物名称"
                  />
                </h3>
                <div className="food-weight-edit">
                  <input 
                    type="number" 
                    value={editableRecord.weight}
                    onChange={(e) => handleWeightChange(e.target.value)} 
                    min="1"
                    step="1"
                    className="weight-input"
                  /> 克
                </div>
              </div>
              
              <div className="nutrition-info">
                <div className="nutrition-section">
                  <h4>基本营养素</h4>
                  <div className="nutrition-grid">
                    {renderEditableNutrientField('热量', 'calories', editableRecord.calories, '卡路里')}
                    {renderEditableNutrientField('蛋白质', 'protein', editableRecord.protein, '克')}
                    {renderEditableNutrientField('总脂肪', 'total_fat', editableRecord.total_fat, '克')}
                    {renderEditableNutrientField('碳水化合物', 'carbohydrates', editableRecord.carbohydrates, '克')}
                    {renderEditableNutrientField('糖', 'sugar', editableRecord.sugar, '克')}
                    {renderEditableNutrientField('纤维', 'fiber', editableRecord.fiber, '克')}
                  </div>
                </div>
                
                <div className="nutrition-section">
                  <h4>脂肪详情</h4>
                  <div className="nutrition-grid">
                    {renderEditableNutrientField('饱和脂肪', 'saturated_fat', editableRecord.saturated_fat, '克')}
                    {renderEditableNutrientField('反式脂肪', 'trans_fat', editableRecord.trans_fat, '克')}
                    {renderEditableNutrientField('不饱和脂肪', 'unsaturated_fat', editableRecord.unsaturated_fat, '克')}
                  </div>
                </div>
                
                <div className="nutrition-section">
                  <h4>维生素</h4>
                  <div className="nutrition-grid">
                    {renderEditableNutrientField('维生素A', 'vitamin_a', editableRecord.vitamin_a, '微克')}
                    {renderEditableNutrientField('维生素C', 'vitamin_c', editableRecord.vitamin_c, '毫克')}
                    {renderEditableNutrientField('维生素D', 'vitamin_d', editableRecord.vitamin_d, '微克')}
                    {renderEditableNutrientField('维生素B1', 'vitamin_b1', editableRecord.vitamin_b1, '毫克')}
                    {renderEditableNutrientField('维生素B2', 'vitamin_b2', editableRecord.vitamin_b2, '毫克')}
                  </div>
                </div>
                
                <div className="nutrition-section">
                  <h4>矿物质</h4>
                  <div className="nutrition-grid">
                    {renderEditableNutrientField('钙', 'calcium', editableRecord.calcium, '毫克')}
                    {renderEditableNutrientField('铁', 'iron', editableRecord.iron, '毫克')}
                    {renderEditableNutrientField('钠', 'sodium', editableRecord.sodium, '毫克')}
                    {renderEditableNutrientField('钾', 'potassium', editableRecord.potassium, '毫克')}
                  </div>
                </div>
                
                {/* 添加记录表单 */}
                <div className="record-form">
                  <h4>记录详情</h4>
                  
                  <div className="form-group">
                    <label>餐食类型:</label>
                    <select 
                      value={mealType} 
                      onChange={(e) => setMealType(e.target.value)}
                      className="meal-type-select"
                      disabled={loading}
                    >
                      <option value="早餐">早餐</option>
                      <option value="午餐">午餐</option>
                      <option value="晚餐">晚餐</option>
                      <option value="加餐">加餐</option>
                    </select>
                  </div>
                  
                  <div className="form-group">
                    <label>备注:</label>
                    <textarea 
                      value={notes}
                      onChange={(e) => setNotes(e.target.value)}
                      placeholder="添加关于这顿饭的备注..."
                      className="notes-textarea"
                      rows="3"
                      disabled={loading}
                    />
                  </div>
                  
                  <div className="record-buttons">
                    <button 
                      className="save-button wide-button"
                      onClick={handleSave}
                      disabled={loading}
                    >
                      {loading ? '保存中...' : '保存修改'}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default FoodRecordEdit; 