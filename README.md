# ResumeAI - 智能简历管理系统

一个基于 Go 和 Vue.js 开发的现代化简历管理与面试追踪系统，支持简历文件管理、面试记录追踪、岗位管理、招聘来源管理和面试官信息管理。

## ✨ 功能特性

### 📄 简历库管理
- **文件上传**：支持 PDF、DOC、DOCX 格式简历文件
- **批量操作**：支持多文件批量上传和删除
- **文件预览**：在线预览和下载简历文件
- **文件筛选**：按文件类型（PDF/Word）筛选
- **文件搜索**：按文件名关键字搜索
- **文件统计**：显示文件总数和总存储占用

### 🎯 面试管理
- **面试记录创建**：记录候选人基本信息、面试信息
- **简历关联**：上传简历文件并与面试记录关联
- **面试状态追踪**：通过/待定/淘汰三种状态
- **面试评分**：0-10分评分系统
- **多维度筛选**：按面试建议（通过/待定/淘汰）筛选
- **关键字搜索**：按姓名、岗位、面试官等搜索
- **面试统计**：总面试数、待定数、通过数统计

### 💼 岗位管理
- **岗位创建/编辑/删除**：维护公司招聘岗位列表
- **岗位搜索**：按岗位名称搜索
- **岗位统计**：总岗位数、今日新增统计

### 📡 招聘来源管理
- **来源创建/编辑/删除**：维护招聘渠道信息
- **来源搜索**：按来源名称搜索
- **来源统计**：总来源数、今日新增统计

### 👥 面试官管理
- **面试官信息管理**：姓名、邮箱、部门
- **邮箱验证**：自动验证邮箱格式
- **搜索功能**：按姓名、邮箱、部门搜索
- **统计信息**：总人数、有效邮箱数、有部门信息人数

## 🛠️ 技术栈

### 后端
- **语言**：Go 1.16+
- **框架**：Gin Web Framework
- **ORM**：GORM
- **数据库**：MySQL
- **文件存储**：本地文件系统
- **配置管理**：godotenv

### 前端
- **框架**：Vue.js 2.7
- **UI组件库**：Element UI
- **样式**：自定义 CSS + Font Awesome
- **HTTP 请求**：Fetch API

## 📋 系统要求

- Go 1.16 或更高版本
- MySQL 5.7 或更高版本
- 现代浏览器（Chrome, Firefox, Safari, Edge）

## 🚀 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/yourusername/Resume.git
cd Resume
```

### 2. 配置文件

在项目根目录创建 `.env` 文件：

```env
# 服务器配置
HOST=0.0.0.0
PORT=8080
DEBUG=true

# 存储配置
STORAGE_PATH=../storage/data
BASE_URL=http://0.0.0.0:8080

# 文件限制
MAX_FILE_SIZE=52428800  # 50MB

# 数据库配置 (MySQL)
DB_DRIVER=mysql
DB_HOST=192.168.222.170
DB_PORT=3306
DB_USER=root
DB_PASSWORD=123456
DB_NAME=Resume
DB_CHARSET=utf8mb4
DB_PARSE_TIME=true
DB_LOC=Local
```

### 3. 创建数据库

```sql
CREATE DATABASE resume_ai CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. 安装依赖

```bash
# 下载 Go 依赖
go mod download
```

### 5. 运行服务

```bash
go run main.go
```

服务启动后，访问 `http://localhost:8080` 即可进入系统。

## 📁 项目结构

```
Resume/
├── main.go                     # 主程序入口
├── go.mod                      # Go模块文件
├── .env                        # 环境配置文件
├── backend/
│   ├── config/                 # 配置管理
│   │   └── config.go          # 配置加载
│   ├── database/               # 数据库连接
│   │   └── database.go         # GORM 初始化
│   ├── handlers/               # 请求处理器
│   │   ├── upload.go           # 文件上传处理
│   │   ├── interview.go        # 面试管理
│   │   ├── company_position.go # 岗位管理
│   │   ├── recruitment_source.go # 招聘来源管理
│   │   └── employee.go         # 面试官管理
│   ├── middleware/             # 中间件
│   │   └── cors.go             # CORS 跨域处理
│   ├── models/                 # 数据模型
│   │   ├── file.go             # 文件模型
│   │   ├── interview.go        # 面试模型
│   │   ├── company_position.go # 岗位模型
│   │   ├── recruitment_source.go # 招聘来源模型
│   │   ├── employee.go         # 面试官模型
│   │   └── responses.go        # 响应结构
│   ├── routers/                # 路由配置
│   │   └── router.go           # API 路由
│   └── services/               # 业务逻辑
│       ├── storage.go          # 文件存储服务
│       ├── interview_service.go # 面试服务
│       ├── company_position_service.go # 岗位服务
│       ├── recruitment_source_service.go # 招聘来源服务
│       └── employee_service.go # 面试官服务
├── frontend/                   # 前端静态文件
│   ├── static/
│   │   ├── css/                # 样式文件
│   │   ├── js/                 # JavaScript 文件
│   │   │   ├── vue.min.js      # Vue.js 核心
│   │   │   ├── main.js         # 简历库页面逻辑
│   │   │   ├── interview-form.js # 面试管理逻辑
│   │   │   ├── company-positions.js # 岗位管理逻辑
│   │   │   ├── recruitment-sources.js # 招聘来源逻辑
│   │   │   └── interviewers.js # 面试官管理逻辑
│   │   └── fonts/              # 字体文件
│   └── templates/              # HTML 模板
│       ├── index.html          # 简历库页面
│       ├── interview-form.html # 面试管理页面
│       ├── company-positions.html # 岗位管理页面
│       ├── recruitment-sources.html # 招聘来源页面
│       └── interviewers.html   # 面试官管理页面
└── storage/                    # 文件存储目录
    └── data/                   # 按年月组织的文件
```

