import requests
import json
import os
import sys
import time
from datetime import datetime, timedelta

# API基础URL
BASE_URL = 'http://localhost:8080/api'

# 全局存储
session = requests.Session()
test_record_id = None  # 用于存储创建的食物记录ID


def print_result(name, success):
    """打印测试结果"""
    mark = '✓' if success else '✗'
    print(f"{mark} {name}")


def login():
    """登录并获取认证令牌"""
    print("\n===== 登录获取Token =====")
    try:
        credentials = {
            "email": "965377515@qq.com",
            "password": "123456"
        }
        
        # 先清除之前的会话信息
        session.cookies.clear()
        session.headers.pop('Authorization', None)
        
        response = session.post(f"{BASE_URL}/login", json=credentials)
        print(f"登录响应状态码: {response.status_code}")
        
        if response.status_code == 200:
            try:
                response_json = response.json()
                print(f"登录响应JSON: {json.dumps(response_json, indent=2, ensure_ascii=False)}")
                
                # 1. 尝试从JSON响应中提取token
                token = None
                
                # 可能的token位置
                if 'token' in response_json:
                    token = response_json.get('token')
                elif 'data' in response_json and 'token' in response_json.get('data', {}):
                    token = response_json.get('data', {}).get('token')
                else:
                    # 查看响应中是否有其他可能的token字段
                    for key, value in response_json.items():
                        if 'token' in key.lower() and isinstance(value, str):
                            token = value
                            break
                            
                if token:
                    # 将token设置到session的headers中
                    session.headers.update({
                        'Authorization': f'Bearer {token}'
                    })
                    print(f"已将Bearer Token设置到请求头: {token[:10]}...（已截断）")
                
                # 2. 检查 Cookie 认证
                auth_cookie_found = False
                for name, value in session.cookies.items():
                    if any(keyword in name.lower() for keyword in ['token', 'auth', 'jwt', 'session']):
                        auth_cookie_found = True
                        print(f"发现认证相关cookie: {name}")
                
                if token or auth_cookie_found:
                    print_result("认证信息设置成功", True)
                    return True
                else:
                    # 3. 假设使用会话cookie认证
                    print("未发现明确的token或认证cookie，将假设使用会话认证")
                    print_result("使用会话认证", True)
                    return True
                    
            except json.JSONDecodeError:
                print("响应不是有效的JSON格式")
                print(f"原始响应内容: {response.text}")
                # 尝试使用会话cookie继续
                print_result("假设使用会话Cookie认证", True)
                return True
        else:
            print_result("登录失败", False)
            print(f"状态码: {response.status_code}, 响应: {response.text}")
            return False
    except Exception as e:
        print_result("登录请求异常", False)
        print(f"异常: {str(e)}")
        return False


def api_request(method, endpoint, **kwargs):
    """
    通用API请求函数，带有自动认证刷新功能
    method: 请求方法（'get', 'post', 'put', 'delete'等）
    endpoint: API端点路径
    **kwargs: 传递给requests的其他参数
    """
    url = f"{BASE_URL}/{endpoint.lstrip('/')}"
    
    # 第一次请求
    response = getattr(session, method.lower())(url, **kwargs)
    
    # 如果遇到401错误，尝试刷新认证
    if response.status_code == 401:
        print(f"请求 {method.upper()} {endpoint} 收到401未授权错误，尝试刷新认证...")
        if try_refresh_auth():
            print("认证已刷新，重新尝试请求")
            # 重新发送请求
            response = getattr(session, method.lower())(url, **kwargs)
        else:
            print("认证刷新失败")
    
    return response


def test_analyze_food():
    """测试食物分析API"""
    print("\n===== 测试食物分析 =====")
    try:
        # 使用根目录下的tangyuan.png
        image_path = "tangyuan.png"
        if not os.path.exists(image_path):
            print_result("找不到tangyuan.png图片", False)
            print("请确保根目录下有tangyuan.png文件")
            return None
        
        # 创建表单数据
        files = {
            'image': ('tangyuan.png', open(image_path, 'rb'), 'image/png')
        }
        
        # 询问是否提供图片描述
        image_description = input("请输入食物图片描述（可选，直接回车跳过）：")
        data = {}
        if image_description.strip():
            data['image_description'] = image_description
            print(f"使用图片描述: {image_description}")
        
        # 使用通用API请求函数
        response = api_request('post', 'analyze-food', files=files, data=data)
        
        if response.status_code == 200:
            result = response.json()
            print_result("食物分析成功", True)
            print(f"分析结果: {json.dumps(result, indent=2, ensure_ascii=False)}")
            return result
        else:
            print_result("食物分析失败", False)
            print(f"状态码: {response.status_code}, 响应: {response.text}")
            return None
    except Exception as e:
        print_result("食物分析请求异常", False)
        print(f"异常: {str(e)}")
        return None


