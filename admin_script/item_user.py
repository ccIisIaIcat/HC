import requests
import json
from datetime import datetime

# 全局配置
BASE_URL = "https://jesuvukndxpo.sealoshzh.site/api"

# 用户配置
USER_CONFIG = {
    "email": "zwj19980717@gmail.com",
    "password": "123456"
}

# 存储登录后的token和用户ID
user_token = None
user_id = None

def login(email, password):
    """登录并获取token和用户信息"""
    response = requests.post(
        f"{BASE_URL}/login",
        json={"email": email, "password": password}
    )
    if response.status_code == 200:
        data = response.json()
        print("登录响应:", json.dumps(data, ensure_ascii=False, indent=2))
        return data["token"], data["user"]["ID"]
    else:
        raise Exception(f"登录失败: {response.text}")

def get_headers(token):
    """获取带认证的请求头"""
    return {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json"
    }

def get_all_items(token):
    """获取所有物品信息"""
    print("\n获取所有物品信息...")
    response = requests.get(
        f"{BASE_URL}/items",
        headers=get_headers(token)
    )
    
    if response.status_code == 200:
        data = response.json()
        items = data.get("data", [])
        print(f"共找到 {len(items)} 个物品:")
        for item in items:
            print(f"ID: {item.get('ID')}, " + 
                  f"名称: {item.get('name', item.get('Name'))}, " + 
                  f"来源: {item.get('source', item.get('Source'))}")
        return items
    else:
        raise Exception(f"获取物品列表失败: {response.text}")

def get_item_info(item_id, token):
    """获取物品详细信息"""
    print(f"\n获取物品信息 (ID: {item_id})...")
    response = requests.get(
        f"{BASE_URL}/items/{item_id}",
        headers=get_headers(token)
    )
    
    if response.status_code == 200:
        data = response.json()
        item = data.get("data", {})
        print(f"物品名称: {item.get('name', item.get('Name'))}")
        print(f"物品描述: {item.get('description', item.get('Description'))}")
        print(f"物品来源: {item.get('source', item.get('Source'))}")
        print(f"图标URL: {item.get('icon_url', item.get('IconURL'))}")
        print(f"图片URL: {item.get('image_url', item.get('ImageURL'))}")
        return item
    else:
        raise Exception(f"获取物品信息失败: {response.text}")

def check_user_item(user_id, item_id, token):
    """检查用户是否拥有该物品"""
    print(f"\n检查用户物品状态...")
    print(f"检查用户ID: {user_id}, 物品ID: {item_id}")
    response = requests.get(
        f"{BASE_URL}/user-items/{user_id}/check/{item_id}",
        headers=get_headers(token)
    )
    
    if response.status_code == 200:
        data = response.json()
        if data.get("has_item"):
            print(f"用户已拥有该物品")
        else:
            print("用户还未拥有该物品")
        return data
    else:
        raise Exception(f"检查物品状态失败: {response.text}")

def get_user_items(user_id, token):
    """获取用户的所有物品"""
    print(f"\n获取用户 {user_id} 的所有物品...")
    response = requests.get(
        f"{BASE_URL}/user-items/{user_id}",
        headers=get_headers(token)
    )
    
    if response.status_code == 200:
        data = response.json()
        items = data.get("data", [])
        print(f"用户共有 {len(items)} 个物品:")
        for item in items:
            print(f"物品ID: {item.get('id')}, " + 
                  f"名称: {item.get('name')}, " + 
                  f"数量: {item.get('total_quantity')}, " + 
                  f"获取来源: {item.get('obtained_from')}")
            
            # 显示物品详细信息
            print(f"    物品描述: {item.get('description')}")
            print(f"    物品来源: {item.get('source')}")
            print(f"    图标URL: {item.get('icon_url')}")
            print(f"    图片URL: {item.get('image_url')}")
        return items
    else:
        raise Exception(f"获取用户物品列表失败: {response.text}")

def add_item_to_user(user_id, item_id, quantity, obtained_from, token):
    """给用户添加物品"""
    print(f"\n正在给用户 {user_id} 添加物品...")
    print(f"添加参数: item_id={item_id}, quantity={quantity}, obtained_from={obtained_from}")
    
    user_item_data = {
        "item_id": item_id,
        "quantity": quantity,
        "obtained_from": obtained_from
    }
    
    print("请求数据:", json.dumps(user_item_data, ensure_ascii=False, indent=2))
    
    response = requests.post(
        f"{BASE_URL}/user-items/{user_id}",
        headers=get_headers(token),
        json=user_item_data
    )
    
    print(f"响应状态码: {response.status_code}")
    print(f"响应内容: {response.text}")
    
    if response.status_code == 200:
        data = response.json()
        result = data.get("data", {})
        print("物品添加成功！")
        return result
    else:
        raise Exception(f"添加物品失败: {response.text}")

def main():
    try:
        global user_token, user_id
        
        # 1. 登录获取token和user_id
        print("正在登录...")
        user_token, user_id = login(USER_CONFIG["email"], USER_CONFIG["password"])
        print(f"登录成功！管理员ID: {user_id}")
        
        # 2. 获取所有物品信息
        all_items = get_all_items(user_token)
        
        if not all_items:
            print("没有找到任何物品")
            return
            
        # 3. 让用户选择要添加的物品
        print("\n请选择要添加的物品ID（直接回车退出）:")
        selected_id = input("物品ID: ").strip()
        
        if not selected_id:
            print("未选择物品，程序退出")
            return
            
        if not selected_id.isdigit():
            print("无效的物品ID，程序退出")
            return
            
        item_id = int(selected_id)
        
        # 4. 获取选中物品的详细信息
        item = get_item_info(item_id, user_token)
        
        # 5. 选择目标用户
        print("\n请输入要发送物品的目标用户ID（直接回车使用当前登录用户）:")
        target_user_id = input("目标用户ID: ").strip()
        
        if not target_user_id:
            target_user_id = user_id
        elif not target_user_id.isdigit():
            print("无效的用户ID，程序退出")
            return
        else:
            target_user_id = int(target_user_id)
            
        print(f"\n将发送物品给用户ID: {target_user_id}")
        
        # 6. 检查用户是否已有该物品
        check_result = check_user_item(target_user_id, item_id, user_token)
        
        # 7. 输入物品数量
        print("\n请输入要发送的物品数量（直接回车默认为1）:")
        quantity_input = input("数量: ").strip()
        quantity = int(quantity_input) if quantity_input.isdigit() else 1
        
        # 8. 输入获取来源
        print("\n请输入物品获取来源（直接回车默认为'管理员赠送'）:")
        obtained_from = input("获取来源: ").strip() or "管理员赠送"
        
        # 9. 添加物品
        result = add_item_to_user(target_user_id, item_id, quantity, obtained_from, user_token)
        print("\n添加结果:")
        print(json.dumps(result, ensure_ascii=False, indent=2))
        
        # 10. 显示用户当前所有物品
        get_user_items(target_user_id, user_token)
        
    except Exception as e:
        print(f"\n错误: {str(e)}")

if __name__ == "__main__":
    main()
