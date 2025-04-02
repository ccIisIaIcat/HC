# 健康检查系统 API 文档

本文档详细说明了健康检查系统的所有API接口，包括请求方法、URL、请求参数和响应数据格式，以及对应的Python测试代码。

## 基本信息

- **基础URL**: `http://localhost:8080/api`
- **认证方式**: JWT令牌，在需要认证的接口中，请求头需要包含 `Authorization: Bearer {token}`
- **标准响应格式**:
  ```json
  {
    "code": 200,      // 状态码，200表示成功
    "message": "success", // 状态描述信息
    "data": {         // 实际返回的数据
      // 具体数据字段
    }
  }
  ```

## 用户类型定义

```json
{
  "ID": 1,
  "name": "用户名",
  "email": "example@example.com",
  "role": "user",     // "admin" 或 "user"
  "verified": true,   // 账号是否已验证
  "CreatedAt": "2023-01-01T00:00:00Z", // 创建时间
  "UpdatedAt": "2023-01-01T00:00:00Z", // 更新时间
  "DeletedAt": null   // 删除时间，如果未删除则为null
}
```

## 1. 认证接口

### 1.1 用户注册

- **URL**: `/register`
- **方法**: `POST`
- **描述**: 创建新用户账号并发送验证邮件
- **请求参数**:
  ```json
  {
    "name": "用户名",
    "email": "example@example.com",
    "password": "your_password"
  }
  ```
- **成功响应** (状态码: 200):
  ```json
  {
    "message": "请查收验证邮件完成注册",
    "email": "example@example.com"
  }
  ```
- **错误响应**:
  - 400: 请求格式错误
  - 409: 邮箱已被注册
  - 500: 服务器内部错误

**Python测试代码**:
```python
import requests

def test_register():
    url = 'http://localhost:8080/api/register'
    data = {
        'name': '测试用户',
        'email': 'test@example.com',
        'password': 'password123'
    }
    
    response = requests.post(url, json=data)
    print('注册状态码:', response.status_code)
    print('注册响应:', response.json())
    
    return response.json()
```

### 1.2 用户登录

- **URL**: `/login`
- **方法**: `POST`
- **描述**: 用户登录并获取认证令牌
- **请求参数**:
  ```json
  {
    "email": "example@example.com",
    "password": "your_password"
  }
  ```
- **成功响应** (状态码: 200):
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "ID": 1,
      "name": "用户名",
      "email": "example@example.com",
      "role": "user",
      "verified": true,
      "CreatedAt": "2023-01-01T00:00:00Z",
      "UpdatedAt": "2023-01-01T00:00:00Z",
      "DeletedAt": null
    }
  }
  ```
- **错误响应**:
  - 400: 请求格式错误
  - 401: 邮箱或密码错误
  - 403: 邮箱未验证

**Python测试代码**:
```python
import requests

def test_login():
    url = 'http://localhost:8080/api/login'
    data = {
        'email': 'admin@example.com',  # 使用系统创建的管理员账户
        'password': 'admin123'
    }
    
    response = requests.post(url, json=data)
    print('登录状态码:', response.status_code)
    
    if response.status_code == 200:
        result = response.json()
        print('登录成功，获取令牌:', result['token'])
        return result['token']
    else:
        print('登录失败:', response.text)
        return None
```

### 1.3 邮箱验证

- **URL**: `/verify-email`
- **方法**: `POST`
- **描述**: 验证用户注册时收到的验证码
- **请求参数**:
  ```json
  {
    "email": "example@example.com",
    "code": "abcd1234"
  }
  ```
- **成功响应** (状态码: 200):
  ```json
  {
    "message": "邮箱验证成功",
    "user": {
      "ID": 1,
      "name": "用户名",
      "email": "example@example.com",
      "role": "user",
      "verified": true,
      "CreatedAt": "2023-01-01T00:00:00Z",
      "UpdatedAt": "2023-01-01T00:00:00Z",
      "DeletedAt": null
    }
  }
  ```
- **错误响应**:
  - 400: 请求格式错误或验证码无效
  - 404: 用户不存在
  - 410: 验证码已过期

**Python测试代码**:
```python
import requests

