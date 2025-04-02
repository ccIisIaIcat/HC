# 健康检查系统 - 前端API模块

这个项目包含了健康检查系统的前端API模块，用于与后端服务进行通信。

## 目录结构

```
src/
├── api/          # API模块目录
│   ├── index.js  # API导出文件
│   ├── axios.js  # Axios实例配置
│   ├── auth.js   # 认证相关API
│   ├── user.js   # 用户相关API
│   └── food.js   # 食物分析相关API
└── tests/        # 测试目录
    ├── api.test.js      # 自动化单元测试
    └── api.manual.test.js  # 手动测试工具
```

## API模块

API模块被拆分为以下几个部分:

1. `axios.js` - 创建和配置axios实例，包含请求/响应拦截器
2. `auth.js` - 处理用户认证相关的API请求
3. `user.js` - 处理用户管理(管理员)相关的API请求
4. `food.js` - 处理食物分析相关的API请求
5. `index.js` - 导出所有API模块供应用使用

## 如何使用

### 安装依赖

```bash
npm install
```

### 导入和使用

```javascript
// 导入整个API模块
import api from './api';

// 调用认证API
api.auth.login({email: 'user@example.com', password: 'password'})
  .then(response => {
    console.log('登录成功:', response);
  })
  .catch(error => {
    console.error('登录失败:', error);
  });

// 也可以单独导入特定模块
import { auth, food } from './api';

// 分析食物图片
const fileInput = document.querySelector('input[type="file"]');
fileInput.addEventListener('change', async (e) => {
  const file = e.target.files[0];
  try {
    const result = await food.analyzeFood(file);
    console.log('食物分析结果:', result);
  } catch (error) {
    console.error('分析失败:', error);
  }
});
```

## 测试

项目包含两种测试方式:

### 自动化单元测试

使用Jest和axios-mock-adapter进行单元测试:

```bash
npm test
```

### 手动测试工具

在浏览器控制台中可以运行手动测试:

```javascript
import { runManualTests } from './tests/api.manual.test';

// 运行所有手动测试
runManualTests();

// 或者单独运行某个测试
import manualTests from './tests/api.manual.test';
manualTests.testAuth();
```

## 后端API文档

完整的后端API文档见`API文档.md`文件。

# DeepSeek V3 API测试脚本

这是一个简单的DeepSeek V3 API测试脚本，发送"你好"并显示返回结果。

## 安装依赖

```bash
pip install requests
```

## 使用方法

### 基本用法

```bash
python test_deepseek_v3.py --api-key YOUR_API_KEY
```

或者自定义参数：

```bash
python test_deepseek_v3.py --api-key YOUR_API_KEY --prompt "你想说的话" --model "deepseek-chat" --api-url "https://api.deepseek.com/v1/chat/completions"
```

### 可用参数

- `--api-key`: DeepSeek API密钥
- `--prompt`: 输入提示词，默认为"你好"
- `--model`: 使用的模型，默认为"deepseek-chat"（即DeepSeek-V3）
- `--api-url`: API端点URL，默认为"https://api.deepseek.com/v1/chat/completions"

### API密钥设置

脚本会按以下顺序查找API密钥：
1. 命令行参数 `--api-key`
2. 环境变量 `DEEPSEEK_API_KEY`
3. 交互式提示输入 