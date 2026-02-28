<p align="center">
  <img src="./doc/assets/logo.png" alt="NexusID Logo" width="200" />
</p>

# NexusID

*[English Version](README.md)* ｜ *[📖 新手使用指引](doc/user_guide.md)* ｜ *[🔧 开发者接入指南](doc/developer_guide.md)*

NexusID 是一个基于 Go（后端）和 React（前端）构建的全栈单点登录 (SSO) 平台。它使用 OpenID Connect (OIDC) 协议提供统一的身份验证和访问管理功能。

## 功能特性

### 核心能力
- **国际化 (i18n)**: 完全支持中英文双语无缝切换
- **OIDC 协议**: 完整的 OpenID Connect (OAuth 2.0 +) 实现
- **多租户架构**: 租户级别的数据完全隔离
- **RS256 JWT 令牌**: 非对称加密机制及提供 JWKS 端点
- **刷新令牌轮换**: 安全的刷新令牌轮换和过期机制
- **用户管理**: 包含注册、登录、角色和权限分配
- **防暴力破解**: 登录限流策略和账户自动锁定机制
- **单点登出**: 全局会话注销及 Token 吊销支持

### 技术栈

**后端：**
- Go 1.23+ 结合 Gin Web 框架
- MySQL 8.0（为提升性能不使用外键）
- Redis 7.0 结合 Singleflight 缓存封装
- GORM 用于数据库对象关系映射
- Zap 用于结构化日志记录
- Viper 用于配置管理
- RS256 JWT 令牌

**前端：**
- React 18 与 TypeScript
- Vite 构建工具
- Tailwind CSS 用于样式渲染
- shadcn/ui 组件库
- React Router 用于页面路由
- react-i18next 用于完整的国际化支持

## 项目结构

```
NexusID/
├── backend/                 # Go 后端服务
│   ├── internal/
│   │   ├── cache/          # 结合 Singleflight 的 Redis 缓存层
│   │   ├── config/         # 配置文件映射管理
│   │   ├── database/       # 数据库连接
│   │   ├── dto/            # 数据传输对象
│   │   ├── handler/        # HTTP 路由控制器
│   │   ├── jwt/            # JWT 生成及密钥管理
│   │   ├── logger/         # Zap 日志封装
│   │   ├── middleware/     # Gin 中间件（包含限流与认证）
│   │   ├── models/         # GORM 实体模型
│   │   ├── repository/     # 数据库仓储层
│   │   └── service/        # 核心业务逻辑层
│   ├── keys/               # RSA 密钥对存储目录
│   ├── main.go             # 应用程序入口
│   └── config.yaml         # 服务器配置文件
├── frontend/               # React 前端
│   ├── src/
│   │   ├── components/     # React 基础组件（包含 shadcn/ui）
│   │   ├── i18n/           # 国际化配置及语言包
│   │   ├── pages/          # 视图级页面组件
│   │   ├── layouts/        # 通用布局框架
│   │   └── lib/            # 工具函数与配置
│   └── package.json
├── deployments/
│   ├── migrations/         # golang-migrate 迁移脚本
│   └── mysql/init/         # MySQL 数据库初始化脚本
└── docker-compose.yml      # 本地环境一键启动配置
```

## 快速入门

### 环境准备
- Go 1.23+
- Node.js 22+
- MySQL 8.0+
- Redis 7.0+
- Docker & Docker Compose (可选)

### 使用 Docker Compose (推荐)

1. 克隆并进入项目目录:
```bash
cd NexusID
```

2. 启动各项服务:
```bash
docker-compose up -d
```

这将启动以下服务:
- MySQL 运行在 3306 端口
- Redis 运行在 6379 端口
- Backend API 运行在 8080 端口

3. 执行数据库迁移:
```bash
# 安装 golang-migrate
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 执行迁移脚本
migrate -path deployments/migrations -database "mysql://nexus_user:nexus_password@tcp(localhost:3306)/nexus_id" up
```

4. 启动前端:
```bash
cd frontend
npm install
npm run dev
```

前端服务将在 `http://localhost:5173` 启动。

5. 使用默认管理员账号登录:

| 字段   | 值                    |
|--------|----------------------|
| 邮箱   | `admin@nexusid.com`  |
| 密码   | `admin12345678`      |

> ⚠️ **重要提示**：在生产环境中，请务必在首次登录后立即修改默认密码。

### 手动部署指南

#### 后端

1. 安装依赖:
```bash
cd backend
go mod download
```

2. 配置环境变量 (修改 `config.yaml` 或通过系统环境变量设定):
```yaml
server:
  port: 8080
  mode: debug

database:
  host: localhost
  port: 3306
  user: nexus_user
  password: nexus_password
  dbname: nexus_id

redis:
  host: localhost
  port: 6379

jwt:
  issuer: http://localhost:8080
  private_key_path: ./keys/private.pem
  public_key_path: ./keys/public.pem
```

3. 启动服务端:
```bash
go run main.go
```

后端 API 将可以通过 `http://localhost:8080` 访问。

#### 前端

1. 安装包依赖:
```bash
cd frontend
npm install
```

2. 启动开发服务器:
```bash
npm run dev
```

前端界面将可以通过 `http://localhost:5173` 访问。

## API 文档摘要

### 健康检查
- `GET /health` - 存活检查
- `GET /ready` - 就绪探测

### OIDC 核心端点
- `GET /.well-known/openid-configuration` - OIDC 服务发现
- `GET /.well-known/jwks.json` - JWKS 公钥分发
- `GET /oauth/authorize` - 授权确认端点
- `POST /oauth/token` - Token 交换/刷新端点
- `POST /oauth/revoke` - Token 吊销端点
- `GET /oauth/logout` - 单点登出端点

### 身份认证
- `POST /api/v1/auth/login` - 账户登录
- `POST /api/v1/auth/register` - 账户注册
- `POST /api/v1/auth/refresh` - 手动刷新访问令牌

### 其他管理 API
*（详情查阅后续接口文档，支持租户、用户、角色及 OIDC 客户端的完整 CRUD 操作）*

## 生产部署建议

### 构建产物

**后端二进制:**
```bash
cd backend
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nexus-id .
```

**前端静态文件:**
```bash
cd frontend
npm run build
```

### 环境注意事项

生产环境清单：
- 将 `GIN_MODE` 设定为 `release`
- 启用 HTTPS/TLS
- 更新并使用强数据库密码
- 修改 `config.yaml` 关闭跨域 (CORS) 允许所有来源的设定
- 配置 RS256 安全密钥的正确挂载
- 建立数据库定时备份机制

## 许可协议

详情请参阅 LICENSE 文件。