def test_verify_email(email, code):
    url = 'http://localhost:8080/api/verify-email'
    data = {
        'email': email,
        'code': code
    }
    
    response = requests.post(url, json=data)
    print('验证状态码:', response.status_code)
    print('验证响应:', response.json())
    
    return response.json()
```

### 1.4 重新发送验证邮件

- **URL**: `/resend-verification`
- **方法**: `POST`
- **描述**: 重新发送验证邮件
- **请求参数**:
  ```json
  {
    "email": "example@example.com"
  }
  ```
- **成功响应** (状态码: 200):
  ```json
  {
    "message": "验证邮件已重新发送"
  }
  ```
- **错误响应**:
  - 400: 请求格式错误
  - 404: 用户不存在
  - 409: 用户已验证
  - 500: 发送邮件失败

**Python测试代码**:
```python
import requests

def test_resend_verification(email):
    url = 'http://localhost:8080/api/resend-verification'
    data = {
        'email': email
    }
    
    response = requests.post(url, json=data)
    print('重发验证邮件状态码:', response.status_code)
    print('重发验证邮件响应:', response.json())
    
    return response.json()
```

### 1.5 获取当前用户信息

- **URL**: `/me`
- **方法**: `GET`
- **描述**: 获取当前登录用户的信息
- **认证**: 需要JWT令牌
- **请求参数**: 无
- **成功响应** (状态码: 200):
  ```json
  {
    "user": {
      "ID": 1,
      "name": "用户名",
      "email": "example@example.com",
      "role": "user",
      "verified": true,
      "CreatedAt": "2023-01-01T00:00:00Z",
      "UpdatedAt": "2023-01-01T00:00:00Z",
      "DeletedAt": null
    }
  }
  ```
- **错误响应**:
  - 401: 未授权，令牌无效或已过期

**Python测试代码**:
```python
import requests

def test_get_current_user(token):
    url = 'http://localhost:8080/api/me'
    headers = {'Authorization': f'Bearer {token}'}
    
    response = requests.get(url, headers=headers)
    print('获取用户信息状态码:', response.status_code)
    print('用户信息:', response.json())
    
    return response.json()
```

## 2. 食物分析接口

### 2.1 上传并分析食物图片

- **URL**: `/analyze-food`
- **方法**: `POST`
- **描述**: 上传食物图片并使用AI分析食物的营养成分
- **认证**: 需要JWT令牌
- **请求参数**: 使用`multipart/form-data`格式，包含以下字段：
  - `image`: 食物图片文件（支持JPG和PNG，最大10MB）
  - `image_description`: (可选) 图片描述，用于辅助AI分析食物内容
- **成功响应** (状态码: 200):
  ```json
  {
    "hasFood": true,
    "foodType": "汤圆",
    "weight": 150,
    "nutrition": {
      "calories": 350,
      "protein": 5.2,
      "totalFat": 7.5,
      "saturatedFat": 1.2,
      "transFat": 0,
      "unsaturatedFat": 6.3,
      "carbohydrates": 62.8,
      "sugar": 15.3,
      "fiber": 1.5,
      "vitamins": {
        "vitaminA": 0,
        "vitaminC": 0,
        "vitaminD": 0,
        "vitaminE": 0.5,
        "vitaminK": 0,
        "vitaminB": {
          "b1": 0.1,
          "b2": 0.05,
          "b6": 0.1,
          "b12": 0
        }
      },
      "minerals": {
        "calcium": 25,
        "iron": 1.2,
        "sodium": 180,
        "potassium": 120,
        "zinc": 0.5,
        "magnesium": 15
      }
    }
  }
  ```
- **错误响应**:
  - 400: 请求格式错误或不支持的文件类型
  - 401: 未授权，令牌无效或已过期
  - 413: 文件过大
  - 500: 服务器内部错误或分析失败

**Python测试代码**:
```python
import requests

