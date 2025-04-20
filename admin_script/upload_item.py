import requests
import os
import json
import mimetypes

# 服务器配置
BASE_URL = "https://jesuvukndxpo.sealoshzh.site"
API_BASE_URL = BASE_URL

# 管理员账户信息
ADMIN_USERNAME = "admin@yourdomain.com"
ADMIN_PASSWORD = "admin888"

# 图片配置
ITEM_IMAGE_PATH = "G:/cc/HC2/HC/tangyuan.png"  # 物品大图路径
ITEM_ICON_PATH = "G:/cc/HC2/HC/tangyuan.png"    # 物品图标路径

ITEM_NAME = "汤圆3"
ITEM_DESCRIPTION = "测试使用的汤圆"
ITEM_SOURCE = "糯米"

def get_file_mime_type(file_path):
    """自动检测文件的MIME类型"""
    mime_type, _ = mimetypes.guess_type(file_path)
    if mime_type is None:
        ext = os.path.splitext(file_path)[1].lower()
        mime_map = {
            '.png': 'image/png',
            '.jpg': 'image/jpeg',
            '.jpeg': 'image/jpeg',
            '.gif': 'image/gif'
        }
        mime_type = mime_map.get(ext, 'application/octet-stream')
    return mime_type

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

def upload_image(file_path, directory, token):
    """上传图片并返回URL"""
    print(f"\n正在上传文件: {file_path}")
    print(f"目标目录: {directory}")
    
    if not os.path.exists(file_path):
        raise Exception(f"文件不存在: {file_path}")
    
    with open(file_path, 'rb') as f:
        mime_type = get_file_mime_type(file_path)
        files = {
            'file': (
                os.path.basename(file_path),
                f,
                mime_type
            )
        }
        data = {'directory': directory}
        headers = {'Authorization': f'Bearer {token}'}
        
        response = requests.post(
            f"{API_BASE_URL}/api/admin/upload/image",
            files=files,
            data=data,
            headers=headers
        )
        
        if response.status_code == 200:
            url = response.json()["url"]
            print(f"上传成功！URL: {url}")
            return url
        else:
            raise Exception(f"上传失败: {response.text}")

def create_item(token, item_data):
    """创建物品"""
    print("\n正在创建物品...")
    print(f"物品数据: {json.dumps(item_data, ensure_ascii=False, indent=2)}")
    
    headers = {
        'Authorization': f'Bearer {token}',
        'Content-Type': 'application/json'
    }
    
    response = requests.post(
        f"{API_BASE_URL}/api/admin/items",
        json=item_data,
        headers=headers
    )
    
    if response.status_code == 200:
        print("物品创建成功！")
        return response.json()
    else:
        raise Exception(f"创建物品失败: {response.text}")

def main():
    try:
        # 1. 登录管理员账户
        token = login_admin()
        
        # 2. 上传物品图片和图标
        image_url = upload_image(ITEM_IMAGE_PATH, "item/pic", token)
        icon_url = upload_image(ITEM_ICON_PATH, "item/icon", token)
        
        # 3. 准备物品数据
        item_data = {
            "name": ITEM_NAME,
            "description": ITEM_DESCRIPTION,
            "source": ITEM_SOURCE,
            "icon_url": icon_url,
            "image_url": image_url
        }
        
        # 4. 创建物品
        result = create_item(token, item_data)
        print("\n完整的物品信息:")
        print(json.dumps(result, ensure_ascii=False, indent=2))
        
    except Exception as e:
        print(f"\n错误: {str(e)}")

if __name__ == "__main__":
    main()