def test_analyze_and_save_food(description=None):
    """测试分析并保存食物记录API"""
    print("\n===== 测试分析并保存食物记录 =====")
    
    # 先获取token
    token = get_token()
    if not token:
        print("❌ 无法获取token，测试失败")
        return
    
    # 检查测试图片是否存在
    if not os.path.exists(TEST_IMAGE_PATH):
        print(f"❌ 错误: 图片文件不存在 {TEST_IMAGE_PATH}")
        return None
    
    url = f'{BASE_URL}/analyze-and-save-food'
    headers = {'Authorization': f'Bearer {token}'}
    
    with open(TEST_IMAGE_PATH, 'rb') as img:
        files = {'image': img}
        data = {
            'meal_type': '午餐',
            'notes': '测试添加的食物记录'
        }
        
        if description:
            data['image_description'] = description
            print(f"使用图片描述: {description}")
        
        response = requests.post(url, headers=headers, files=files, data=data)
    
    print_response(response, "分析并保存食物记录")
    
    if response.status_code == 200:
        result = response.json()
        print("✅ 分析并保存食物记录成功")
        print(f"- 记录ID: {result.get('record_id', 'N/A')}")
        
        analysis = result.get('analysis', {})
        if analysis.get('hasFood'):
            print(f"- 食物类型: {analysis.get('foodType', 'N/A')}")
            print(f"- 估计重量: {analysis.get('weight', 'N/A')}克")
            
            nutrition = analysis.get('nutrition', {})
            print(f"- 热量: {nutrition.get('calories', 'N/A')}卡路里")
            print(f"- 蛋白质: {nutrition.get('protein', 'N/A')}克")
        return result
    
    print("❌ 分析并保存食物记录失败")
    return None


def test_create_food_record():
    """测试创建食物记录API"""
    print("\n===== 测试创建食物记录 =====")
    try:
        # 构建食物记录数据
        now = datetime.now()
        record_data = {
            "food_name": "测试食物",
            "record_time": now.strftime("%Y-%m-%dT%H:%M:%SZ"),
            "weight": 200.0,
            
            # 基本营养素
            "calories": 300.0,
            "protein": 10.5,
            "total_fat": 5.0,
            "saturated_fat": 1.2,
            "trans_fat": 0.1,
            "unsaturated_fat": 3.7,
            "carbohydrates": 45.0,
            "sugar": 8.5,
            "fiber": 3.5,
            
            # 维生素
            "vitamin_a": 120.0,
            "vitamin_c": 15.0,
            "vitamin_d": 1.2,
            "vitamin_b1": 0.3,
            "vitamin_b2": 0.4,
            
            # 矿物质
            "calcium": 45.0,
            "iron": 2.5,
            "sodium": 180.0,
            "potassium": 350.0,
            
            # 记录信息
            "meal_type": "早餐",
            "notes": "手动创建的完整测试记录"
        }
        
        # 使用通用API请求函数
        response = api_request('post', 'food-records', json=record_data)
        
        if response.status_code == 200:
            result = response.json()
            global test_record_id
            if not test_record_id:  # 如果之前的测试没有设置ID，使用这个
                test_record_id = result.get('record', {}).get('ID')
            print_result("创建食物记录成功", True)
            print(f"记录ID: {result.get('record', {}).get('ID')}")
            return result
        else:
            print_result("创建食物记录失败", False)
            print(f"状态码: {response.status_code}, 响应: {response.text}")
            return None
    except Exception as e:
        print_result("创建食物记录请求异常", False)
        print(f"异常: {str(e)}")
        return None


def test_get_food_records():
    """测试获取食物记录列表API"""
    print("\n===== 测试获取食物记录列表 =====")
    try:
        # 获取过去30天的记录
        today = datetime.now().strftime("%Y-%m-%d")
        month_ago = (datetime.now() - timedelta(days=30)).strftime("%Y-%m-%d")
        
        # 使用通用API请求函数
        response = api_request('get', f"food-records?start_date={month_ago}&end_date={today}")
        
        if response.status_code == 200:
            result = response.json()
            print_result("获取食物记录列表成功", True)
            print(f"记录总数: {result.get('total', 0)}")
            if result.get('total', 0) > 0:
                record = result.get('records', [])[0]
                print(f"示例记录: {record.get('food_name')} - {record.get('calories')} 卡路里")
            return result
        else:
            print_result("获取食物记录列表失败", False)
            print(f"状态码: {response.status_code}, 响应: {response.text}")
            return None
    except Exception as e:
        print_result("获取食物记录列表请求异常", False)
        print(f"异常: {str(e)}")
        return None