def test_food_analysis(token, image_path, image_description=""):
    url = 'http://localhost:8080/api/analyze-food'
    headers = {'Authorization': f'Bearer {token}'}
    files = {'image': open(image_path, 'rb')}
    data = {}
    
    if image_description:
        data['image_description'] = image_description
    
    response = requests.post(url, headers=headers, files=files, data=data)
    print('食物分析状态码:', response.status_code)
    
    if response.status_code == 200:
        analysis = response.json()
        print(f'食物类型: {analysis["foodType"]}')
        print(f'估计重量: {analysis["weight"]}克')
        print(f'热量: {analysis["nutrition"]["calories"]}卡路里')
        print(f'蛋白质: {analysis["nutrition"]["protein"]}克')
        return analysis
    else:
        print('分析失败:', response.text)
        return None
```

## 3. 管理员接口

### 3.1 获取所有用户列表

- **URL**: `/admin/users`
- **方法**: `GET`
- **描述**: 获取系统中的所有用户列表
- **认证**: 需要JWT令牌，且用户角色为admin
- **请求参数**: 无
- **成功响应** (状态码: 200):
  ```json
  {
    "users": [
      {
        "ID": 1,
        "name": "管理员",
        "email": "admin@example.com",
        "role": "admin",
        "verified": true,
        "CreatedAt": "2023-01-01T00:00:00Z",
        "UpdatedAt": "2023-01-01T00:00:00Z",
        "DeletedAt": null
      },
      {
        "ID": 2,
        "name": "普通用户",
        "email": "user@example.com",
        "role": "user",
        "verified": true,
        "CreatedAt": "2023-01-01T00:00:00Z",
        "UpdatedAt": "2023-01-01T00:00:00Z",
        "DeletedAt": null
      }
    ]
  }
  ```
- **错误响应**:
  - 401: 未授权，令牌无效或已过期
  - 403: 没有管理员权限

**Python测试代码**:
```python
import requests

def test_get_users(admin_token):
    url = 'http://localhost:8080/api/admin/users'
    headers = {'Authorization': f'Bearer {admin_token}'}
    
    response = requests.get(url, headers=headers)
    print('获取用户列表状态码:', response.status_code)
    
    if response.status_code == 200:
        users = response.json()['users']
        print(f'共有 {len(users)} 个用户:')
        for user in users:
            print(f"- {user['name']} ({user['email']}), 角色: {user['role']}")
        return users
    else:
        print('获取用户列表失败:', response.text)
        return None
```

### 3.2 获取系统统计信息

- **URL**: `/admin/stats`
- **方法**: `GET`
- **描述**: 获取系统使用统计信息
- **认证**: 需要JWT令牌，且用户角色为admin
- **请求参数**: 无
- **成功响应** (状态码: 200):
  ```json
  {
    "stats": {
      "userCount": 10,
      "verifiedUserCount": 8,
      "adminCount": 1,
      "analysisCount": 25,
      "activeUsersLast7Days": 5
    }
  }
  ```
- **错误响应**:
  - 401: 未授权，令牌无效或已过期
  - 403: 没有管理员权限

**Python测试代码**:
```python
import requests

def test_get_stats(admin_token):
    url = 'http://localhost:8080/api/admin/stats'
    headers = {'Authorization': f'Bearer {admin_token}'}
    
    response = requests.get(url, headers=headers)
    print('获取统计信息状态码:', response.status_code)
    
    if response.status_code == 200:
        stats = response.json()['stats']
        print('系统统计信息:')
        print(f"- 总用户数: {stats['userCount']}")
        print(f"- 已验证用户数: {stats['verifiedUserCount']}")
        print(f"- 管理员数: {stats['adminCount']}")
        print(f"- 分析次数: {stats['analysisCount']}")
        print(f"- 最近7天活跃用户: {stats['activeUsersLast7Days']}")
        return stats
    else:
        print('获取统计信息失败:', response.text)
        return None
