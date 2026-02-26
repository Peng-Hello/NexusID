---
tags:
  - project
  - sso
  - golang
  - oidc
  - active
created: 2026-02-25
status: 🟡 规划中
timeline: 2026-02-25 --- 2026-03-11
---

# NexusID (SSO单点登录平台)

## 📋 项目概述

创建一个遵循云原生架构与行业安全的单点登录（SSO）平台，支持多个业务系统共用登录功能。基于OpenID Connect (OIDC)协议实现，使用Go语言开发，采用MySQL + Redis作为数据存储与高可用缓存方案。

**定位**：专注身份认证（Authentication / IdP）以及宏观身份下发，不再深度捆绑各业务子系统的复杂鉴权体系。

**项目周期**: 2026-02-25 --- 2026-03-11 (2周)

---

## 🎯 核心目标

- [x] 项目架构与数据库严谨规划（修正反模式，拥抱高并发与非对称加密）
- [ ] 完整的OIDC协议实现（含JWKS发现与单点登出）
- [ ] 多租户架构支持的用户注册/登录/管理
- [ ] 基础身份与角色管理后台
- [ ] 分布式环境下的令牌管理 (Refresh Token Rotation)
- [ ] 生产可用的无状态 Docker 部署方案

---

## 🏗️ 技术架构

### 技术栈
- **语言**: Go 1.21+ (后端), TypeScript (前端)
- **后端框架**: Gin (Web), GORM (DB)
- **前端框架**: React, Tailwind CSS, shadcn/ui
- **底层硬件**: MySQL 8.0 (去外键设计) + Redis 7.0 (Singleflight防击穿)
- **安全**: JWT (RS256 非对称加密), bcrypt 密码哈希
- **配置管理**: Viper
- **日志**: zap
- **数据库迁移**: golang-migrate

### 架构图
```
┌─────────────────────────────────────────────────────────────┐
│                    Client Applications                      │
│  (业务系统获取 Token 后，本地自行解析提取角色并完成资源鉴权)    │
└──────────────────────┬──────────────────────────────────────┘
                       │ OIDC Protocol (Authorization Code)
                       │ Pull from /.well-known/jwks.json
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                    SSO Platform (Go)                         │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────────┐   │
│  │ Auth Core   │  │ User Mgmt    │  │ Identity System  │   │
│  │ - OIDC/JWK  │  │ - CRUD       │  │ - Groups/Roles   │   │
│  │ - JWT/RS256 │  │ - MultiTenant│  │ - Assignments    │   │
│  └─────────────┘  └──────────────┘  └──────────────────┘   │
└──────────────────────┬──────────────────────────────────────┘
                       │
        ┌──────────────┴──────────────┐
        ▼                             ▼
┌──────────────┐            ┌────────────────┐
│   MySQL      │            │     Redis      │
│ - users      │            │ - Sessions     │
│ - clients    │            │ - Singleflight │
│ - roles      │            │ - Blacklist    │
└──────────────┘            └────────────────┘
```

---

## 📂 项目结构 (Full-Stack Monorepo)

本项目采用全栈 Monorepo 结构，前端和后端代码位于同一个代码仓库中，方便统一版本管理和部署。

```
nexus-id/
├── backend/              # Go 后端服务
│   ├── cmd/server/main.go
│   ├── internal/         # 核心业务逻辑 (Auth, OIDC, User, Identity)
│   ├── api/openapi.yaml  # API文档
│   ├── migrations/       # 数据库迁移脚本
│   └── configs/          # 后端配置文件
├── frontend/             # React + shadcn/ui 前端应用
│   ├── src/
│   │   ├── app/          # 页面路由 (Login, Admin, Dashboard 等)
│   │   ├── components/   # shadcn/ui 和自定义组件
│   │   ├── lib/          # 工具函数 (API 客户端, Token 管理)
│   │   └── hooks/        # 自定义 React Hooks
│   ├── package.json
│   └── tailwind.config.js
└── docker/               # Docker & 部署配置
```

---

## 🚀 快速开始

### 开发环境要求
- Go 1.21+
- MySQL 8.0+
- Redis 7.0+
- Docker & Docker Compose

### 启动步骤
```bash
# 1. 克隆项目
git clone <repo>
cd nexus-id

# 2. 启动依赖服务
docker-compose up -d mysql redis

# 3. 配置环境变量
cp configs/config.yaml configs/config.local.yaml

# 4. 运行数据库迁移
make migrate

# 5. 启动服务
make run
```

---

## 📊 进度追踪

- Week 1: [[Week 1 - 基础框架与OIDC核心]]
- Week 2: [[Week 2 - 权限系统与管理后台]] (现已转型为基础身份管理)

---

## 📝 相关文档

- [[技术架构设计]]
- [[数据库设计]]
- [[API文档]]
- [[部署文档]]
- [[安全指南]]

---

## 🔗 外部资源

- [OIDC官方规范](https://openid.net/connect/)
- [Gin框架文档](https://gin-gonic.com/docs/)
- [GORM文档](https://gorm.io/docs/)
