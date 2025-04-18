import requests
import json
from datetime import datetime

# API基础URL
BASE_URL = "https://jesuvukndxpo.sealoshzh.site/api"

# 测试用户信息
TEST_USER = {
    "email": "965377515@qq.com",
    "password": "123456"
}

def login():
    """登录并获取token"""
    response = requests.post(f"{BASE_URL}/login", json=TEST_USER)
    if response.status_code == 200:
        return response.json()["token"]
    else:
        raise Exception("登录失败:", response.json())

def test_check_in_flow():
    """测试完整的打卡流程"""
    # 1. 登录获取token
    token = login()
    headers = {"Authorization": f"Bearer {token}"}
    
    print("\n=== 开始测试打卡功能 ===")

    # 2. 获取今日打卡状态
    print("\n1. 检查今日打卡状态")
    response = requests.get(f"{BASE_URL}/check-in/today", headers=headers)
    print("状态码:", response.status_code)
    print("响应:", json.dumps(response.json(), ensure_ascii=False, indent=2))

    # 3. 尝试打卡
    print("\n2. 尝试打卡")
    check_in_data = {
        "content": f"今天的打卡内容 - {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}"
    }
    response = requests.post(f"{BASE_URL}/check-in", headers=headers, json=check_in_data)
    print("状态码:", response.status_code)
    print("响应:", json.dumps(response.json(), ensure_ascii=False, indent=2))

    # 4. 再次获取今日打卡状态
    print("\n3. 再次检查今日打卡状态")
    response = requests.get(f"{BASE_URL}/check-in/today", headers=headers)
    print("状态码:", response.status_code)
    print("响应:", json.dumps(response.json(), ensure_ascii=False, indent=2))

    # 5. 获取所有打卡记录
    print("\n4. 获取所有打卡记录")
    response = requests.get(f"{BASE_URL}/check-ins", headers=headers)
    print("状态码:", response.status_code)
    print("响应:", json.dumps(response.json(), ensure_ascii=False, indent=2))

    # 6. 测试分页
    print("\n5. 测试分页获取打卡记录")
    params = {
        "page": 1,
        "page_size": 5
    }
    response = requests.get(f"{BASE_URL}/check-ins", headers=headers, params=params)
    print("状态码:", response.status_code)
    print("响应:", json.dumps(response.json(), ensure_ascii=False, indent=2))

def main():
    try:
        test_check_in_flow()
    except Exception as e:
        print("测试过程中出现错误:", str(e))

if __name__ == "__main__":
    main() 