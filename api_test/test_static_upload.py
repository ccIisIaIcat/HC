import requests
import os
import json
import mimetypes  # 添加mimetypes模块用于文件类型检测

# 服务器配置
BASE_URL = "https://jesuvukndxpo.sealoshzh.site"
API_BASE_URL = f"{BASE_URL}/api"

# 测试账户信息
ADMIN_USERNAME = "admin@yourdomain.com"
ADMIN_PASSWORD = "admin888"
USER_USERNAME = "zwj19980717@gmail.com"
USER_PASSWORD = "123456"

# 测试图片配置
TEST_IMAGE_PATH = "G:/cc/HC2/HC/tangyuan.png"  # 测试用图片路径
TEST_IMAGE_DIRECTORY = "pic"  # 上传目录名称
TEST_IMAGE_FILENAME = "demo.png"  # 自定义上传后的文件名

# 全局变量
admin_token = None
user_token = None

def get_file_mime_type(file_path):
    """自动检测文件的MIME类型"""
    mime_type, _ = mimetypes.guess_type(file_path)
    if mime_type is None:
        # 如果无法检测，则根据扩展名判断
        ext = os.path.splitext(file_path)[1].lower()
        mime_map = {
            '.png': 'image/png',
            '.jpg': 'image/jpeg',
            '.jpeg': 'image/jpeg',
            '.gif': 'image/gif',
            '.bmp': 'image/bmp',
            '.webp': 'image/webp'
        }
        mime_type = mime_map.get(ext, 'application/octet-stream')
    return mime_type

def login(username, password):
    """登录并获取token"""
    response = requests.post(
        f"{API_BASE_URL}/login",
        json={
            "email": username,
            "password": password
        }
    )
    if response.status_code == 200:
        return response.json()["token"]
    else:
        raise Exception(f"登录失败: {response.text}")

def setup():
    """测试设置：登录获取token"""
    global admin_token, user_token
    
    # 确保测试图片目录存在
    os.makedirs(os.path.dirname(TEST_IMAGE_PATH), exist_ok=True)
    
    # 如果测试图片不存在，创建一个测试图片
    if not os.path.exists(TEST_IMAGE_PATH):
        create_test_image()
    
    # 获取管理员和普通用户的token
    admin_token = login(ADMIN_USERNAME, ADMIN_PASSWORD)
    user_token = login(USER_USERNAME, USER_PASSWORD)
    
    print("测试设置完成")
    print(f"管理员token: {admin_token[:20]}...")
    print(f"用户token: {user_token[:20]}...")

def create_test_image():
    """创建测试用图片"""
    from PIL import Image
    
    # 创建一个100x100的红色图片
    img = Image.new('RGB', (100, 100), color='red')
    img.save(TEST_IMAGE_PATH)
    print(f"创建测试图片: {TEST_IMAGE_PATH}")

def test_upload_image_as_admin():
    """测试管理员上传图片"""
    print("\n测试管理员上传图片...")
    
    with open(TEST_IMAGE_PATH, 'rb') as f:
        # 自动检测文件类型
        mime_type = get_file_mime_type(TEST_IMAGE_PATH)
        
        # 构建文件上传参数
        files = {
            'file': (                    # 表单字段名称，对应后端的 c.FormFile("file")
                os.path.basename(TEST_IMAGE_PATH),  # 原始文件名（用于content-type判断）
                f,                       # 文件对象
                mime_type               # 自动检测的文件类型
            )
        }
        
        # 其他参数
        data = {
            'directory': TEST_IMAGE_DIRECTORY,  # 指定上传目录
            'filename': TEST_IMAGE_FILENAME     # 指定保存后的文件名
        }
        headers = {'Authorization': f'Bearer {admin_token}'}  # 认证信息
        
        print("\n请求详情:")
        print(f"URL: {API_BASE_URL}/admin/upload/image")
        print(f"文件名: {os.path.basename(TEST_IMAGE_PATH)}")
        print(f"MIME类型: {mime_type}")
        print(f"目标目录: {TEST_IMAGE_DIRECTORY}")
        print(f"期望的文件名: {TEST_IMAGE_FILENAME}")
        print(f"Headers: {headers}")
        print(f"Form Data: {data}")
        
        # 发送请求
        response = requests.post(
            f"{API_BASE_URL}/admin/upload/image",
            files=files,
            data=data,
            headers=headers
        )
        
        print(f"\n响应详情:")
        print(f"状态码: {response.status_code}")
        print(f"响应内容: {response.text}")
        
        assert response.status_code == 200, "管理员上传图片应该成功"
        return response.json()["url"]

def test_upload_image_as_user():
    """测试普通用户上传图片（应该失败）"""
    print("\n测试普通用户上传图片...")
    
    with open(TEST_IMAGE_PATH, 'rb') as f:
        # 自动检测文件类型
        mime_type = get_file_mime_type(TEST_IMAGE_PATH)
        
        # 构建文件上传参数
        files = {
            'file': (                    # 表单字段名称，对应后端的 c.FormFile("file")
                os.path.basename(TEST_IMAGE_PATH),  # 原始文件名（用于content-type判断）
                f,                       # 文件对象
                mime_type               # 自动检测的文件类型
            )
        }
        
        # 其他参数
        data = {
            'directory': TEST_IMAGE_DIRECTORY,  # 指定上传目录
            'filename': TEST_IMAGE_FILENAME     # 指定保存后的文件名
        }
        headers = {'Authorization': f'Bearer {user_token}'}  # 认证信息
        
        # 发送请求
        response = requests.post(
            f"{API_BASE_URL}/admin/upload/image",
            files=files,
            data=data,
            headers=headers
        )
        
        print(f"状态码: {response.status_code}")
        print(f"响应: {response.text}")
        
        assert response.status_code in [401, 403], "普通用户上传图片应该被拒绝"

def test_get_image_info():
    """测试获取图片信息"""
    print("\n测试获取图片信息...")
    
    # 先上传一张图片
    upload_url = test_upload_image_as_admin()
    print(f"上传的图片URL: {upload_url}")
    
    # 从URL中提取文件路径
    # 例如从 "/static/pic/xxx.png" 提取 "pic/xxx.png"
    image_path = upload_url.replace("/static/", "", 1)
    print(f"提取的图片路径: {image_path}")
    
    # 测试普通用户获取图片信息
    headers = {'Authorization': f'Bearer {user_token}'}
    response = requests.get(
        f"{API_BASE_URL}/image/{image_path}",
        headers=headers
    )
    
    print(f"状态码: {response.status_code}")
    print(f"响应内容: {response.text}")
    
    assert response.status_code == 200, "获取图片信息应该成功"
    
    # 验证返回的URL是否正确
    response_data = response.json()
    assert response_data["url"] == upload_url, "返回的URL应该与上传时的URL匹配"

def run_all_tests():
    """运行所有测试"""
    try:
        setup()
        test_upload_image_as_admin()
        test_upload_image_as_user()
        test_get_image_info()
        print("\n所有测试完成！")
    except Exception as e:
        print(f"\n测试过程中出现错误: {str(e)}")

if __name__ == "__main__":
    run_all_tests() 