#!/usr/bin/env python
# -*- coding: utf-8 -*-

import requests
import json
import datetime
import sys

# API基础URL
BASE_URL = "http://localhost:8080/api"

# 测试登录函数
def login(email, password):
    """
    测试登录功能并返回认证token
    """
    login_url = f"{BASE_URL}/login"
    login_data = {
        "email": email,
        "password": password
    }
    
    try:
        response = requests.post(login_url, json=login_data)
        response.raise_for_status()
        result = response.json()
        print(f"登录成功! 用户: {email}")
        return result.get("token")
    except requests.exceptions.RequestException as e:
        print(f"登录失败: {e}")
        if hasattr(response, 'text'):
            print(f"错误详情: {response.text}")
        return None

# 测试健康分析API
def test_health_analysis(token, start_date, end_date, analysis_type="comprehensive", description=""):
    """
    测试健康分析API
    
    参数:
        token (str): 认证token
        start_date (str): 开始日期，格式为YYYY-MM-DD
        end_date (str): 结束日期，格式为YYYY-MM-DD
        analysis_type (str): 分析类型，可选值: comprehensive, nutrition, calories, macros
        description (str): 用户描述(可选)
    """
    analysis_url = f"{BASE_URL}/health-analysis"
    
    # 准备请求数据
    analysis_data = {
        "start_date": start_date,
        "end_date": end_date,
        "analysis_type": analysis_type,
        "description": description
    }
    
    # 准备请求头(带认证token)
    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json"
    }
    
    print("\n===== 健康分析请求 =====")
    print(f"时间范围: {start_date} 至 {end_date}")
    print(f"分析类型: {analysis_type}")
    if description:
        print(f"用户描述: {description}")
    
    try:
        response = requests.post(analysis_url, json=analysis_data, headers=headers)
        response.raise_for_status()
        result = response.json()
        print("\n===== 分析结果 =====")
        print(result.get("analysis", "无分析结果"))
        return result
    except requests.exceptions.RequestException as e:
        print(f"\n请求失败: {e}")
        if hasattr(response, 'text'):
            print(f"错误详情: {response.text}")
        return None

# 主程序
def main():
    # 默认参数
    email = "965377515@qq.com"
    password = "123456"
    
    # 如果提供了登录信息命令行参数
    if len(sys.argv) > 6 and sys.argv[5] == "--login":
        email = sys.argv[6]
        if len(sys.argv) > 7:
            password = sys.argv[7]
    
    # 获取当前日期
    today = datetime.date.today()
    default_end_date = (today + datetime.timedelta(days=2)).strftime("%Y-%m-%d")
    
    # 默认获取最近7天的数据
    default_start_date = (today - datetime.timedelta(days=7)).strftime("%Y-%m-%d")
    
    # 从命令行参数获取日期范围(如果提供)
    if len(sys.argv) > 2:
        start_date = sys.argv[1]
        end_date = sys.argv[2]
    else:
        start_date = default_start_date
        end_date = default_end_date
    
    # 从命令行参数获取分析类型(如果提供)
    analysis_type = sys.argv[3] if len(sys.argv) > 3 else "comprehensive"
    
    # 从命令行参数获取用户描述(如果提供)
    description = sys.argv[4] if len(sys.argv) > 4 else "我正在尝试减肥，并增加蛋白质摄入。"
    
    # 登录并获取token
    token = login(email, password)
    if not token:
        print("未能获取认证token，测试终止")
        return
    
    # 测试健康分析API
    test_health_analysis(token, start_date, end_date, analysis_type, description)
    
if __name__ == "__main__":
    print("=== 健康分析API测试脚本 ===")
    print("使用方法: python test_health_api.py [起始日期 结束日期 分析类型 用户描述] [--login 邮箱 密码]")
    print("日期格式: YYYY-MM-DD")
    print("分析类型: comprehensive, nutrition, calories, macros")
    print("示例: python test_health_api.py 2023-01-01 2023-01-31 nutrition '我想增加肌肉'")
    print("使用自定义账号: python test_health_api.py 2023-01-01 2023-01-31 nutrition '描述' --login example@email.com password123")
    print("\n默认登录用户: 965377515@qq.com")
    print("如果不提供参数，将使用默认参数。")
    print("=====================\n")
    main() 