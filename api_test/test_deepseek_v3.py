#!/usr/bin/env python
# -*- coding: utf-8 -*-

import requests
import json
import os
import argparse
import logging
import base64
from typing import Optional, Dict, Any

# 配置日志
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

# 全局变量
API_URL = "https://api.deepseek.com/v1/chat/completions"  # API 端点URL
API_KEY = "sk-2ff5bcd3d6d046e8aae5bc566ff5e968"  # API密钥
MODEL = "deepseek-chat"  # 默认模型
DEFAULT_PROMPT = "你好"  # 默认提示词
SYSTEM_PROMPT = "你是由DeepSeek开发的AI助手，可以提供简洁、准确、有用的回答。"  # 系统角色提示

# 食物分析相关
FOOD_ANALYSIS_PROMPT = """分析这张图片中的食物。请提供以下信息，必须严格按照指定的JSON格式返回：
{
    "hasFood": true,  // 布尔值，表示是否包含食物
    "foodType": "食物名称",  // 字符串，食物的类型
    "weight": 100,  // 数字，估计重量（克）
    "nutrition": {
        "calories": 0,  // 数字，热量（卡路里）
        "protein": 0,  // 数字，蛋白质（克）
        "totalFat": 0,  // 数字，总脂肪（克）
        "saturatedFat": 0,  // 数字，饱和脂肪（克）
        "transFat": 0,  // 数字，反式脂肪（克）
        "unsaturatedFat": 0,  // 数字，不饱和脂肪（克）
        "carbohydrates": 0,  // 数字，碳水化合物（克）
        "sugar": 0,  // 数字，糖分（克）
        "fiber": 0,  // 数字，膳食纤维（克）
        "vitamins": {
            "vitaminA": 0,  // 数字，维生素A（毫克）
            "vitaminC": 0,  // 数字，维生素C（毫克）
            "vitaminD": 0,  // 数字，维生素D（微克）
            "vitaminE": 0,  // 数字，维生素E（毫克）
            "vitaminK": 0,  // 数字，维生素K（微克）
            "vitaminB": {
                "b1": 0,  // 数字，维生素B1（毫克）
                "b2": 0,  // 数字，维生素B2（毫克）
                "b6": 0,  // 数字，维生素B6（毫克）
                "b12": 0  // 数字，维生素B12（微克）
            }
        },
        "minerals": {
            "calcium": 0,  // 数字，钙（毫克）
            "iron": 0,  // 数字，铁（毫克）
            "sodium": 0,  // 数字，钠（毫克）
            "potassium": 0,  // 数字，钾（毫克）
            "zinc": 0,  // 数字，锌（毫克）
            "magnesium": 0  // 数字，镁（毫克）
        }
    }
}"""

def test_deepseek_v3():
    """测试DeepSeek V3模型API"""
    global API_KEY, DEFAULT_PROMPT, MODEL, API_URL, SYSTEM_PROMPT
    
    # 请求头
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {API_KEY}"
    }
    
    # 请求数据
    data = {
        "model": MODEL,
        "messages": [
            {"role": "system", "content": SYSTEM_PROMPT},
            {"role": "user", "content": DEFAULT_PROMPT}
        ],
        "max_tokens": 1000  # 限制生成的最大token数量
    }
    
    logging.info(f"发送请求: '{DEFAULT_PROMPT}'")
    logging.info(f"使用模型: {MODEL}")
    logging.info(f"API URL: {API_URL}")
    
    try:
        # 发送请求
        logging.info("正在发送请求到DeepSeek API...")
        response = requests.post(API_URL, headers=headers, json=data)
        
        # 检查响应
        response.raise_for_status()
        
        # 解析结果
        result = response.json()
        
        # 检查API错误
        if "error" in result and result["error"] is not None:
            logging.error(f"API返回错误: {result['error'].get('message', '未知错误')}")
            return None
        
        # 检查响应中是否有choices
        if "choices" not in result or len(result["choices"]) == 0:
            logging.error("API响应中没有choices")
            return None
        
        # 提取生成的文本
        generated_text = result["choices"][0]["message"]["content"]
        
        print("\n===== 返回结果 =====")
        print(generated_text)
        print("=====================\n")
        
        # 输出token使用情况
        if "usage" in result:
            usage = result["usage"]
            logging.info(f"输入tokens: {usage.get('prompt_tokens', 'N/A')}")
            logging.info(f"输出tokens: {usage.get('completion_tokens', 'N/A')}")
            logging.info(f"总tokens: {usage.get('total_tokens', 'N/A')}")
        
        return generated_text
    
    except requests.exceptions.HTTPError as e:
        logging.error(f"HTTP错误: {e}")
        if 'response' in locals() and hasattr(response, 'text'):
            logging.error(f"响应状态码: {response.status_code}, 错误详情: {response.text}")
        return None
    except requests.exceptions.ConnectionError as e:
        logging.error(f"连接错误: {e}")
        return None
    except requests.exceptions.Timeout as e:
        logging.error(f"请求超时: {e}")
        return None
    except requests.exceptions.RequestException as e:
        logging.error(f"请求异常: {e}")
        if 'response' in locals() and hasattr(response, 'text'):
            logging.error(f"错误详情: {response.text}")
        return None
    except json.JSONDecodeError as e:
        logging.error(f"JSON解析错误: {e}")
        if 'response' in locals() and hasattr(response, 'text'):
            logging.error(f"响应内容: {response.text}")
        return None
    except Exception as e:
        logging.error(f"未预期的错误: {e}")
        return None

