# 文章分析系统

一个基于Vue.js和Go开发的智能文章分析系统，支持文件上传、AI分析和结果展示。

## 系统功能

- 📄 **文件上传** - 支持txt、pdf、doc、docx格式文件
- 🤖 **AI智能分析** - 基于大模型技术进行文章深度分析
- 📊 **分析结果展示** - 提供摘要、关键要点、情感分析等
- 🔍 **搜索功能** - 支持文章标题、内容和作者的模糊搜索
- 📱 **响应式设计** - 适配各种设备屏幕

## 技术栈

### 前端
- Vue.js 3 + TypeScript
- Element Plus UI组件库
- Vite构建工具
- Axios HTTP客户端

### 后端
- Go语言
- Gin Web框架
- GORM ORM框架
- SQLite数据库
- 大模型API集成

## 项目结构

```
article_analysis/
├── backend/          # 后端服务
│   ├── cmd/         # 主程序入口
│   ├── internal/    # 内部模块
│   ├── data/        # 数据库文件
│   └── scripts/     # SQL脚本
├── frontend/        # 前端应用
│   ├── src/         # 源代码
│   ├── public/      # 静态资源
│   └── package.json # 依赖配置
└── 系统方案文档.md   # 系统设计方案
```

## 快速启动

### 前置要求
- Node.js 16+ 
- Go 1.19+
- Git

### 1. 克隆项目
```bash
git clone <repository-url>
cd article_analysis
```

### 2. 启动后端服务
```bash
cd backend
go mod download
go run cmd/main.go
```
后端服务将运行在 http://localhost:8080

### 3. 启动前端服务
```bash
cd frontend
npm install
npm run dev
```
前端应用将运行在 http://localhost:5173

### 4. 访问系统
打开浏览器访问 http://localhost:5173 即可使用系统

## 使用说明

### 上传文章
1. 点击"上传文章"按钮
2. 选择txt、pdf、doc或docx格式的文件
3. 等待上传和分析完成
4. 查看分析结果

### 搜索文章
1. 在文章列表页面的搜索框中输入关键词
2. 支持按标题、内容、作者进行模糊搜索
3. 可与作者筛选功能组合使用

### 查看分析结果
1. 在文章列表中点击文章标题
2. 查看详细的AI分析结果，包括：
   - 文章摘要
   - 关键要点
   - 情感分析
   - 分类标签

## API接口

后端提供RESTful API，主要接口包括：

- `GET /api/v1/articles` - 获取文章列表（支持分页和搜索）
- `GET /api/v1/articles/:id` - 获取文章详情
- `POST /api/v1/articles/upload` - 上传文章文件

## 配置说明

### 后端配置
编辑 `backend/config.yaml` 文件：
```yaml
server:
  port: 8080
  mode: debug

database:
  driver: sqlite
  source: data/article_analysis.db

ai:
  api_key: your-api-key
  api_url: your-api-url
```

### 前端配置
编辑 `frontend/.env` 文件（如需要）：
```
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

## 开发说明

### 后端开发
```bash
cd backend
go mod tidy                    # 更新依赖
go test ./...                  # 运行测试
go run cmd/main.go            # 启动服务
```

### 前端开发
```bash
cd frontend
npm install                    # 安装依赖
npm run dev                    # 开发模式
npm run build                  # 构建生产版本
npm run lint                   # 代码检查
```

## 部署说明

### 生产环境部署
1. 构建前端：
```bash
cd frontend
npm run build
```

2. 构建后端：
```bash
cd backend
go build -o app cmd/main.go
```

3. 配置环境变量和配置文件
4. 启动服务

## 注意事项

- 首次运行时会自动创建数据库表结构
- 上传的文件会保存在 `backend/web/uploads` 目录
- 建议定期备份数据库文件
- 生产环境请修改默认配置和密钥

## 许可证

MIT License

## 技术支持

如有问题，请在项目仓库中提交Issue。