def test_get_food_record():
    """测试获取单个食物记录API"""
    print("\n===== 测试获取单个食物记录 =====")
    try:
        global test_record_id
        if not test_record_id:
            print_result("获取单个食物记录失败", False)
            print("没有可用的记录ID，请先创建记录")
            return None
            
        # 使用通用API请求函数
        response = api_request('get', f"food-records/{test_record_id}")
        
        if response.status_code == 200:
            result = response.json()
            print_result("获取单个食物记录成功", True)
            print(f"食物名称: {result.get('record', {}).get('food_name')}")
            return result
        else:
            print_result("获取单个食物记录失败", False)
            print(f"状态码: {response.status_code}, 响应: {response.text}")
            return None
    except Exception as e:
        print_result("获取单个食物记录请求异常", False)
        print(f"异常: {str(e)}")
        return None


def test_update_food_record():
    """测试更新食物记录API"""
    print("\n===== 测试更新食物记录 =====")
    try:
        global test_record_id
        if not test_record_id:
            print_result("更新食物记录失败", False)
            print("没有可用的记录ID，请先创建记录")
            return None
        
        # 先获取原始记录
        response = api_request('get', f"food-records/{test_record_id}")
        if response.status_code != 200:
            print_result("获取原始记录失败", False)
            print(f"状态码: {response.status_code}, 响应: {response.text}")
            return None
            
        original_record = response.json().get('record', {})
        
        # 修改记录数据
        update_data = original_record.copy()
        update_data["notes"] = f"更新的测试记录 {datetime.now().strftime('%H:%M:%S')}"
        update_data["calories"] = float(original_record.get("calories", 0)) + 50.0
        
        # 使用通用API请求函数
        response = api_request('put', f"food-records/{test_record_id}", json=update_data)
        
        if response.status_code == 200:
            result = response.json()
            print_result("更新食物记录成功", True)
            print(f"更新后的热量: {result.get('record', {}).get('calories')}")
            return result
        else:
            print_result("更新食物记录失败", False)
            print(f"状态码: {response.status_code}, 响应: {response.text}")
            return None
    except Exception as e:
        print_result("更新食物记录请求异常", False)
        print(f"异常: {str(e)}")
        return None


def test_delete_food_record():
    """测试删除食物记录API"""
    print("\n===== 测试删除食物记录 =====")
    try:
        global test_record_id
        if not test_record_id:
            print_result("删除食物记录失败", False)
            print("没有可用的记录ID，请先创建记录")
            return None
            
        # 使用通用API请求函数
        response = api_request('delete', f"food-records/{test_record_id}")
        
        if response.status_code == 200:
            result = response.json()
            print_result("删除食物记录成功", True)
            print(f"响应消息: {result.get('message')}")
            
            # 验证记录是否真的被删除
            verify_response = api_request('get', f"food-records/{test_record_id}")
            if verify_response.status_code == 404:
                print_result("验证删除成功", True)
            else:
                print_result("验证删除失败", False)
                print(f"记录可能未被删除，状态码: {verify_response.status_code}")
            
            return result
        else:
            print_result("删除食物记录失败", False)
            print(f"状态码: {response.status_code}, 响应: {response.text}")
            return None
    except Exception as e:
        print_result("删除食物记录请求异常", False)
        print(f"异常: {str(e)}")
        return None


def check_auth_cookies():
    """检查是否有认证相关的cookie并打印出来"""
    print("\n===== 检查认证Cookie =====")
    cookies_dict = session.cookies.get_dict()
    
    # 打印所有cookie
    print(f"当前会话的所有Cookie: {cookies_dict}")
    
    # 检查是否存在可能的认证相关cookie
    auth_keywords = ['token', 'auth', 'jwt', 'session']
    found_auth_cookie = False
    
    for name, value in cookies_dict.items():
        for keyword in auth_keywords:
            if keyword.lower() in name.lower():
                print(f"发现可能的认证Cookie: {name}={value[:5]}...")
                found_auth_cookie = True
    
    return found_auth_cookie


def send_test_request():
    """发送测试请求，检查认证状态"""
    print("\n===== 发送测试认证请求 =====")
    try:
        # 获取当前用户信息作为测试
        response = api_request('get', 'me')
        
        print(f"测试请求状态码: {response.status_code}")
        print(f"测试请求响应: {response.text}")
        
        if response.status_code == 200:
            print_result("认证状态检查成功", True)
            return True
        else:
            print_result("认证状态检查失败", False)
            return False
    except Exception as e:
        print_result("测试请求异常", False)
        print(f"异常: {str(e)}")
        return False