```

## 4. 食物记录接口

以下是食物记录相关的API接口，用于记录、查询、修改和删除用户的饮食记录。

### 4.1 创建食物记录

- **URL**: `/food-records`
- **方法**: `POST`
- **描述**: 手动创建新的食物记录
- **认证**: 需要JWT令牌
- **请求参数**:
  ```json
  {
    "food_name": "测试食物",
    "record_time": "2023-07-15T12:30:00Z",  // 可选，默认为当前时间
    "weight": 200.0,
    
    // 基本营养素
    "calories": 300.0,
    "protein": 10.5,
    "total_fat": 5.0,
    "saturated_fat": 1.2,
    "trans_fat": 0.1,
    "unsaturated_fat": 3.7,
    "carbohydrates": 45.0,
    "sugar": 8.5,
    "fiber": 3.5,
    
    // 维生素（可选）
    "vitamin_a": 120.0,  // μg
    "vitamin_c": 15.0,   // mg
    "vitamin_d": 1.2,    // μg
    "vitamin_b1": 0.3,   // mg
    "vitamin_b2": 0.4,   // mg
    
    // 矿物质（可选）
    "calcium": 45.0,     // mg
    "iron": 2.5,         // mg
    "sodium": 180.0,     // mg
    "potassium": 350.0,  // mg
    
    // 记录信息
    "meal_type": "早餐", // 可选值: 早餐、午餐、晚餐、加餐
    "notes": "手动创建的测试记录",
    "image_path": "/path/to/image.jpg" // 可选
  }
  ```
- **成功响应** (状态码: 200):
  ```json
  {
    "message": "食物记录创建成功",
    "record": {
      "ID": 1,
      "CreatedAt": "2023-07-15T12:30:00Z",
      "UpdatedAt": "2023-07-15T12:30:00Z",
      "DeletedAt": null,
      "user_id": 2,
      "record_time": "2023-07-15T12:30:00Z",
      "food_name": "测试食物",
      "weight": 200.0,
      "calories": 300.0,
      "protein": 10.5,
      // ... 其他字段
    }
  }
  ```
- **错误响应**:
  - 400: 请求数据格式错误
  - 401: 未授权，令牌无效或已过期
  - 500: 服务器内部错误，保存食物记录失败

**Python测试代码**:
```python
import requests
from datetime import datetime

def create_food_record(token, food_data=None):
    url = f'{BASE_URL}/food-records'
    headers = {'Authorization': f'Bearer {token}'}
    
    if food_data is None:
        # 默认测试数据
        food_data = {
            "food_name": "测试食物",
            "record_time": datetime.now().isoformat(),
            "weight": 200.0,
            "calories": 300.0,
            "protein": 10.5,
            "total_fat": 5.0,
            "carbohydrates": 45.0,
            "fiber": 3.5,
            "meal_type": "早餐",
            "notes": "测试记录"
        }
    
    response = requests.post(url, headers=headers, json=food_data)
    print('创建食物记录状态码:', response.status_code)
    
    if response.status_code == 200:
        result = response.json()
        print('食物记录创建成功!')
        print(f'记录ID: {result["record"]["ID"]}')
        return result["record"]
    else:
        print('创建食物记录失败:', response.text)
        return None
