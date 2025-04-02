#!/usr/bin/env python
# -*- coding: utf-8 -*-

import requests
import os
import json
import time

# 全局配置参数
BASE_URL = 'https://jesuvukndxpo.sealoshzh.site/api'

# 测试账户信息
ADMIN_EMAIL = 'admin@yourdomain.com'
ADMIN_PASSWORD = 'admin888'
TEST_USER_NAME = '测试用户'
TEST_USER_EMAIL = 'zwj19980717@gmail.com'
TEST_USER_PASSWORD = '123456'

# 测试图片路径
TEST_IMAGE_PATH = 'tangyuan.png'  # 确保此图片文件存在于脚本同目录下

# 保存测试状态
test_data = {
    'admin_token': None,
    'user_token': None,
    'verification_code': None
}

def print_response(response, description):
    """打印响应信息"""
    print(f"\n{description} - 状态码: {response.status_code}")
    try:
        print(f"响应内容: {json.dumps(response.json(), ensure_ascii=False, indent=2)}")
    except:
        print(f"响应内容: {response.text}")

def login_admin():
    """管理员登录"""
    print("\n===== 测试管理员登录 =====")
    url = f'{BASE_URL}/login'
    data = {
        'email': ADMIN_EMAIL,
        'password': ADMIN_PASSWORD
    }
    
    response = requests.post(url, json=data)
    print_response(response, "管理员登录")
    
    if response.status_code == 200:
        token = response.json().get('token')
        if token:
            test_data['admin_token'] = token
            print("✅ 管理员登录成功")
            return token
    
    print("❌ 管理员登录失败")
    return None

def get_current_user(token, is_admin=False):
    """获取当前用户信息"""
    user_type = "管理员" if is_admin else "用户"
    print(f"\n===== 测试获取{user_type}信息 =====")
    
    url = f'{BASE_URL}/me'
    headers = {'Authorization': f'Bearer {token}'}
    
    response = requests.get(url, headers=headers)
    print_response(response, f"获取{user_type}信息")
    
    if response.status_code == 200:
        user_info = response.json().get('user')
        if user_info:
            print(f"✅ 获取{user_type}信息成功: {user_info.get('name')} ({user_info.get('email')})")
            return user_info
    
    print(f"❌ 获取{user_type}信息失败")
    return None

def register_user():
    """注册新用户"""
    print("\n===== 测试用户注册 =====")
    url = f'{BASE_URL}/register'
    data = {
        'name': TEST_USER_NAME,
        'email': TEST_USER_EMAIL,
        'password': TEST_USER_PASSWORD
    }
    
    response = requests.post(url, json=data)
    print_response(response, "用户注册")
    
    if response.status_code == 200:
        print("✅ 用户注册成功，等待邮箱验证")
        
        # 直接获取验证码（在实际环境中应该从邮件中获取）
        # 这里我们需要手动输入验证码
        verification_code = input("请输入收到的验证码: ")
        test_data['verification_code'] = verification_code
        return True
    
    print("❌ 用户注册失败")
    return False

def verify_email():
    """验证用户邮箱"""
    print("\n===== 测试邮箱验证 =====")
    
    if not test_data['verification_code']:
        print("❌ 没有验证码，无法验证邮箱")
        return False
    
    url = f'{BASE_URL}/verify-email'
    data = {
        'email': TEST_USER_EMAIL,
        'code': test_data['verification_code']
    }
    
    response = requests.post(url, json=data)
    print_response(response, "邮箱验证")
    
    if response.status_code == 200:
        print("✅ 邮箱验证成功")
        return True
    
    print("❌ 邮箱验证失败")
    return False

def resend_verification():
    """重新发送验证邮件"""
    print("\n===== 测试重发验证邮件 =====")
    url = f'{BASE_URL}/resend-verification'
    data = {
        'email': TEST_USER_EMAIL
    }
    
    response = requests.post(url, json=data)
    print_response(response, "重发验证邮件")
    
    if response.status_code == 200:
        print("✅ 验证邮件重发成功")
        # 更新验证码
        verification_code = input("请输入新收到的验证码: ")
        test_data['verification_code'] = verification_code
        return True
    
    print("❌ 验证邮件重发失败")
    return False

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
    
    print("❌ 用户登录失败")
    return None

def get_all_users(admin_token):
    """获取所有用户列表（管理员功能）"""
    print("\n===== 测试获取所有用户列表 =====")
    url = f'{BASE_URL}/admin/users'
    headers = {'Authorization': f'Bearer {admin_token}'}
    
    response = requests.get(url, headers=headers)
    print_response(response, "获取用户列表")
    
    if response.status_code == 200:
        users = response.json().get('users', [])
        print(f"✅ 获取用户列表成功，共 {len(users)} 个用户")
        return users
    
    print("❌ 获取用户列表失败")
    return None

