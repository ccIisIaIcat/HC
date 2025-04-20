import requests
import json
from datetime import datetime

# 全局配置
BASE_URL = "http://localhost:8080/api"

# 测试账号信息
ADMIN_USER = {
    "email": "admin@yourdomain.com",
    "password": "admin888"
}

NORMAL_USER = {
    "email": "965377515@qq.com",
    "password": "123456"
}

# 存储登录后的token和用户ID
admin_token = None
user_token = None
admin_user_id = None
user_user_id = None

def login(email, password):
    """登录并获取token和用户信息"""
    response = requests.post(
        f"{BASE_URL}/login",
        json={"email": email, "password": password}
    )
    if response.status_code == 200:
        data = response.json()
        return data["token"], data["user"]["ID"] # 假设登录响应中包含user_id
    else:
        raise Exception(f"登录失败: {response.text}")

def test_setup():
    """初始化测试环境，登录获取token"""
    global admin_token, user_token, admin_user_id, user_user_id
    print("\n=== 开始测试设置 ===")
    
    admin_token, admin_user_id = login(ADMIN_USER["email"], ADMIN_USER["password"])
    user_token, user_user_id = login(NORMAL_USER["email"], NORMAL_USER["password"])
    
    print(f"管理员(ID: {admin_user_id})和普通用户(ID: {user_user_id})登录成功")

def get_headers(token):
    """获取带认证的请求头"""
    return {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json"
    }

# 管理员API测试
def test_admin_item_operations():
    """测试管理员物品管理API"""
    print("\n=== 测试管理员物品管理API ===")
    
    # 创建物品
    item_data = {
        "name": "测试物品",
        "description": "这是一个测试物品",
        "source": "测试来源",
        "icon_url": "http://example.com/icon.png",
        "image_url": "http://example.com/image.png"
    }
    
    response = requests.post(
        f"{BASE_URL}/admin/items",
        headers=get_headers(admin_token),
        json=item_data
    )
    print(f"创建物品结果: {response.status_code}")
    if response.status_code == 200:
        item_id = response.json()["data"]["ID"]
        print(f"创建的物品ID: {item_id}")
        
        # 更新物品
        item_data["description"] = "更新后的描述"
        response = requests.put(
            f"{BASE_URL}/admin/items/{item_id}",
            headers=get_headers(admin_token),
            json=item_data
        )
        print(f"更新物品结果: {response.status_code}")
        
        # 删除物品
        response = requests.delete(
            f"{BASE_URL}/admin/items/{item_id}",
            headers=get_headers(admin_token)
        )
        print(f"删除物品结果: {response.status_code}")

# 普通用户API测试
def test_user_item_queries():
    """测试普通用户物品查询API"""
    print("\n=== 测试普通用户物品查询API ===")
    
    # 获取物品列表
    response = requests.get(
        f"{BASE_URL}/items",
        headers=get_headers(user_token)
    )
    print(f"获取物品列表结果: {response.status_code}")
    
    # 搜索物品
    response = requests.get(
        f"{BASE_URL}/items/search?name=测试",
        headers=get_headers(user_token)
    )
    print(f"搜索物品结果: {response.status_code}")
    
    # 获取物品来源列表
    response = requests.get(
        f"{BASE_URL}/items/sources",
        headers=get_headers(user_token)
    )
    print(f"获取物品来源列表结果: {response.status_code}")

def test_user_item_management():
    """测试用户物品管理API"""
    print("\n=== 测试用户物品管理API ===")
    
    # 获取用户物品列表
    response = requests.get(
        f"{BASE_URL}/user-items/{user_user_id}",
        headers=get_headers(user_token)
    )
    print(f"获取用户物品列表结果: {response.status_code}")
    
    # 创建用户物品关联
    user_item_data = {
        "item_id": 1,  # 假设存在ID为1的物品
        "quantity": 1,
        "obtained_from": "每日签到奖励"  # 指定物品获得来源
    }
    
    response = requests.post(
        f"{BASE_URL}/user-items/{user_user_id}",
        headers=get_headers(user_token),
        json=user_item_data
    )
    print(f"创建用户物品关联结果: {response.status_code}")
    
    if response.status_code == 200:
        # 更新物品数量
        response = requests.put(
            f"{BASE_URL}/user-items/{user_user_id}/{user_item_data['item_id']}/quantity",
            headers=get_headers(user_token),
            json={"quantity": 2}
        )
        print(f"更新物品数量结果: {response.status_code}")
        
        # 获取物品详情
        response = requests.get(
            f"{BASE_URL}/user-items/{user_user_id}/{user_item_data['item_id']}",
            headers=get_headers(user_token)
        )
        print(f"获取物品详情结果: {response.status_code}")
        
        # 检查物品拥有状态
        response = requests.get(
            f"{BASE_URL}/user-items/{user_user_id}/check/{user_item_data['item_id']}",
            headers=get_headers(user_token)
        )
        print(f"检查物品拥有状态结果: {response.status_code}")
        
        # 获取物品总数
        response = requests.get(
            f"{BASE_URL}/user-items/{user_user_id}/count",
            headers=get_headers(user_token)
        )
        print(f"获取物品总数结果: {response.status_code}")
        
        # 按来源获取物品
        response = requests.get(
            f"{BASE_URL}/user-items/{user_user_id}/source?source=测试来源",
            headers=get_headers(user_token)
        )
        print(f"按来源获取物品结果: {response.status_code}")
        
        # 删除用户物品
        response = requests.delete(
            f"{BASE_URL}/user-items/{user_user_id}/{user_item_data['item_id']}",
            headers=get_headers(user_token)
        )
        print(f"删除用户物品结果: {response.status_code}")

def test_user_permission():
    """测试用户权限（尝试访问其他用户的资源）"""
    print("\n=== 测试用户权限 ===")
    
    # 尝试获取其他用户的物品列表
    response = requests.get(
        f"{BASE_URL}/user-items/{admin_user_id}",  # 尝试访问管理员的物品
        headers=get_headers(user_token)
    )
    print(f"尝试访问其他用户物品列表结果: {response.status_code}")

def run_all_tests():
    """运行所有测试"""
    try:
        test_setup()
        test_admin_item_operations()
        test_user_item_queries()
        test_user_item_management()
        test_user_permission()
        print("\n=== 所有测试完成 ===")
    except Exception as e:
        print(f"\n测试过程中出现错误: {str(e)}")

if __name__ == "__main__":
    run_all_tests() 