```

### 4.2 获取食物记录列表

- **URL**: `/food-records`
- **方法**: `GET`
- **描述**: 获取用户的食物记录列表，支持日期范围过滤
- **认证**: 需要JWT令牌
- **请求参数**:
  - `start_date`: (可选) 开始日期，格式为 `YYYY-MM-DD`
  - `end_date`: (可选) 结束日期，格式为 `YYYY-MM-DD`
- **成功响应** (状态码: 200):
  ```json
  {
    "records": [
      {
        "ID": 1,
        "CreatedAt": "2023-07-15T12:30:00Z",
        "UpdatedAt": "2023-07-15T12:30:00Z",
        "DeletedAt": null,
        "user_id": 2,
        "record_time": "2023-07-15T12:30:00Z",
        "food_name": "测试食物",
        "weight": 200.0,
        // ... 其他字段
      },
      // ... 更多记录
    ],
    "total": 1
  }
  ```
- **错误响应**:
  - 400: 日期格式错误
  - 401: 未授权，令牌无效或已过期
  - 500: 服务器内部错误，获取记录失败

**Python测试代码**:
```python
import requests

def get_food_records(token, start_date=None, end_date=None):
    url = f'{BASE_URL}/food-records'
    headers = {'Authorization': f'Bearer {token}'}
    params = {}
    
    if start_date:
        params['start_date'] = start_date
    if end_date:
        params['end_date'] = end_date
    
    response = requests.get(url, headers=headers, params=params)
    print('获取食物记录列表状态码:', response.status_code)
    
    if response.status_code == 200:
        result = response.json()
        print(f'获取食物记录成功，共 {result["total"]} 条记录')
        return result["records"]
    else:
        print('获取食物记录失败:', response.text)
        return None
```

### 4.3 获取单个食物记录

- **URL**: `/food-records/:id`
- **方法**: `GET`
- **描述**: 获取指定ID的食物记录详情
- **认证**: 需要JWT令牌
- **请求参数**: 
  - `:id`: 记录ID，作为URL参数
- **成功响应** (状态码: 200):
  ```json
  {
    "record": {
      "ID": 1,
      "CreatedAt": "2023-07-15T12:30:00Z",
      "UpdatedAt": "2023-07-15T12:30:00Z",
      "DeletedAt": null,
      "user_id": 2,
      "record_time": "2023-07-15T12:30:00Z",
      "food_name": "测试食物",
      "weight": 200.0,
      "calories": 300.0,
      // ... 其他字段
    }
  }
  ```
- **错误响应**:
  - 400: 无效的记录ID
  - 401: 未授权，令牌无效或已过期
  - 403: 无权访问此记录（不属于当前用户）
  - 404: 记录不存在
  - 500: 服务器内部错误

**Python测试代码**:
```python
import requests

def get_food_record(token, record_id):
    url = f'{BASE_URL}/food-records/{record_id}'
    headers = {'Authorization': f'Bearer {token}'}
    
    response = requests.get(url, headers=headers)
    print('获取单个食物记录状态码:', response.status_code)
    
    if response.status_code == 200:
        result = response.json()
        print(f'成功获取记录: {result["record"]["food_name"]}')
        return result["record"]
    else:
        print('获取食物记录失败:', response.text)
        return None
```

### 4.4 更新食物记录

- **URL**: `/food-records/:id`
- **方法**: `PUT`
- **描述**: 更新指定ID的食物记录
- **认证**: 需要JWT令牌
- **请求参数**: 
  - `:id`: 记录ID，作为URL参数
  - 请求体包含与创建记录相同的JSON字段，可以只包含需要更新的字段
- **成功响应** (状态码: 200):
  ```json
  {
    "message": "记录更新成功",
    "record": {
      "ID": 1,
      "CreatedAt": "2023-07-15T12:30:00Z",
      "UpdatedAt": "2023-07-15T13:45:00Z",
      "DeletedAt": null,
      "user_id": 2,
      "record_time": "2023-07-15T12:30:00Z",
      "food_name": "更新后的食物名称",
      "weight": 250.0,
      // ... 其他更新后的字段
    }
  }
  ```
- **错误响应**:
  - 400: 无效的记录ID或请求数据格式错误
  - 401: 未授权，令牌无效或已过期
  - 403: 无权修改此记录（不属于当前用户）
  - 404: 记录不存在
  - 500: 服务器内部错误，更新失败

**Python测试代码**:
```python
import requests

