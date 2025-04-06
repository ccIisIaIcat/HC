#!/usr/bin/env python
# -*- coding: utf-8 -*-

import requests
import json
import os
from datetime import datetime, timedelta

# API基础URL
BASE_URL = "https://jesuvukndxpo.sealoshzh.site/api"

# 测试用户信息
TEST_USER_EMAIL = "965377515@qq.com"
TEST_USER_PASSWORD = "123456"

# 存储测试数据
test_data = {}

def print_response(response, operation_name):
    """打印响应信息"""
    print(f"\n{operation_name}响应:")
    print(f"状态码: {response.status_code}")
    try:
        print(f"响应内容: {json.dumps(response.json(), ensure_ascii=False, indent=2)}")
    except:
        print(f"响应内容: {response.text}")

def login_user():
    """用户登录"""
    print("\n===== 测试用户登录 =====")
    url = f'{BASE_URL}/login'
    data = {
        'email': TEST_USER_EMAIL,
        'password': TEST_USER_PASSWORD
    }
    
    response = requests.post(url, json=data)
    print_response(response, "用户登录")
    
    if response.status_code == 200:
        token = response.json().get('token')
        if token:
            test_data['user_token'] = token
            print("✅ 用户登录成功")
            return token
    
    print("❌ 用户登录失败",response.text)
    return None

def get_headers(token):
    """获取带有认证token的请求头"""
    return {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json"
    }

def test_food_record_api():
    """测试食物记录API"""
    print("\n=== 测试食物记录API ===")
    
    # 登录获取token
    token = login_user()
    if not token:
        return
    
    # 测试数据
    test_record = {
        "food_name": "测试食物",
        "food_type": "主食",
        "calories": 300,
        "protein": 10,
        "fat": 5,
        "carbohydrate": 50,
        "meal_type": "午餐",
        "record_time": datetime.now().strftime("%Y-%m-%dT%H:%M:%S+08:00"),
        "notes": "这是一个测试记录"
    }
    
    # 1. 创建食物记录
    print("\n1. 创建食物记录")
    url = f"{BASE_URL}/food-records"
    print(get_headers(token))
    response = requests.post(url, json=test_record, headers=get_headers(token))
    print_response(response, "创建食物记录")
    
    if response.status_code != 200:
        print("❌ 创建食物记录失败")
        return
    
    result = response.json()
    if "record" not in result:
        print("❌ 响应中没有record字段")
        return
    
    record_id = result["record"]["ID"]
    test_data['food_record_id'] = record_id
    print(f"✅ 创建的食物记录ID: {record_id}")
    
    # 2. 获取食物记录列表
    print("\n2. 获取食物记录列表")
    today = datetime.now()
    thirty_days_ago = today - timedelta(days=30)
    
    params = {
        "start_date": thirty_days_ago.strftime("%Y-%m-%d"),
        "end_date": today.strftime("%Y-%m-%d")
    }
    
    response = requests.get(url, params=params, headers=get_headers(token))
    print_response(response, "获取食物记录列表")
    
    if response.status_code != 200:
        print("❌ 获取食物记录列表失败")
        return
    
    result = response.json()
    if "records" not in result:
        print("❌ 响应中没有records字段")
        return
    
    found = False
    for record in result["records"]:
        if record["ID"] == record_id:
            found = True
            break
    
    if not found:
        print("❌ 未找到创建的测试记录")
    else:
        print("✅ 成功找到创建的测试记录")
    
    # 3. 获取单个食物记录
    print("\n3. 获取单个食物记录")
    url = f"{BASE_URL}/food-records/{record_id}"
    response = requests.get(url, headers=get_headers(token))
    print_response(response, "获取单个食物记录")
    
    if response.status_code != 200:
        print("❌ 获取单个食物记录失败")
        return
    
    result = response.json()
    if "record" not in result:
        print("❌ 响应中没有record字段")
        return
    
    record = result["record"]
    print(f"✅ 食物名称: {record['food_name']}")
    print(f"✅ 食物类型: {record['food_type']}")
    print(f"✅ 卡路里: {record['calories']}")
    
    # 4. 更新食物记录
    print("\n4. 更新食物记录")
    updated_data = {
        "food_name": "更新后的测试食物",
        "calories": 350,
        "notes": "这是更新后的测试记录"
    }
    
    response = requests.put(url, json=updated_data, headers=get_headers(token))
    print_response(response, "更新食物记录")
    
    if response.status_code != 200:
        print("❌ 更新食物记录失败")
        return
    
    result = response.json()
    if "record" not in result:
        print("❌ 响应中没有record字段")
        return
    
    record = result["record"]
    print(f"✅ 更新后的食物名称: {record['food_name']}")
    print(f"✅ 更新后的卡路里: {record['calories']}")
    print(f"✅ 更新后的备注: {record['notes']}")
    
    # 5. 删除食物记录
    print("\n5. 删除食物记录")
    response = requests.delete(url, headers=get_headers(token))
    print_response(response, "删除食物记录")
    
    if response.status_code != 200:
        print("❌ 删除食物记录失败")
        return
    
    # 验证记录已被删除
    response = requests.get(url, headers=get_headers(token))
    if response.status_code != 404:
        print("❌ 记录应该已被删除")
    else:
        print("✅ 记录已成功删除")
    
    # 6. 分析并保存食物记录
    print("\n6. 分析并保存食物记录")
    url = f"{BASE_URL}/analyze-and-save"
    
    # 准备测试图片
    image_path = os.path.join(os.path.dirname(__file__), "test_food_image.jpg")
    
    # 检查测试图片是否存在
    if not os.path.exists(image_path):
        print(f"⚠️ 警告: 测试图片 {image_path} 不存在，跳过此测试")
        return
    
    # 准备表单数据
    files = {
        "image": ("test_food_image.jpg", open(image_path, "rb"), "image/jpeg")
    }
    data = {
        "meal_type": "晚餐",
        "notes": "这是一个测试分析记录",
        "image_description": "一盘水果沙拉"
    }
    
    # 发送请求
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.post(url, files=files, data=data, headers=headers)
    
    # 关闭文件
    files["image"][1].close()
    
    print_response(response, "分析并保存食物记录")
    
    if response.status_code != 200:
        print("❌ 分析并保存食物记录失败")
        return
    
    result = response.json()
    if "record" not in result or "analysis" not in result:
        print("❌ 响应中缺少必要字段")
        return
    
    print("✅ 分析并保存食物记录成功")
    print(f"✅ 记录ID: {result['record']['ID']}")
    print(f"✅ 分析结果: {result['analysis']}")

