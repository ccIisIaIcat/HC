#!/usr/bin/env python
# -*- coding: utf-8 -*-

import requests
import json
import time
import os
from datetime import datetime, timedelta

# 测试配置
BASE_URL = "https://jesuvukndxpo.sealoshzh.site"  # 根据实际情况修改
# 测试账户信息
ADMIN_EMAIL = 'admin@yourdomain.com'
ADMIN_PASSWORD = 'admin888'
TEST_USER_NAME = '测试用户'
TEST_USER_EMAIL = '965377515@qq.com'
TEST_USER_PASSWORD = '123456'

# 全局变量
token = None
health_state_id = None


def login():
    """登录并获取token"""
    global token
    print("\n===== 测试用户登录 =====")
    url = f'{BASE_URL}/api/login'
    data = {
        'email': TEST_USER_EMAIL,
        'password': TEST_USER_PASSWORD
    }
    
    response = requests.post(url, json=data)
    print(f"登录响应状态码: {response.status_code}")
    print(f"登录响应内容: {response.text[:100]}...")
    
    if response.status_code == 200:
        token = response.json().get('token')
        if token:
            print(f"✅ 用户登录成功，获取到token: {token[:10]}...")
            return token
    
    print("❌ 用户登录失败")
    return None


def get_headers():
    """获取带有认证token的请求头"""
    return {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json"
    }


def test_create_health_state():
    """测试创建健康状态记录"""
    global health_state_id
    url = f"{BASE_URL}/api/health-states"
    
    # 准备测试数据
    data = {
        "height": 175.0,
        "weight": 70.0,
        "bmi": 22.86,
        "body_fat_percentage": 18.5,
        "temperature": 36.5,
        "heart_rate": 75,
        "respiratory_rate": 16,
        "fasting_glucose": 5.2,
        "postprandial_glucose": 6.8,
        "total_cholesterol": 4.8,
        "notes": "测试健康状态记录"
    }
    
    response = requests.post(url, headers=get_headers(), json=data)
    print(f"创建健康状态记录响应状态码: {response.status_code}")
    print(f"创建健康状态记录响应内容: {response.text[:100]}...")
    
    assert response.status_code == 200, f"创建健康状态记录失败: {response.text}"
    
    result = response.json()
    assert "record" in result, "响应中没有record字段"
    assert "ID" in result["record"], "记录中没有ID字段"
    
    health_state_id = result["record"]["ID"]
    print(f"✅ 创建健康状态记录成功，ID: {health_state_id}")
    return health_state_id


def test_get_health_states():
    """测试获取健康状态记录列表"""
    url = f"{BASE_URL}/api/health-states"
    
    # 测试默认查询（最近30天）
    response = requests.get(url, headers=get_headers())
    print(f"获取健康状态记录列表响应状态码: {response.status_code}")
    print(f"获取健康状态记录列表响应内容: {response.text[:100]}...")
    
    assert response.status_code == 200, f"获取健康状态记录列表失败: {response.text}"
    
    result = response.json()
    assert "records" in result, "响应中没有records字段"
    assert "total" in result, "响应中没有total字段"
    
    print(f"✅ 获取健康状态记录列表成功，共{result['total']}条记录")
    
    # 测试日期范围查询
    today = datetime.now()
    start_date = (today - timedelta(days=7)).strftime("%Y-%m-%d")
    end_date = today.strftime("%Y-%m-%d")
    
    params = {
        "start_date": start_date,
        "end_date": end_date
    }
    
    response = requests.get(url, headers=get_headers(), params=params)
    print(f"按日期范围获取健康状态记录列表响应状态码: {response.status_code}")
    print(f"按日期范围获取健康状态记录列表响应内容: {response.text[:100]}...")
    
    assert response.status_code == 200, f"按日期范围获取健康状态记录列表失败: {response.text}"
    
    result = response.json()
    print(f"✅ 按日期范围获取健康状态记录列表成功，共{result['total']}条记录")
    
    return result["records"]