def update_food_record(token, record_id, update_data):
    url = f'{BASE_URL}/food-records/{record_id}'
    headers = {'Authorization': f'Bearer {token}'}
    
    response = requests.put(url, headers=headers, json=update_data)
    print('更新食物记录状态码:', response.status_code)
    
    if response.status_code == 200:
        result = response.json()
        print('记录更新成功!')
        return result["record"]
    else:
        print('更新食物记录失败:', response.text)
        return None
```

### 4.5 删除食物记录

- **URL**: `/food-records/:id`
- **方法**: `DELETE`
- **描述**: 删除指定ID的食物记录
- **认证**: 需要JWT令牌
- **请求参数**: 
  - `:id`: 记录ID，作为URL参数
- **成功响应** (状态码: 200):
  ```json
  {
    "message": "记录已成功删除"
  }
  ```
- **错误响应**:
  - 400: 无效的记录ID
  - 401: 未授权，令牌无效或已过期
  - 403: 无权删除此记录（不属于当前用户）
  - 404: 记录不存在
  - 500: 服务器内部错误，删除失败

**Python测试代码**:
```python
import requests

def delete_food_record(token, record_id):
    url = f'{BASE_URL}/food-records/{record_id}'
    headers = {'Authorization': f'Bearer {token}'}
    
    response = requests.delete(url, headers=headers)
    print('删除食物记录状态码:', response.status_code)
    
    if response.status_code == 200:
        print('记录删除成功!')
        return True
    else:
        print('删除食物记录失败:', response.text)
        return False
```

### 4.6 分析并保存食物记录

- **URL**: `/analyze-and-save`
- **方法**: `POST`
- **描述**: 上传食物图片进行分析，并将分析结果保存为食物记录
- **认证**: 需要JWT令牌
- **请求参数**: 使用`multipart/form-data`格式，包含以下字段：
  - `image`: 食物图片文件（支持JPG和PNG，最大10MB）
  - `meal_type`: 餐食类型（早餐/午餐/晚餐/加餐）
  - `notes`: (可选) 备注信息
  - `image_description`: (可选) 图片描述，用于辅助AI分析食物内容
- **成功响应** (状态码: 200):
  ```json
  {
    "message": "食物分析和记录保存成功",
    "record": {
      "ID": 1,
      "CreatedAt": "2023-07-15T12:30:00Z",
      "UpdatedAt": "2023-07-15T12:30:00Z",
      "DeletedAt": null,
      "user_id": 2,
      "record_time": "2023-07-15T12:30:00Z",
      "food_name": "汤圆",
      "weight": 150.0,
      "calories": 350.0,
      // ... 其他字段
    },
    "analysis": {
      "hasFood": true,
      "foodType": "汤圆",
      "weight": 150,
      "nutrition": {
        // ... 详细的分析结果
      }
    }
  }
  ```
- **错误响应**:
  - 400: 请求格式错误、不支持的文件类型、未能检测到食物
  - 401: 未授权，令牌无效或已过期
  - 413: 文件过大
  - 500: 服务器内部错误，分析或保存失败

**Python测试代码**:
```python
import requests
import os

def analyze_and_save_food(token, image_path, meal_type="午餐", notes="", image_description=""):
    url = f'{BASE_URL}/analyze-and-save'
    headers = {'Authorization': f'Bearer {token}'}
    
    if not os.path.exists(image_path):
        print(f'错误: 图片文件不存在 {image_path}')
        return None
    
    with open(image_path, 'rb') as img:
        files = {'image': img}
        data = {
            'meal_type': meal_type,
            'notes': notes,
            'image_description': image_description
        }
        response = requests.post(url, headers=headers, files=files, data=data)
    
    print('分析并保存食物记录状态码:', response.status_code)
    
    if response.status_code == 200:
        result = response.json()
        print('食物分析和记录保存成功!')
        print(f'记录ID: {result["record"]["ID"]}')
        print(f'食物名称: {result["record"]["food_name"]}')
        print(f'卡路里: {result["record"]["calories"]}')
        return result
    else:
        print('分析并保存食物记录失败:', response.text)
        return None