def test_health_state_api():
    """测试健康状态记录API"""
    print("\n=== 测试健康状态记录API ===")
    
    # 登录获取token
    token = login_user()
    if not token:
        return
    
    # 测试数据
    test_record = {
        "height": 175.0,
        "weight": 70.0,
        "bmi": 22.86,
        "body_fat_percentage": 18.5,
        "temperature": 36.5,
        "heart_rate": 75,
        "respiratory_rate": 16,
        "fasting_glucose": 5.2,
        "postprandial_glucose": 7.1,
        "total_cholesterol": 4.8,
        "ldl_cholesterol": 2.9,
        "hdl_cholesterol": 1.2,
        "triglycerides": 1.5,
        "notes": "这是一个测试健康状态记录"
    }
    
    # 1. 创建健康状态记录
    print("\n1. 创建健康状态记录")
    url = f"{BASE_URL}/health-states"
    response = requests.post(url, json=test_record, headers=get_headers(token))
    print_response(response, "创建健康状态记录")
    
    if response.status_code != 200:
        print("❌ 创建健康状态记录失败")
        return
    
    result = response.json()
    if "record" not in result:
        print("❌ 响应中没有record字段")
        return
    
    record_id = result["record"]["ID"]
    test_data['health_state_id'] = record_id
    print(f"✅ 创建的健康状态记录ID: {record_id}")
    
    # 2. 获取健康状态记录列表
    print("\n2. 获取健康状态记录列表")
    today = datetime.now()
    thirty_days_ago = today - timedelta(days=30)
    
    params = {
        "start_date": thirty_days_ago.strftime("%Y-%m-%d"),
        "end_date": today.strftime("%Y-%m-%d")
    }
    
    response = requests.get(url, params=params, headers=get_headers(token))
    print_response(response, "获取健康状态记录列表")
    
    if response.status_code != 200:
        print("❌ 获取健康状态记录列表失败")
        return
    
    result = response.json()
    if "records" not in result:
        print("❌ 响应中没有records字段")
        return
    
    found = False
    for record in result["records"]:
        if record["ID"] == record_id:
            found = True
            break
    
    if not found:
        print("❌ 未找到创建的测试记录")
    else:
        print("✅ 成功找到创建的测试记录")
    
    # 3. 获取单个健康状态记录
    print("\n3. 获取单个健康状态记录")
    url = f"{BASE_URL}/health-states/{record_id}"
    response = requests.get(url, headers=get_headers(token))
    print_response(response, "获取单个健康状态记录")
    
    if response.status_code != 200:
        print("❌ 获取单个健康状态记录失败")
        return
    
    result = response.json()
    if "record" not in result:
        print("❌ 响应中没有record字段")
        return
    
    record = result["record"]
    print(f"✅ 身高: {record['height']}")
    print(f"✅ 体重: {record['weight']}")
    print(f"✅ BMI: {record['bmi']}")
    
    # 4. 更新健康状态记录
    print("\n4. 更新健康状态记录")
    updated_data = {
        "weight": 72.0,
        "bmi": 23.51,
        "body_fat_percentage": 19.0,
        "notes": "这是更新后的测试记录"
    }
    
    response = requests.put(url, json=updated_data, headers=get_headers(token))
    print_response(response, "更新健康状态记录")
    
    if response.status_code != 200:
        print("❌ 更新健康状态记录失败")
        return
    
    result = response.json()
    if "record" not in result:
        print("❌ 响应中没有record字段")
        return
    
    record = result["record"]
    print(f"✅ 更新后的体重: {record['weight']}")
    print(f"✅ 更新后的BMI: {record['bmi']}")
    print(f"✅ 更新后的体脂率: {record['body_fat_percentage']}")
    print(f"✅ 更新后的备注: {record['notes']}")
    
    # 5. 删除健康状态记录
    print("\n5. 删除健康状态记录")
    response = requests.delete(url, headers=get_headers(token))
    print_response(response, "删除健康状态记录")
    
    if response.status_code != 200:
        print("❌ 删除健康状态记录失败")
        return
    
    # 验证记录已被删除
    response = requests.get(url, headers=get_headers(token))
    if response.status_code != 404:
        print("❌ 记录应该已被删除")
    else:
        print("✅ 记录已成功删除")
    
    # 6. 获取最新的健康状态记录
    print("\n6. 获取最新的健康状态记录")
    url = f"{BASE_URL}/health-states/latest"
    response = requests.get(url, headers=get_headers(token))
    print_response(response, "获取最新健康状态记录")
    
    if response.status_code != 200:
        print("❌ 获取最新健康状态记录失败")
        return
    
    result = response.json()
    if "record" not in result:
        print("❌ 响应中没有record字段")
        return
    
    record = result["record"]
    print(f"✅ 最新记录ID: {record['ID']}")

if __name__ == "__main__":
    test_food_record_api()
