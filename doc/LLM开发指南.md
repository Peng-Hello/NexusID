---
tags:
  - llm
  - guide
  - sso
created: 2026-02-25
---

# 🤖 LLM 开发指南 (LLM Friendly Development Guide)

> **给AI助手的系统提示 (System Prompt for AI Assistants)**
> 
> 当你(AI)被要求协助开发本 SSO(单点登录) 项目时，**必须**遵循本文档中定义的架构约束和设计哲学。本指南是项目中最核心的高优先级规则，其重要性高于通常的最佳实践。

## 1. 核心架构约束 (Core Architecture Constraints)

### 1.1 身份管理 vs 细粒度权限
这个 SSO 平台的核心定位是 **IdP (Identity Provider)**。
*   **绝对禁止 (MUST NOT)**: 在 SSO 数据库中建立复杂的 `permissions`, `role_permissions`，或处理类似于“某个用户能否点击商品发布按钮”的业务逻辑。
*   **必须遵守 (MUST DO)**: SSO 仅维护用户的宏观身份（Role / Group）。具体的 Resource 和 Action 鉴权校验必须 **下放给接入的业务子系统自身去实现**。

### 1.2 JWT 与安全体系
*   **绝对禁止 (MUST NOT)**: 使用对称加密 (如 `HS256`) 签发 JWT Access/ID Token！这会导致边缘子系统泄露 Secret 时全盘崩溃。
*   **必须遵守 (MUST DO)**: 
    *   必须使用 **非对称加密 (如 RS256 / ES256)**。
    *   必须提供 `/.well-known/jwks.json` 发现端点，供子系统拉取公钥进行无秘钥的本地验签。
    *   必须在生成 Token 和刷新 Token 的相关代码中支持 **Refresh Token Rotation (刷新令牌轮换)**，检测到旧刷新令牌被滥用时，必须使关联的整个令牌树失效。

### 1.3 数据库设计反模式红线 (Database Anti-Patterns)
*   **绝对禁止 (MUST NOT)**: 在 MySQL 表结构设计中使用物理外键（`FOREIGN KEY ... ON DELETE CASCADE`）。
*   **必须遵守 (MUST DO)**: 
    *   依赖业务代码（Service 层）通过代码逻辑或最终一致性任务来维护数据隔离和清理。
    *   所有核心表（`users`, `roles`, `clients`, `audit_logs`）必须预留 `tenant_id` 字段以支持 SaaS 化。默认全局租户的 ID 为 `0`。
    *   软删除必须使用 `deleted_at` (时间戳) 而不是简单的 `status`。必须在唯一约束索引中加入 `deleted_at` 字段，以解决同名用户被逻辑删除后，新用户无法注册的冲突。

### 1.4 高并发缓存保护 (Cache Protection)
*   **绝对禁止 (MUST NOT)**: 直接执行 `Get(Redis) -> if empty -> Get(MySQL) -> Set(Redis)` 这种脆弱的逻辑，它在高并发下必定引起缓存击穿。
*   **必须遵守 (MUST DO)**: 
    *   在 Go 代码中查询数据库前，必须使用 `golang.org/x/sync/singleflight` 把瞬间的高并发请求合并为一个。
    *   查询不到的记录也必须写入 Redis（设置短时间 TTL 占位符如 `NULL`），有效阻断针对不存在 ID 的恶意穿透攻击。

## 2. 技术栈规定 (Tech Stack Specifications)

在生成代码时，优先使用以下技术栈：

### 2.1 后端 (Backend)
*   **语言**: Go `1.21+`
*   **Web 框架**: `gin-gonic/gin`
*   **ORM**: `gorm.io/gorm` (注意：GORM 的表关联不要自动生成外键约束，使用 `gorm:"-"` 忽略物理外键)
*   **缓存**: `redis/go-redis/v9`  
*   **配置**: `spf13/viper`
*   **日志**: `go.uber.org/zap`
*   **路由/安全**: 强制要求基于 JWT RS256 进行鉴权拦截。

### 2.2 前端 (Frontend)
*   **仓库结构**: 本项目是一个 **全栈 Monorepo**。后端代码在 `/backend` 目录下，前端代码在 `/frontend` 目录下。AI 必须将对应语言的代码生成到指定的模块目录中。
*   **UI 框架**: React (如果在规划阶段使用了 Next.js，需明确是 App Router 还是 Pages Router，但标准 SPA 使用 Vite 构建器也可)。
*   **样式方案**: Tailwind CSS。
*   **组件库**: **[shadcn/ui](https://ui.shadcn.com/)**。任何前端组件优先思考使用 `shadcn/ui` 风格的基础组件（如 Card, Button, Input, Table 等）进行组合开发。
*   **API 请求**: `axios`。前端在收到 Token 后，必须建立专门的 Token Service（管理在 HttpOnly Cookie 或安全的 local storage 中），并配置统一的 Axios Interceptor 来全局处理 401 响应并无缝调用 `/api/auth/refresh` 进行令牌轮换续期。

## 3. 代码生成示例风格 (Code Generation Expected Style)

当被要求生成 `User` 模型时，**你生成的代码必须长这样**：

```go
type User struct {
    ID            uint64         `gorm:"primarykey;autoIncrement"`
    TenantID      uint64         `gorm:"not null;default:0;uniqueIndex:uk_tenant_username_deleted;uniqueIndex:uk_tenant_email_deleted"`
    Username      string         `gorm:"type:varchar(50);not null;uniqueIndex:uk_tenant_username_deleted"`
    Email         string         `gorm:"type:varchar(100);not null;uniqueIndex:uk_tenant_email_deleted"`
    PasswordHash  string         `gorm:"type:varchar(255);not null"`
    Status        int8           `gorm:"type:tinyint;default:1"` // 1: Active, 0: Disabled
    DeletedAt     uint64         `gorm:"default:0;uniqueIndex:uk_tenant_username_deleted;uniqueIndex:uk_tenant_email_deleted"` // 0 means active
    CreatedAt     time.Time      `gorm:"autoCreateTime"`
    UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
}
```
*注意到了吗？没有 `gorm.DeletedAt` 结构体，手动管理 Unix 时间戳 `DeletedAt` 以配合混合唯一索引。没有任何物理外键声明。*

## 4. 目录流转与协作方式 (Workflow)
1. **优先提取上下文**: 处理任何需求前，请首先 grep/阅读本项目的 `数据库设计.md` 以及 `技术架构设计.md` 了解全局约束。
2. **渐进式实施**: 提供代码前，必须进行思考并在必要时制定短期的 implementation plan。
3. **安全自检**: 部署代码或结束对话前，自行走查是否存在上文提及的“反模式红线”。