## 🔧 API 接口

### 文件管理

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/upload` | 单文件上传 |
| POST | `/api/v1/upload/batch` | 批量文件上传 |
| GET | `/api/v1/files` | 获取文件列表 |
| GET | `/api/v1/files/:id` | 获取文件信息 |
| DELETE | `/api/v1/files/:id` | 删除文件 |
| DELETE | `/api/v1/files/batch` | 批量删除文件 |
| GET | `/api/v1/files/:id/download` | 下载文件 |
| GET | `/api/v1/files/search` | 搜索文件 |
| GET | `/api/v1/stats` | 获取统计信息 |

### 面试管理

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/interviews` | 创建面试记录 |
| GET | `/api/v1/interviews` | 获取面试列表 |
| GET | `/api/v1/interviews/:id` | 获取面试详情 |
| PUT | `/api/v1/interviews/:id` | 更新面试记录 |
| DELETE | `/api/v1/interviews/:id` | 删除面试记录 |
| POST | `/api/v1/interviews/upload-resume` | 上传简历文件 |

### 岗位管理

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/company-positions` | 创建岗位 |
| GET | `/api/v1/company-positions` | 获取岗位列表 |
| GET | `/api/v1/company-positions/:id` | 获取岗位详情 |
| PUT | `/api/v1/company-positions/:id` | 更新岗位 |
| DELETE | `/api/v1/company-positions/:id` | 删除岗位 |

### 招聘来源管理

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/recruitment-sources` | 创建招聘来源 |
| GET | `/api/v1/recruitment-sources` | 获取来源列表 |
| GET | `/api/v1/recruitment-sources/:id` | 获取来源详情 |
| PUT | `/api/v1/recruitment-sources/:id` | 更新来源 |
| DELETE | `/api/v1/recruitment-sources/:id` | 删除来源 |

### 面试官管理

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/employees` | 创建面试官 |
| GET | `/api/v1/employees` | 获取面试官列表 |
| GET | `/api/v1/employees/:id` | 获取面试官详情 |
| PUT | `/api/v1/employees/:id` | 更新面试官 |
| DELETE | `/api/v1/employees/:id` | 删除面试官 |

## 💡 使用说明

### 简历库

1. **上传文件**：点击"上传文件"按钮或拖拽文件到指定区域
2. **筛选文件**：使用顶部分类标签按 PDF/Word 筛选
3. **搜索文件**：在搜索框中输入文件名关键字
4. **查看详情**：点击文件行的"详情"按钮查看文件信息
5. **下载文件**：点击"下载"按钮获取文件
6. **删除文件**：点击"删除"按钮，支持单个或批量删除

### 面试管理

1. **创建面试记录**：点击"新增面试"按钮
2. **填写信息**：
   - 候选人基本信息（姓名、手机号、邮箱）
   - 选择应聘岗位
   - 上传或选择简历文件
   - 选择招聘来源
   - 设置面试轮次、形式、时间
   - 选择面试官
   - 填写评分、评价和建议
3. **查看记录**：点击"详情"查看完整面试信息
4. **编辑记录**：点击"编辑"修改面试信息
5. **删除记录**：点击"删除"移除记录

### 岗位管理

1. **新增岗位**：点击"新增岗位"按钮，输入岗位名称
2. **编辑岗位**：点击编辑图标修改岗位名称
3. **删除岗位**：点击删除图标移除岗位

### 招聘来源管理

1. **新增来源**：点击"新增来源"按钮，输入来源名称
2. **编辑来源**：点击编辑图标修改来源名称
3. **删除来源**：点击删除图标移除来源

### 面试官管理

1. **新增面试官**：点击"新增面试官"按钮
2. **填写信息**：姓名（必填）、邮箱（必填）、部门（可选）
3. **编辑面试官**：点击编辑图标修改信息
4. **删除面试官**：点击删除图标移除面试官

## 🔒 数据验证

- **简历文件唯一性**：同一简历文件不能被多个面试记录使用
- **邮箱格式验证**：面试官和候选人邮箱需符合格式要求
- **评分范围验证**：面试评分限制在 0-10 分
- **必填项验证**：关键字段不能为空

## 📊 数据库表结构

| 表名 | 描述 |
|------|------|
| `file_infos` | 文件信息（ID、名称、路径、大小等）|
| `interviews` | 面试记录（候选人信息、面试详情等）|
| `company_positions` | 岗位信息 |
| `recruitment_sources` | 招聘来源信息 |
| `employees` | 面试官信息 |

## 🐛 常见问题

### 1. 数据库连接失败？
- 检查 MySQL 服务是否启动
- 确认 `.env` 中的数据库配置正确
- 验证数据库用户权限

### 2. 文件上传失败？
- 检查存储目录权限
- 确认文件大小不超过配置限制（默认 50MB）
- 验证文件格式（仅支持 PDF、DOC、DOCX）

### 3. 简历文件无法关联？
- 确保简历文件已上传到系统
- 检查简历是否已被其他面试记录使用
---

**注意**：生产环境使用时，请确保：
- 配置安全的数据库密码
- 设置合理的文件上传大小限制
- 定期备份数据库和存储文件
- 使用 HTTPS 协议部署