def try_refresh_auth():
    """如果认证失败，尝试刷新token或重新登录"""
    print("\n===== 尝试刷新认证 =====")
    
    # 首先尝试发送测试请求检查认证状态
    try:
        test_response = session.get(f"{BASE_URL}/me")
        if test_response.status_code == 200:
            print_result("认证有效，无需刷新", True)
            return True
    except Exception:
        pass  # 忽略错误，继续刷新
    
    print("认证可能已失效，尝试刷新...")
    
    # 尝试方法1: 查找并使用刷新token端点
    try:
        refresh_response = session.post(f"{BASE_URL}/refresh-token")
        if refresh_response.status_code == 200:
            try:
                refresh_data = refresh_response.json()
                if 'token' in refresh_data:
                    token = refresh_data['token']
                    session.headers.update({'Authorization': f'Bearer {token}'})
                    print_result("成功刷新token", True)
                    return True
            except:
                pass
    except Exception:
        pass  # 忽略错误，尝试下一种方法
    
    # 如果刷新失败，尝试重新登录
    print("刷新token失败，尝试重新登录...")
    return login()


def run_all_tests():
    """运行所有测试"""
    print("======== 健康检查系统API测试 - 食物记录功能 ========")
    
    # 登录获取认证
    if not login():
        print("登录失败，无法继续测试")
        return
    
    # 检查认证cookie
    check_auth_cookies()
    
    # 发送测试请求
    if not send_test_request():
        print("认证状态检查失败，尝试刷新认证...")
        if not try_refresh_auth():
            print("刷新认证失败，无法继续测试")
            return
    
    # 测试食物分析
    test_analyze_food()
    
    # 测试带有图片描述的食物分析
    test_food_analysis_with_description()
    
    # 测试分析并保存
    test_analyze_and_save_food()
    
    # 测试带有图片描述的分析并保存
    test_analyze_and_save_food_with_description()
    
    # 测试创建记录
    test_create_food_record()
    
    # # 测试获取记录列表
    # test_get_food_records()
    
    # # 测试获取单个记录
    # test_get_food_record()
    
    # # 测试更新记录
    # test_update_food_record()
    
    # # 测试删除记录
    # test_delete_food_record()
    
    print("\n======== 测试完成 ========")


# 添加一个测试图片描述的函数
def test_food_analysis_with_description():
    """测试带有图片描述的食物分析接口"""
    print("\n===== 测试带有图片描述的食物分析 =====")
    
    # 先获取token
    token = get_token()
    if not token:
        print("❌ 无法获取token，测试失败")
        return
    
    # 检查测试图片是否存在
    if not os.path.exists(TEST_IMAGE_PATH):
        print(f"❌ 错误: 图片文件不存在 {TEST_IMAGE_PATH}")
        return None
    
    # 提供图片描述
    description = "这是一盘炒饭，里面有鸡蛋、胡萝卜、青豆和火腿丁"
    
    url = f'{BASE_URL}/analyze-food'
    headers = {'Authorization': f'Bearer {token}'}
    
    with open(TEST_IMAGE_PATH, 'rb') as img:
        files = {'image': img}
        data = {'image_description': description}
        print(f"使用图片描述: {description}")
        
        response = requests.post(url, headers=headers, files=files, data=data)
    
    print_response(response, "带描述的食物分析")
    
    if response.status_code == 200:
        result = response.json()
        print("✅ 带描述的食物分析成功")
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
    
    print("❌ 带描述的食物分析失败")
    return None


# 添加一个新的测试函数，测试带有图片描述的分析并保存
def test_analyze_and_save_food_with_description():
    """测试带有图片描述的分析并保存食物记录API"""
    print("\n===== 测试带有图片描述的分析并保存食物记录 =====")
    
    # 提供图片描述
    description = "这是一盘蔬菜沙拉，包含生菜、番茄、黄瓜和少量橄榄油"
    
    # 调用测试函数并传递描述
    return test_analyze_and_save_food(description)


if __name__ == "__main__":
    print("\n========== 食物分析API测试开始 ==========\n")
    
    # 测试基本食物分析
    test_food_analysis()
    
    # 测试带有图片描述的食物分析
    test_food_analysis_with_description()
    
    # 测试分析并保存食物记录
    test_analyze_and_save_food()
    
    # 测试带有图片描述的分析并保存食物记录
    test_analyze_and_save_food_with_description()
    
    print("\n========== 食物分析API测试完成 ==========")