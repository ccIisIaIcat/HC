import React, { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import '../styles/FoodAnalysis.css';
import api from '../api';

function FoodAnalysis() {
  const [image, setImage] = useState(null);
  const [preview, setPreview] = useState(null);
  const [dragging, setDragging] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [result, setResult] = useState(null);
  const [notes, setNotes] = useState('');
  const [mealType, setMealType] = useState('午餐'); // 默认午餐
  const [imageDescription, setImageDescription] = useState(''); // 添加图片描述状态
  const [editableResult, setEditableResult] = useState(null); // 可编辑的结果数据
  const [originalWeight, setOriginalWeight] = useState(0); // 原始重量，用于比例计算
  const fileInputRef = useRef(null);
  const navigate = useNavigate();

  // 当分析结果更新时，创建一个可编辑的副本
  useEffect(() => {
    if (result) {
      setEditableResult(JSON.parse(JSON.stringify(result))); // 深拷贝
      setOriginalWeight(result.weight || 100);
    }
  }, [result]);

  // 处理文件选择
  const handleFileSelect = (file) => {
    if (!file) return;

    // 验证文件类型
    if (!file.type.match('image/jpeg') && !file.type.match('image/png')) {
      setError('只支持JPG和PNG格式图片');
      return;
    }

    // 验证文件大小 (最大10MB)
    if (file.size > 10 * 1024 * 1024) {
      setError('图片大小不能超过10MB');
      return;
    }

    // 清除错误信息
    setError('');
    setImage(file);

    // 创建预览
    const reader = new FileReader();
    reader.onload = (e) => {
      setPreview(e.target.result);
    };
    reader.readAsDataURL(file);
  };

  // 处理文件点击选择
  const handleFileInputChange = (e) => {
    handleFileSelect(e.target.files[0]);
  };

  // 处理拖拽事件
  const handleDragOver = (e) => {
    e.preventDefault();
    setDragging(true);
  };

  const handleDragLeave = (e) => {
    e.preventDefault();
    setDragging(false);
  };

  const handleDrop = (e) => {
    e.preventDefault();
    setDragging(false);
    
    if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
      handleFileSelect(e.dataTransfer.files[0]);
    }
  };

  // 分析食物图片
  const handleAnalyze = async () => {
    if (!image) {
      setError('请先上传食物图片');
      return;
    }

    setLoading(true);
    setError('');
    
    try {
      const formData = new FormData();
      formData.append('image', image);
      // 添加图片描述（如果有）
      if (imageDescription.trim()) {
        formData.append('image_description', imageDescription);
      }
      
      const response = await api.food.analyzeFoodImage(formData);
      
      if (response && response.hasFood) {
        setResult(response);
      } else {
        setError('未能检测到食物，请尝试其他图片');
      }
    } catch (err) {
      console.error('分析失败:', err);
      setError('食物分析失败，请稍后再试');
    } finally {
      setLoading(false);
    }
  };

  // 返回仪表盘
  const handleBack = () => {
    navigate('/dashboard');
  };

  // 清除图片和结果
  const handleReset = () => {
    setImage(null);
    setPreview(null);
    setResult(null);
    setError('');
    setImageDescription(''); // 清除图片描述
  };

  // 处理营养成分值变化
  const handleNutrientChange = (category, key, subKey = null, value) => {
    if (!editableResult) return;

    // 创建新的状态对象
    const newResult = { ...editableResult };
    
    // 转换为数字
    const numValue = parseFloat(value) || 0;
    
    if (subKey) {
      // 嵌套属性，如维生素B
      if (subKey.includes('.')) {
        const [parentKey, childKey] = subKey.split('.');
        newResult.nutrition[category][key][parentKey][childKey] = numValue;
      } else {
        newResult.nutrition[category][key][subKey] = numValue;
      }
    } else if (category) {
      // 类别内的属性，如维生素、矿物质等
      newResult.nutrition[category][key] = numValue;
    } else {
      // 基本营养素
      newResult.nutrition[key] = numValue;
    }
    
    setEditableResult(newResult);
  };

  // 处理重量变化并按比例调整其他营养成分
  const handleWeightChange = (value) => {
    if (!editableResult || !originalWeight) return;
    
    const newWeight = parseFloat(value) || 0;
    if (newWeight <= 0) return;
    
    // 计算比例因子
    const ratio = newWeight / editableResult.weight;
    
    // 创建新的状态对象
    const newResult = JSON.parse(JSON.stringify(editableResult));
    newResult.weight = newWeight;
    
    // 按比例调整基本营养素
    const basicNutrients = ['calories', 'protein', 'totalFat', 'saturatedFat', 'transFat', 
                           'unsaturatedFat', 'carbohydrates', 'sugar', 'fiber'];
    
    basicNutrients.forEach(nutrient => {
      if (typeof newResult.nutrition[nutrient] === 'number') {
        newResult.nutrition[nutrient] = +(newResult.nutrition[nutrient] * ratio).toFixed(2);
      }
    });
    
    // 按比例调整维生素
    for (const key in newResult.nutrition.vitamins) {
      if (key === 'vitaminB') {
        for (const bKey in newResult.nutrition.vitamins.vitaminB) {
          newResult.nutrition.vitamins.vitaminB[bKey] = +(newResult.nutrition.vitamins.vitaminB[bKey] * ratio).toFixed(2);
        }
      } else if (typeof newResult.nutrition.vitamins[key] === 'number') {
        newResult.nutrition.vitamins[key] = +(newResult.nutrition.vitamins[key] * ratio).toFixed(2);
      }
    }
    
    // 按比例调整矿物质
    for (const key in newResult.nutrition.minerals) {
      if (typeof newResult.nutrition.minerals[key] === 'number') {
        newResult.nutrition.minerals[key] = +(newResult.nutrition.minerals[key] * ratio).toFixed(2);
      }
    }
    
    setEditableResult(newResult);
  };

  // 记录食物（使用编辑后的数据）
  const handleSaveRecord = async () => {
    if (!editableResult) {
      setError('没有分析结果可保存');
      return;
    }
    
    try {
      setLoading(true);
      
      // 准备食物记录数据
      const recordData = {
        // 基本信息
        food_name: editableResult.foodType || '未知食物',
        weight: editableResult.weight || 100, // 使用编辑后的重量
        record_time: new Date().toISOString(),
        meal_type: mealType || '午餐',
        notes: notes || '',

        // 基本营养素
        calories: editableResult.nutrition?.calories || 0,
        protein: editableResult.nutrition?.protein || 0,
        total_fat: editableResult.nutrition?.totalFat || 0,
        saturated_fat: editableResult.nutrition?.saturatedFat || 0,
        trans_fat: editableResult.nutrition?.transFat || 0,
        unsaturated_fat: editableResult.nutrition?.unsaturatedFat || 0,
        carbohydrates: editableResult.nutrition?.carbohydrates || 0,
        sugar: editableResult.nutrition?.sugar || 0,
        fiber: editableResult.nutrition?.fiber || 0,

        // 维生素
        vitamin_a: editableResult.nutrition?.vitamins?.vitaminA || 0,
        vitamin_c: editableResult.nutrition?.vitamins?.vitaminC || 0,
        vitamin_d: editableResult.nutrition?.vitamins?.vitaminD || 0,
        vitamin_b1: editableResult.nutrition?.vitamins?.vitaminB?.b1 || 0,
        vitamin_b2: editableResult.nutrition?.vitamins?.vitaminB?.b2 || 0,

        // 矿物质
        calcium: editableResult.nutrition?.minerals?.calcium || 0,
        iron: editableResult.nutrition?.minerals?.iron || 0,
        sodium: editableResult.nutrition?.minerals?.sodium || 0,
        potassium: editableResult.nutrition?.minerals?.potassium || 0
      };

      console.log('保存食物记录数据:', recordData);

      // 调用API保存记录
      const response = await api.food.createFoodRecord(recordData);
      
      console.log('保存食物记录响应:', response);
      
      // 显示成功消息
      alert(`食物记录保存成功！\n食物: ${recordData.food_name}\n餐食类型: ${recordData.meal_type}`);
      
      // 清除结果和笔记
      setResult(null);
      setEditableResult(null);
      setPreview(null);
      setImage(null);
      setNotes('');
    } catch (err) {
      console.error('保存食物记录失败:', err);
      setError(`保存食物记录失败: ${err.message || '未知错误'}`);
      alert('保存食物记录失败，请查看控制台获取详细信息');
    } finally {
      setLoading(false);
    }
  };

  // 放弃记录
  const handleDiscardRecord = () => {
    // 确认是否放弃记录
    if (window.confirm('确定要放弃当前的食物记录吗？分析结果将被清除，但图片将保留。')) {
      // 清除结果和笔记，但保留图片
      setResult(null);
      setEditableResult(null);
      setNotes('');
      setError(''); // 清除可能存在的错误信息
      alert('已放弃食物记录');
    }
  };

  // 直接分析并保存记录（替代方法，使用analyze-and-save接口）
  const handleAnalyzeAndSaveRecord = async () => {
    if (!image || !mealType) return;
    
    try {
      setLoading(true);
      setError('');
      
      const formData = new FormData();
      formData.append('image', image);
      formData.append('meal_type', mealType);
      formData.append('notes', notes);
      // 添加图片描述（如果有）
      if (imageDescription.trim()) {
        formData.append('image_description', imageDescription);
      }
      
      const response = await api.food.analyzeAndSaveFoodRecord(formData);
      
      alert('食物分析和记录保存成功！');
      
      // 清除状态
      setImage(null);
      setPreview(null);
      setResult(null);
      setNotes('');
      setImageDescription(''); // 清除图片描述
    } catch (err) {
      console.error('分析并保存失败:', err);
      setError('食物分析和保存失败，请稍后再试');
    } finally {
      setLoading(false);
    }
  };

  // 创建可编辑的营养素输入框
  const renderEditableNutrientField = (label, value, unit, onChange) => {
    return (
      <div className="nutrition-item editable">
        <span className="nutrition-label">{label}</span>
        <div className="nutrition-value-input">
          <input 
            type="number" 
            value={value} 
            onChange={(e) => onChange(e.target.value)}
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
      <div className="food-analysis-main">
        <div className="food-analysis-header">
          <h2>食物营养分析</h2>
          <button className="back-button" onClick={handleBack}>返回仪表盘</button>
        </div>
        
        <div className="food-analysis-content">
          <div className="upload-section">
            <h3>上传食物图片</h3>
            {error && <div className="error-message">{error}</div>}
            
            <div 
              className={`upload-area ${dragging ? 'dragging' : ''} ${loading ? 'disabled' : ''}`}
              onClick={() => !loading && fileInputRef.current.click()}
              onDragOver={!loading ? handleDragOver : (e) => e.preventDefault()}
              onDragLeave={!loading ? handleDragLeave : (e) => e.preventDefault()}
              onDrop={!loading ? handleDrop : (e) => e.preventDefault()}
            >
              {preview ? (
                <div className="preview-container">
                  <img src={preview} alt="食物预览" className="image-preview" />
                  <div className="preview-overlay">
                    <span>{loading ? '处理中...' : '点击或拖拽更换图片'}</span>
                  </div>
                </div>
              ) : (
                <div className="upload-placeholder">
                  <i className="upload-icon">📷</i>
                  <p>{loading ? '处理中...' : '点击或拖拽上传食物图片'}</p>
                  <span>支持JPG和PNG格式，最大10MB</span>
                </div>
              )}
              <input 
                type="file" 
                ref={fileInputRef} 
                className="file-input" 
                accept="image/jpeg,image/png" 
                onChange={handleFileInputChange}
                disabled={loading}
              />
            </div>
            
            <div className="action-buttons">
              <button 
                className="analyze-button" 
                onClick={handleAnalyze} 
                disabled={!image || loading}
              >
                {loading ? '分析中...' : '分析图片'}
              </button>
              {preview && (
                <button className="reset-button" onClick={handleReset}>
                  清除
                </button>
              )}
            </div>

            {/* 添加图片描述输入框 */}
            {preview && (
              <div className="description-input-container">
                <label htmlFor="image-description">食物描述 (可选):</label>
                <textarea
                  id="image-description"
                  placeholder="请描述图片中的食物，例如：这是一盘炒饭，里面有鸡蛋、胡萝卜和青豆...,最好增加对食物大小或重量的描述"
                  value={imageDescription}
                  onChange={(e) => setImageDescription(e.target.value)}
                  disabled={loading}
                  rows="2"
                  className="image-description-input"
                />
                <small className="description-help-text">
                  提供准确的食物描述可以帮助AI更精确地分析食物营养成分
                </small>
              </div>
            )}
          </div>

          {editableResult && (
            <div className="result-section">
              <h3>分析结果 <small>(值可编辑)</small></h3>
              <div className="result-card">
                <div className="result-header">
                  <h4>检测到的食物: 
                    <input 
                      type="text" 
                      value={editableResult.foodType}
                      onChange={(e) => {
                        const newResult = { ...editableResult };
                        newResult.foodType = e.target.value;
                        setEditableResult(newResult);
                      }}
                      className="food-name-input"
                      placeholder="输入食物名称"
                    />
                  </h4>
                  <div className="food-weight-edit">
                    <input 
                      type="number" 
                      value={editableResult.weight}
                      onChange={(e) => handleWeightChange(e.target.value)} 
                      min="1"
                      step="1"
                      className="weight-input"
                    /> 克
                  </div>
                </div>
                
                <div className="nutrition-info">
                  <div className="nutrition-section">
                    <h5>基本营养素</h5>
                    <div className="nutrition-grid">
                      {renderEditableNutrientField('热量', editableResult.nutrition.calories, '卡路里', 
                        (value) => handleNutrientChange(null, 'calories', null, value))}
                          
                      {renderEditableNutrientField('蛋白质', editableResult.nutrition.protein, '克', 
                        (value) => handleNutrientChange(null, 'protein', null, value))}
                          
                      {renderEditableNutrientField('总脂肪', editableResult.nutrition.totalFat, '克', 
                        (value) => handleNutrientChange(null, 'totalFat', null, value))}
                          
                      {renderEditableNutrientField('碳水化合物', editableResult.nutrition.carbohydrates, '克', 
                        (value) => handleNutrientChange(null, 'carbohydrates', null, value))}
                          
                      {renderEditableNutrientField('糖', editableResult.nutrition.sugar, '克', 
                        (value) => handleNutrientChange(null, 'sugar', null, value))}
                          
                      {renderEditableNutrientField('纤维', editableResult.nutrition.fiber, '克', 
                        (value) => handleNutrientChange(null, 'fiber', null, value))}
                    </div>
                  </div>
                  
                  <div className="nutrition-section">
                    <h5>脂肪详情</h5>
                    <div className="nutrition-grid">
                      {renderEditableNutrientField('饱和脂肪', editableResult.nutrition.saturatedFat, '克', 
                        (value) => handleNutrientChange(null, 'saturatedFat', null, value))}
                          
                      {renderEditableNutrientField('反式脂肪', editableResult.nutrition.transFat, '克', 
                        (value) => handleNutrientChange(null, 'transFat', null, value))}
                          
                      {renderEditableNutrientField('不饱和脂肪', editableResult.nutrition.unsaturatedFat, '克', 
                        (value) => handleNutrientChange(null, 'unsaturatedFat', null, value))}
                    </div>
                  </div>
                  
                  <div className="nutrition-section">
                    <h5>维生素</h5>
                    <div className="nutrition-grid">
                      {renderEditableNutrientField('维生素A', editableResult.nutrition.vitamins.vitaminA, '微克', 
                        (value) => handleNutrientChange('vitamins', 'vitaminA', null, value))}
                          
                      {renderEditableNutrientField('维生素C', editableResult.nutrition.vitamins.vitaminC, '毫克', 
                        (value) => handleNutrientChange('vitamins', 'vitaminC', null, value))}
                          
                      {renderEditableNutrientField('维生素D', editableResult.nutrition.vitamins.vitaminD, '微克', 
                        (value) => handleNutrientChange('vitamins', 'vitaminD', null, value))}
                          
                      {renderEditableNutrientField('维生素E', editableResult.nutrition.vitamins.vitaminE, '毫克', 
                        (value) => handleNutrientChange('vitamins', 'vitaminE', null, value))}
                          
                      {renderEditableNutrientField('维生素K', editableResult.nutrition.vitamins.vitaminK, '微克', 
                        (value) => handleNutrientChange('vitamins', 'vitaminK', null, value))}
                          
                      {renderEditableNutrientField('维生素B1', editableResult.nutrition.vitamins.vitaminB.b1, '毫克', 
                        (value) => handleNutrientChange('vitamins', 'vitaminB', 'b1', value))}
                          
                      {renderEditableNutrientField('维生素B2', editableResult.nutrition.vitamins.vitaminB.b2, '毫克', 
                        (value) => handleNutrientChange('vitamins', 'vitaminB', 'b2', value))}
                          
                      {renderEditableNutrientField('维生素B6', editableResult.nutrition.vitamins.vitaminB.b6, '毫克', 
                        (value) => handleNutrientChange('vitamins', 'vitaminB', 'b6', value))}
                          
                      {renderEditableNutrientField('维生素B12', editableResult.nutrition.vitamins.vitaminB.b12, '微克', 
                        (value) => handleNutrientChange('vitamins', 'vitaminB', 'b12', value))}
                    </div>
                  </div>
                  
                  <div className="nutrition-section">
                    <h5>矿物质</h5>
                    <div className="nutrition-grid">
                      {renderEditableNutrientField('钙', editableResult.nutrition.minerals.calcium, '毫克', 
                        (value) => handleNutrientChange('minerals', 'calcium', null, value))}
                          
                      {renderEditableNutrientField('铁', editableResult.nutrition.minerals.iron, '毫克', 
                        (value) => handleNutrientChange('minerals', 'iron', null, value))}
                          
                      {renderEditableNutrientField('钠', editableResult.nutrition.minerals.sodium, '毫克', 
                        (value) => handleNutrientChange('minerals', 'sodium', null, value))}
                          
                      {renderEditableNutrientField('钾', editableResult.nutrition.minerals.potassium, '毫克', 
                        (value) => handleNutrientChange('minerals', 'potassium', null, value))}
                          
                      {renderEditableNutrientField('锌', editableResult.nutrition.minerals.zinc, '毫克', 
                        (value) => handleNutrientChange('minerals', 'zinc', null, value))}
                          
                      {renderEditableNutrientField('镁', editableResult.nutrition.minerals.magnesium, '毫克', 
                        (value) => handleNutrientChange('minerals', 'magnesium', null, value))}
                    </div>
                  </div>
                </div>
                
                {/* 添加记录表单 */}
                <div className="record-form">
                  <h5>保存食物记录</h5>
                  
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
                      className="save-button"
                      onClick={handleSaveRecord}
                      disabled={loading}
                    >
                      {loading ? '保存中...' : '记录'}
                    </button>
                    <button 
                      className="discard-button"
                      onClick={handleDiscardRecord}
                      disabled={loading}
                    >
                      放弃
                    </button>
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default FoodAnalysis; 