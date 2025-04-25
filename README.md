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