def get_system_stats(admin_token):
    """获取系统统计信息（管理员功能）"""
    print("\n===== 测试获取系统统计信息 =====")
    url = f'{BASE_URL}/admin/stats'
    headers = {'Authorization': f'Bearer {admin_token}'}
    
    response = requests.get(url, headers=headers)
    print_response(response, "获取系统统计")
    
    if response.status_code == 200:
        data = response.json()
        print("✅ 获取系统统计成功")
        
        # 打印返回的所有统计信息，不管字段名是什么
        for key, value in data.items():
            print(f"- {key}: {value}")
            
        # 针对特定字段进行处理（基于实际返回的字段名）
        if 'total_users' in data:
            print(f"- 总用户数: {data.get('total_users', 'N/A')}")
        if 'verified_users' in data:
            print(f"- 已验证用户数: {data.get('verified_users', 'N/A')}")
        if 'unverified_users' in data:
            print(f"- 未验证用户数: {data.get('unverified_users', 'N/A')}")
            
        return data
    
    print("❌ 获取系统统计失败")
    return None

def analyze_food(token, description=None):
    """测试食物分析接口"""
    print("\n===== 测试食物分析 =====")
    
    if not os.path.exists(TEST_IMAGE_PATH):
        print(f"❌ 错误: 图片文件不存在 {TEST_IMAGE_PATH}")
        return None
    
    url = f'{BASE_URL}/analyze-food'
    headers = {'Authorization': f'Bearer {token}'}
    
    with open(TEST_IMAGE_PATH, 'rb') as img:
        files = {'image': img}
        data = {}
        if description:
            data['image_description'] = description
            print(f"使用图片描述: {description}")
        
        response = requests.post(url, headers=headers, files=files, data=data)
    
    print_response(response, "食物分析")
    
    if response.status_code == 200:
        result = response.json()
        print("✅ 食物分析成功")
        if result.get('hasFood'):
            print(f"- 食物类型: {result.get('foodType', 'N/A')}")
            print(f"- 估计重量: {result.get('weight', 'N/A')}克")
            
            nutrition = result.get('nutrition', {})
            print(f"- 热量: {nutrition.get('calories', 'N/A')}卡路里")
            print(f"- 蛋白质: {nutrition.get('protein', 'N/A')}克")
            print(f"- 碳水化合物: {nutrition.get('carbohydrates', 'N/A')}克")
            print(f"- 脂肪: {nutrition.get('totalFat', 'N/A')}克")
        else:
            print("- 未检测到食物")
        return result
    
    print("❌ 食物分析失败")
    return None

def test_unauthorized_access():
    """测试未授权访问的情况"""
    print("\n===== 测试未授权访问 =====")
    
    # 测试无token访问需要认证的接口
    url = f'{BASE_URL}/me'
    response = requests.get(url)
    print_response(response, "无token访问/me")
    
    if response.status_code == 401:
        print("✅ 未授权测试成功: /me接口正确拒绝了无token请求")
    else:
        print("❌ 未授权测试失败: /me接口未能正确拒绝无token请求")
    
    # 测试普通用户访问管理员接口
    if test_data['user_token']:
        url = f'{BASE_URL}/admin/users'
        headers = {'Authorization': f'Bearer {test_data["user_token"]}'}
        response = requests.get(url, headers=headers)
        print_response(response, "普通用户访问管理员接口")
        
        if response.status_code == 403:
            print("✅ 权限测试成功: 管理员接口正确拒绝了普通用户请求")
        else:
            print("❌ 权限测试失败: 管理员接口未能正确拒绝普通用户请求")

def run_all_tests():
    """运行所有测试"""
    print("\n========== 健康检查系统API测试开始 ==========\n")
    
    # 记录测试开始时间
    start_time = time.time()
    
    # 1. 管理员功能测试
    admin_token = login_admin()
    if admin_token:
        admin_info = get_current_user(admin_token, is_admin=True)
        users = get_all_users(admin_token)
        stats = get_system_stats(admin_token)
    
    # 2. 用户注册和验证流程
    registered = register_user()
    if registered:
        verified = verify_email()
        
        # 如果验证失败，尝试重发验证邮件
        if not verified:
            resend_verification()
            verify_email()
    
    # 3. 用户登录和功能测试
    user_token = login_user()
    if user_token:
        user_info = get_current_user(user_token)
        analyze_food(user_token,"这个食物的重量为157克")
    
    # 4. 未授权访问测试
    test_unauthorized_access()
    
    # 计算测试用时
    end_time = time.time()
    duration = end_time - start_time
    
    print("\n========== 健康检查系统API测试完成 ==========")
    print(f"测试用时: {duration:.2f} 秒")

if __name__ == "__main__":
    run_all_tests() 