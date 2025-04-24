import requests
import json

# 服务器配置
BASE_URL = "https://jesuvukndxpo.sealoshzh.site"
API_BASE_URL = BASE_URL

# 管理员账户信息
ADMIN_USERNAME = "admin@yourdomain.com"
ADMIN_PASSWORD = "admin888"

def login_admin():
    """登录管理员账户并获取token"""
    print("正在登录管理员账户...")
    response = requests.post(
        f"{API_BASE_URL}/api/login",
        json={
            "email": ADMIN_USERNAME,
            "password": ADMIN_PASSWORD
        }
    )
    if response.status_code == 200:
        token = response.json()["token"]
        print("管理员登录成功！")
        return token
    else:
        raise Exception(f"管理员登录失败: {response.text}")

def update_item(token, item_id, update_data):
    """根据ID更新物品"""
    print("\n正在更新物品...")
    print(f"物品ID: {item_id}")
    print(f"更新数据: {json.dumps(update_data, ensure_ascii=False, indent=2)}")
    
    headers = {
        'Authorization': f'Bearer {token}',
        'Content-Type': 'application/json'
    }
    
    response = requests.put(
        f"{API_BASE_URL}/api/admin/items/{item_id}",
        json=update_data,
        headers=headers
    )
    
    if response.status_code == 200:
        print("物品更新成功！")
        return response.json()
    else:
        raise Exception(f"更新物品失败: {response.text}")

def main():
    try:
        # 1. 登录管理员账户
        token = login_admin()
        
        # 2. 指定要更新的物品ID和更新数据
        item_id = 58  # 替换为实际的物品ID
        update_data = {
            "name": "不倒翁餐盘",
            "description": "呵～居然会中断打卡？弱爆了好吗！不过看你又重新振作的样子，有点像...不倒翁？【在打卡中断后重新打卡3天】",
            "source": "饭团猫",
        }
        
        # 3. 更新物品
        result = update_item(token, item_id, update_data)
        print("\n更新后的物品信息:")
        print(json.dumps(result, ensure_ascii=False, indent=2))
        
    except Exception as e:
        print(f"\n错误: {str(e)}")

if __name__ == "__main__":
    main()