def test_get_health_state():
    """测试获取单个健康状态记录"""
    global health_state_id
    url = f"{BASE_URL}/api/health-states/{health_state_id}"
    
    response = requests.get(url, headers=get_headers())
    print(f"获取单个健康状态记录响应状态码: {response.status_code}")
    print(f"获取单个健康状态记录响应内容: {response.text[:100]}...")
    
    assert response.status_code == 200, f"获取单个健康状态记录失败: {response.text}"
    
    result = response.json()
    assert "record" in result, "响应中没有record字段"
    assert result["record"]["ID"] == health_state_id, "返回的记录ID不匹配"
    
    print(f"✅ 获取单个健康状态记录成功，ID: {result['record']['ID']}")
    return result["record"]


def test_get_latest_health_state():
    """测试获取最新的健康状态记录"""
    url = f"{BASE_URL}/api/health-states/latest"
    
    response = requests.get(url, headers=get_headers())
    print(f"获取最新健康状态记录响应状态码: {response.status_code}")
    print(f"获取最新健康状态记录响应内容: {response.text[:100]}...")
    
    assert response.status_code == 200, f"获取最新健康状态记录失败: {response.text}"
    
    result = response.json()
    assert "record" in result, "响应中没有record字段"
    
    print(f"✅ 获取最新健康状态记录成功，ID: {result['record']['ID']}")
    return result["record"]


def test_update_health_state():
    """测试更新健康状态记录"""
    global health_state_id
    url = f"{BASE_URL}/api/health-states/{health_state_id}"
    
    # 准备更新数据
    data = {
        "weight": 71.0,
        "bmi": 23.18,
        "notes": "更新后的测试健康状态记录"
    }
    
    response = requests.put(url, headers=get_headers(), json=data)
    print(f"更新健康状态记录响应状态码: {response.status_code}")
    print(f"更新健康状态记录响应内容: {response.text[:100]}...")
    
    assert response.status_code == 200, f"更新健康状态记录失败: {response.text}"
    
    result = response.json()
    assert "record" in result, "响应中没有record字段"
    assert result["record"]["ID"] == health_state_id, "返回的记录ID不匹配"
    assert result["record"]["weight"] == 71.0, "体重未更新"
    assert result["record"]["bmi"] == 23.18, "BMI未更新"
    assert result["record"]["notes"] == "更新后的测试健康状态记录", "备注未更新"
    
    print(f"✅ 更新健康状态记录成功，ID: {result['record']['ID']}")
    return result["record"]


def test_delete_health_state():
    """测试删除健康状态记录"""
    global health_state_id
    url = f"{BASE_URL}/api/health-states/{health_state_id}"
    
    response = requests.delete(url, headers=get_headers())
    print(f"删除健康状态记录响应状态码: {response.status_code}")
    print(f"删除健康状态记录响应内容: {response.text[:100]}...")
    
    assert response.status_code == 200, f"删除健康状态记录失败: {response.text}"
    
    result = response.json()
    assert "message" in result, "响应中没有message字段"
    
    print(f"✅ 删除健康状态记录成功，ID: {health_state_id}")
    
    # 验证记录已被删除
    response = requests.get(url, headers=get_headers())
    print(f"验证记录删除响应状态码: {response.status_code}")
    
    assert response.status_code == 404, "记录未被删除"


def run_all_tests():
    """运行所有测试"""
    print("开始测试用户健康状态API...")
    
    # 登录获取token
    if not login():
        print("❌ 登录失败，无法继续测试")
        return
    
    # 测试创建健康状态记录
    test_create_health_state()
    
    # 测试获取健康状态记录列表
    test_get_health_states()
    
    # 测试获取单个健康状态记录
    test_get_health_state()
    
    # 测试获取最新的健康状态记录
    test_get_latest_health_state()
    
    # 测试更新健康状态记录
    test_update_health_state()
    
    # 测试删除健康状态记录
    test_delete_health_state()
    
    print("✅ 所有测试完成！")


if __name__ == "__main__":
    run_all_tests() 