```

### 4.7 食物记录数据结构

食物记录包含以下字段：

```json
{
  "ID": 1,                           // 记录唯一ID（自动生成）
  "CreatedAt": "2023-07-15T12:30:00Z", // 创建时间（自动生成）
  "UpdatedAt": "2023-07-15T12:30:00Z", // 更新时间（自动生成）
  "DeletedAt": null,                 // 删除时间（软删除，默认为null）
  "user_id": 2,                      // 用户ID（关联User表）
  "record_time": "2023-07-15T12:30:00Z", // 记录时间（用餐时间）
  "food_name": "汤圆",                // 食物名称
  "weight": 150.0,                   // 食物重量(克)
  
  // 基本营养素
  "calories": 350.0,                 // 热量（卡路里）
  "protein": 5.2,                    // 蛋白质含量（克）
  "total_fat": 7.5,                  // 总脂肪（克）
  "saturated_fat": 1.2,              // 饱和脂肪（克）
  "trans_fat": 0.0,                  // 反式脂肪（克）
  "unsaturated_fat": 6.3,            // 不饱和脂肪（克）
  "carbohydrates": 62.8,             // 碳水化合物（克）
  "sugar": 15.3,                     // 糖分（克）
  "fiber": 1.5,                      // 膳食纤维（克）
  
  // 维生素
  "vitamin_a": 0.0,                  // 维生素A（μg）
  "vitamin_c": 0.0,                  // 维生素C（mg）
  "vitamin_d": 0.0,                  // 维生素D（μg）
  "vitamin_b1": 0.1,                 // 维生素B1（mg）
  "vitamin_b2": 0.05,                // 维生素B2（mg）
  
  // 矿物质
  "calcium": 25.0,                   // 钙（mg）
  "iron": 1.2,                       // 铁（mg）
  "sodium": 180.0,                   // 钠（mg）
  "potassium": 120.0,                // 钾（mg）
  
  // 记录类型和备注
  "meal_type": "午餐",                // 餐食类型：早餐/午餐/晚餐/加餐
  "notes": "好吃的汤圆",              // 备注
  "image_path": "/uploads/foods/1.jpg" // 图片路径（可选）
}
```

## 完整测试代码

以下是一个完整的Python测试脚本，可以用来测试所有API接口：

```python
import requests
import os

# 服务器基础URL
BASE_URL = 'http://localhost:8080/api'

# 测试管理员登录
def login_admin():
    url = f'{BASE_URL}/login'
    data = {
        'email': 'admin@example.com',
        'password': 'admin123'
    }
    
    response = requests.post(url, json=data)
    if response.status_code == 200:
        return response.json()['token']
    else:
        print(f'管理员登录失败: {response.text}')
        return None

# 测试注册新用户
def register_user(name, email, password):
    url = f'{BASE_URL}/register'
    data = {
        'name': name,
        'email': email,
        'password': password
    }
    
    response = requests.post(url, json=data)
    print(f'用户注册状态码: {response.status_code}')
    print(f'响应: {response.text}')
    return response.status_code == 200

# 测试验证邮箱
def verify_email(email, code):
    url = f'{BASE_URL}/verify-email'
    data = {
        'email': email,
        'code': code
    }
    
    response = requests.post(url, json=data)
    print(f'邮箱验证状态码: {response.status_code}')
    print(f'响应: {response.text}')
    return response.status_code == 200

# 测试用户登录
def login_user(email, password):
    url = f'{BASE_URL}/login'
    data = {
        'email': email,
        'password': password
    }
    
    response = requests.post(url, json=data)
    if response.status_code == 200:
        return response.json()['token']
    else:
        print(f'用户登录失败: {response.text}')
        return None

