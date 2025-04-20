#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
import sys
import requests
from pathlib import Path

# 配置参数
CONFIG = {
    'apk_path': 'FTM_beta_1.0.8.apk',  # APK文件路径
    'version': '1.0.8',     # 版本号
    'platform': 'android',  # 平台
    'force_update': False,  # 是否强制更新
    'update_notes': '新版本更新说明',  # 更新说明
    'server_url': 'https://jesuvukndxpo.sealoshzh.site',  # 修改为本地地址
    'admin_email': 'admin@yourdomain.com',  # 管理员邮箱
    'admin_password': 'admin888'  # 管理员密码
}


def login_admin():
    """
    管理员登录获取token
    """
    url = f"{CONFIG['server_url']}/api/login"
    data = {
        'email': CONFIG['admin_email'],
        'password': CONFIG['admin_password']
    }
    
    try:
        print(f"正在尝试登录: {url}")
        print(f"登录数据: {data}")
        response = requests.post(url, json=data)
        print(f"登录响应状态码: {response.status_code}")
        print(f"登录响应内容: {response.text}")
        
        if response.status_code == 200:
            token = response.json().get('token')
            if token:
                print("✅ 管理员登录成功")
                print(f"获取到的token: {token[:10]}...")  # 只显示token的前10个字符
                return token
            else:
                print("❌ 登录响应中没有token")
                return None
        else:
            print(f"❌ 登录失败！状态码: {response.status_code}")
            print(f"错误信息: {response.text}")
            return None
    except Exception as e:
        print(f"❌ 登录过程中发生错误: {str(e)}")
        return None

def upload_apk(admin_token):
    """
    上传APK文件到服务器
    """
    # 检查文件是否存在
    if not os.path.exists(CONFIG['apk_path']):
        print(f"错误：文件 {CONFIG['apk_path']} 不存在")
        return False

    # 检查文件扩展名
    if not CONFIG['apk_path'].lower().endswith('.apk'):
        print("错误：文件必须是APK格式")
        return False

    # 检查文件大小
    file_size = os.path.getsize(CONFIG['apk_path'])
    print(f"文件大小: {file_size / 1024 / 1024:.2f} MB")

    # 准备上传数据
    url = f"{CONFIG['server_url']}/api/app-updates/upload"
    
    # 确保token格式正确
    if not admin_token.startswith('Bearer '):
        admin_token = f'Bearer {admin_token}'
    
    headers = {
        "Authorization": admin_token
    }
    
    files = {
        'file': ('app.apk', open(CONFIG['apk_path'], 'rb'), 'application/vnd.android.package-archive')
    }
    
    data = {
        'version': CONFIG['version'],
        'platform': CONFIG['platform'],
        'force_update': str(CONFIG['force_update']).lower(),
        'update_notes': CONFIG['update_notes']
    }

    try:
        print(f"\n正在尝试上传APK到: {url}")
        print(f"完整请求头: {headers}")
        print(f"表单数据: {data}")
        print(f"Token长度: {len(admin_token)}")
        
        # 发送请求，设置较长的超时时间
        response = requests.post(
            url, 
            headers=headers, 
            files=files, 
            data=data,
            timeout=300  # 5分钟超时
        )
        print(f"上传响应状态码: {response.status_code}")
        print(f"上传响应头: {dict(response.headers)}")
        print(f"上传响应内容: {response.text}")
        
        # 检查响应
        if response.status_code == 200:
            print("✅ 上传成功！")
            print(f"版本号: {CONFIG['version']}")
            print(f"平台: {CONFIG['platform']}")
            print(f"强制更新: {CONFIG['force_update']}")
            print(f"更新说明: {CONFIG['update_notes']}")
            return True
        else:
            print(f"❌ 上传失败！状态码: {response.status_code}")
            print(f"错误信息: {response.text}")
            return False
            
    except requests.exceptions.Timeout:
        print("❌ 上传超时，请检查网络连接或文件大小")
        return False
    except requests.exceptions.ConnectionError:
        print("❌ 连接错误，请检查网络连接或服务器状态")
        return False
    except Exception as e:
        print(f"❌ 上传过程中发生错误: {str(e)}")
        return False
    finally:
        files['file'][1].close()

if __name__ == '__main__':
    # 先登录获取token
    print("\n===== 开始上传APK =====")
    admin_token = login_admin()
    if admin_token:
        # 使用token上传APK
        upload_apk(admin_token)
    else:
        print("❌ 无法获取管理员token，上传终止") 