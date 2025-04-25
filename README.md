# 饭团猫 - 健康饮食助手

饭团猫是一个全面的健康饮食管理应用，帮助用户记录和分析日常饮食，追踪营养摄入，实现健康生活目标。

## 功能特点

- 🍲 **食物识别与分析**：AI 自动识别食物图片并分析营养成分
- 📊 **营养数据追踪**：记录并分析每日营养摄入数据
- 📈 **健康状态监测**：追踪体重、BMI等身体指标变化
- 🏆 **成就系统**：完成健康挑战获得成就徽章
- 👥 **社区互动**：分享健康食谱和经验
- 🔔 **健康提醒**：定制个性化的饮食和健康提醒

## 技术栈

### 后端
- Go
- Gin Web 框架
- GORM
- MySQL
- Redis（可选，用于缓存）

### 前端
- uni-app
- Vue.js
- TypeScript
- Vite

## 项目结构

```
HC/
├── backend/             # 后端代码
│   ├── handlers/        # 请求处理器
│   ├── models/          # 数据模型
│   ├── utils/          # 工具函数
│   └── main.go         # 主程序入口
│
├── frontend_uniapp/     # 前端代码
│   ├── src/
│   │   ├── api/        # API 接口
│   │   ├── components/ # 组件
│   │   ├── pages/      # 页面
│   │   └── static/     # 静态资源
│   └── package.json
│
└── doc/                # 文档
```

## 部署说明

### 后端部署

1. 安装依赖
```bash
# 安装 Go (如果未安装)
# Windows: 访问 https://golang.org/dl/ 下载安装包
# Linux:
sudo apt-get update
sudo apt-get install golang-go

# 安装依赖
cd backend
go mod tidy
```

2. 配置数据库
```bash
# 安装 MySQL (如果未安装)
# 创建数据库
mysql -u root -p
CREATE DATABASE healthcheck;
```

3. 配置环境变量
```bash
# 创建 .env 文件
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=healthcheck
JWT_SECRET=your_jwt_secret
```

4. 运行服务
```bash
go run main.go
```

### 前端部署

1. 安装依赖
```bash
cd frontend_uniapp
npm install
```

2. 配置环境
```bash
# 创建 .env 文件
VITE_API_BASE_URL=http://localhost:8080
```

3. 开发模式
```bash
npm run dev
```

4. 构建生产版本
```bash
npm run build
```

### 应用更新部署

1. 上传新版本
```bash
# 使用管理脚本上传新版本
cd admin_script
python upload_apk.py --version "1.0.0" --file "path/to/app.apk"
```

2. 检查更新
- 应用会自动检查新版本
- 访问 `/download` 页面获取最新版本

## API 文档

API 文档见 `doc/api.md`

## 开发指南

### 添加新的帖子
```go
post := &models.Post{
    Title:   "标题",
    Content: "内容",
    Author:  "作者",
}

imageURLs := []string{
    "http://example.com/image1.jpg",
    "http://example.com/image2.jpg",
}

err := models.CreatePostWithImages(post, imageURLs)
```

### 获取帖子
```go
post, err := models.GetPostWithImages(postID)
```

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交改动 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

## 许可证

版权所有 © 2023-2024 饭团猫团队

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