# 测试获取当前用户信息
def get_current_user(token):
    url = f'{BASE_URL}/me'
    headers = {'Authorization': f'Bearer {token}'}
    
    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        return response.json()['user']
    else:
        print(f'获取用户信息失败: {response.text}')
        return None

# 测试食物分析
def analyze_food(token, image_path, image_description=""):
    url = 'http://localhost:8080/api/analyze-food'
    headers = {'Authorization': f'Bearer {token}'}
    files = {'image': open(image_path, 'rb')}
    data = {}
    
    if image_description:
        data['image_description'] = image_description
    
    response = requests.post(url, headers=headers, files=files, data=data)
    print('食物分析状态码:', response.status_code)
    
    if response.status_code == 200:
        analysis = response.json()
        print(f'食物类型: {analysis["foodType"]}')
        print(f'估计重量: {analysis["weight"]}克')
        print(f'热量: {analysis["nutrition"]["calories"]}卡路里')
        print(f'蛋白质: {analysis["nutrition"]["protein"]}克')
        return analysis
    else:
        print('分析失败:', response.text)
        return None

# 测试获取所有用户列表 (管理员功能)
def get_all_users(admin_token):
    url = f'{BASE_URL}/admin/users'
    headers = {'Authorization': f'Bearer {admin_token}'}
    
    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        return response.json()['users']
    else:
        print(f'获取用户列表失败: {response.text}')
        return None

# 测试获取系统统计信息 (管理员功能)
def get_system_stats(admin_token):
    url = f'{BASE_URL}/admin/stats'
    headers = {'Authorization': f'Bearer {admin_token}'}
    
    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        return response.json()['stats']
    else:
        print(f'获取统计信息失败: {response.text}')
        return None

# 运行测试
if __name__ == '__main__':
    print("开始测试健康检查系统API...")
    
    # 测试管理员登录
    print("\n1. 测试管理员登录")
    admin_token = login_admin()
    if admin_token:
        print("管理员登录成功!")
        
        # 测试获取管理员信息
        admin_info = get_current_user(admin_token)
        if admin_info:
            print(f"获取管理员信息成功: {admin_info['name']} ({admin_info['email']})")
        
        # 测试获取所有用户
        print("\n2. 测试获取所有用户")
        users = get_all_users(admin_token)
        if users:
            print(f"获取用户列表成功, 共 {len(users)} 个用户")
        
        # 测试获取系统统计
        print("\n3. 测试获取系统统计")
        stats = get_system_stats(admin_token)
        if stats:
            print("获取系统统计成功!")
    
    # 测试注册新用户
    print("\n4. 测试注册新用户")
    test_email = "test_user@example.com"
    register_success = register_user("测试用户", test_email, "password123")
    
    if register_success:
        print("用户注册成功，请检查收到的验证邮件")
        verification_code = input("请输入收到的验证码: ")
        
        # 测试验证邮箱
        print("\n5. 测试验证邮箱")
        verify_success = verify_email(test_email, verification_code)
        
        if verify_success:
            print("邮箱验证成功!")
            
            # 测试新用户登录
            print("\n6. 测试新用户登录")
            user_token = login_user(test_email, "password123")
            
            if user_token:
                print("用户登录成功!")
                
                # 测试获取用户信息
                user_info = get_current_user(user_token)
                if user_info:
                    print(f"获取用户信息成功: {user_info['name']} ({user_info['email']})")
                
                # 测试食物分析
                print("\n7. 测试食物分析")
                image_path = input("请输入食物图片路径: ")
                image_description = input("请输入食物描述: ")
                analysis_result = analyze_food(user_token, image_path, image_description)
                
                if analysis_result:
                    print("食物分析成功!")
                    print(f"食物类型: {analysis_result['foodType']}")
                    print(f"估计重量: {analysis_result['weight']}克")
                    print(f"热量: {analysis_result['nutrition']['calories']}卡路里")
    
    print("\n测试完成!") 