def analyze_food_image(image_path: str, description: Optional[str] = None) -> Optional[Dict[str, Any]]:
    """
    分析食物图片并返回营养成分分析结果
    
    Args:
        image_path: 图片文件路径
        description: 可选的图片描述
    
    Returns:
        Dict 包含食物分析结果，如果失败则返回None
    """
    global API_KEY, MODEL, API_URL
    
    # 检查文件是否存在
    if not os.path.exists(image_path):
        logging.error(f"图片文件不存在: {image_path}")
        return None
    
    # 检查文件类型
    file_ext = os.path.splitext(image_path)[1].lower()
    if file_ext not in ['.jpg', '.jpeg', '.png']:
        logging.error(f"不支持的图片格式: {file_ext}")
        return None
    
    # 检查文件大小（限制为10MB）
    file_size = os.path.getsize(image_path)
    if file_size > 10 * 1024 * 1024:
        logging.error(f"图片文件过大: {file_size} 字节")
        return None
    
    try:
        # 读取图片文件
        with open(image_path, 'rb') as f:
            image_bytes = f.read()
        
        # 将图片转换为base64
        base64_image = base64.b64encode(image_bytes).decode('utf-8')
        
        # 构建提示词
        if description:
            prompt = f"""分析这张图片中的食物，分析时尤其注意热量和质量的比例关系，
            可以先考虑食物中的主要成分和高热量成分的质量，再按照比例推算热量，用户描述: {description}
            {FOOD_ANALYSIS_PROMPT}"""
        else:
            prompt = FOOD_ANALYSIS_PROMPT
        
        # 请求头
        headers = {
            "Content-Type": "application/json",
            "Authorization": f"Bearer {API_KEY}"
        }
        
        # 构建请求数据
        # 注意：DeepSeek API可能与OpenAI API的格式略有不同，可能需要调整
        data = {
            "model": MODEL,
            "messages": [
                {
                    "role": "user",
                    "content": [
                        {
                            "type": "text",
                            "text": prompt
                        },
                        {
                            "type": "image_url",
                            "image_url": {
                                "url": f"data:image/jpeg;base64,{base64_image}"
                            }
                        }
                    ]
                }
            ],
            "max_tokens": 4000
        }
        
        logging.info(f"正在发送图片分析请求，图片路径: {image_path}")
        logging.info(f"使用模型: {MODEL}")
        
        # 发送请求
        response = requests.post(API_URL, headers=headers, json=data)
        response.raise_for_status()
        
        # 解析结果
        result = response.json()
        
        # 检查API错误
        if "error" in result and result["error"] is not None:
            logging.error(f"API返回错误: {result['error'].get('message', '未知错误')}")
            return None
        
        # 检查响应中是否有choices
        if "choices" not in result or len(result["choices"]) == 0:
            logging.error("API响应中没有choices")
            return None
        
        # 提取生成的文本
        content = result["choices"][0]["message"]["content"]
        logging.info("API返回的原始内容")
        
        # 如果内容被Markdown代码块包裹，提取其中的JSON
        content = content.strip()
        if content.startswith("```json"):
            content = content[7:]
            if "```" in content:
                content = content[:content.find("```")]
        elif content.startswith("```"):
            content = content[3:]
            if "```" in content:
                content = content[:content.find("```")]
        content = content.strip()
        
        # 解析JSON内容
        try:
            analysis = json.loads(content)
            
            # 打印分析结果
            print("\n===== 食物分析结果 =====")
            if analysis.get("hasFood", False):
                print(f"检测到食物: {analysis.get('foodType', '未知')}")
                print(f"估计重量: {analysis.get('weight', 0)} 克")
                
                nutrition = analysis.get("nutrition", {})
                print(f"\n营养成分:")
                print(f"热量: {nutrition.get('calories', 0)} 卡路里")
                print(f"蛋白质: {nutrition.get('protein', 0)} 克")
                print(f"总脂肪: {nutrition.get('totalFat', 0)} 克")
                print(f"碳水化合物: {nutrition.get('carbohydrates', 0)} 克")
            else:
                print("未检测到食物")
            print("=========================\n")
            
            return analysis
            
        except json.JSONDecodeError as e:
            logging.error(f"解析JSON内容失败: {e}")
            logging.error(f"原始内容: {content}")
            return None
            
    except Exception as e:
        logging.error(f"分析食物图片时发生错误: {e}")
        return None

# 命令行参数解析
def parse_args():
    parser = argparse.ArgumentParser(description="DeepSeek API测试工具")
    parser.add_argument("--chat", action="store_true", help="执行普通聊天测试")
    parser.add_argument("--food", type=str, help="进行食物分析的图片路径")
    parser.add_argument("--description", type=str, help="可选的食物图片描述")
    return parser.parse_args()

# 主函数
if __name__ == "__main__":
    args = parse_args()
    
    if args.food:
        # 分析食物图片
        analyze_food_image(args.food, args.description)
    else:
        # 默认执行普通聊天测试
        test_deepseek_v3()