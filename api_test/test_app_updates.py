#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import requests
import json

# 配置参数
CONFIG = {
    'server_url': 'https://jesuvukndxpo.sealoshzh.site',  # 服务器地址
    'platform': 'android',  # 平台
    'version': '1.0.0'     # 当前版本号
}

def test_get_app_updates():
    """
    测试获取所有应用更新
    """
    url = f"{CONFIG['server_url']}/api/app-updates"
    print(f"\n正在测试获取所有应用更新: {url}")
    
    try:
        response = requests.get(url)
        print(f"响应状态码: {response.status_code}")
        print(f"响应内容: {json.dumps(response.json(), ensure_ascii=False, indent=2)}")
        return response.status_code == 200
    except Exception as e:
        print(f"❌ 请求失败: {str(e)}")
        return False

def test_check_update():
    """
    测试检查更新
    """
    url = f"{CONFIG['server_url']}/api/app-updates/check"
    params = {
        'platform': CONFIG['platform'],
        'current_version': CONFIG['version']
    }
    print(f"\n正在测试检查更新: {url}")
    print(f"参数: {params}")
    
    try:
        response = requests.get(url, params=params)
        print(f"响应状态码: {response.status_code}")
        print(f"响应内容: {json.dumps(response.json(), ensure_ascii=False, indent=2)}")
        return response.status_code == 200
    except Exception as e:
        print(f"❌ 请求失败: {str(e)}")
        return False

def test_download_update():
    """
    测试下载更新
    """
    url = f"{CONFIG['server_url']}/api/app-updates/download/{CONFIG['platform']}/{CONFIG['version']}"
    print(f"\n正在测试下载更新: {url}")
    
    try:
        response = requests.get(url, stream=True)
        print(f"响应状态码: {response.status_code}")
        
        if response.status_code == 200:
            # 获取文件名
            filename = response.headers.get('Content-Disposition', '').split('filename=')[-1].strip('"')
            if not filename:
                filename = f"app-{CONFIG['platform']}-{CONFIG['version']}.apk"
            
            # 保存文件
            with open(filename, 'wb') as f:
                for chunk in response.iter_content(chunk_size=8192):
                    if chunk:
                        f.write(chunk)
            print(f"✅ 文件已保存为: {filename}")
            return True
        else:
            print(f"❌ 下载失败: {response.text}")
            return False
    except Exception as e:
        print(f"❌ 请求失败: {str(e)}")
        return False

if __name__ == '__main__':
    print("===== 开始测试应用更新公开路由 =====")
    
    # 测试获取所有应用更新
    print("\n1. 测试获取所有应用更新")
    test_get_app_updates()
    
    # 测试检查更新
    print("\n2. 测试检查更新")
    test_check_update()
    
    # 测试下载更新
    print("\n3. 测试下载更新")
    test_download_update()
    
    print("\n===== 测试完成 =====") 