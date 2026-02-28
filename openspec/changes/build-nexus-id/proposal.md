## Why

随着业务系统的增多，各个子系统独立维护复杂的鉴权体系导致开发成本高、用户体验碎片化，且难以推行统一的安全标准。构建 NexusID（SSO单点登录平台）旨在专注身份认证（Authentication / IdP）以及宏观身份下发，解决跨系统登录问题，降低各业务系统的接入与维护成本。

## What Changes

- 基于 Golang 开发高性能核心认证服务，替换/移除传统的分散式鉴权逻辑。
- 实现完整的 OpenID Connect (OIDC) 协议规范（含JWKS发现与单点登出）。
- 建立多租户（MultiTenant）架构，支持隔离的用户、客户端和角色数据管理。
- 引入基于 RS256 非对称加密的 JWT 令牌机制与 Refresh Token Rotation 管理机制。
- 构建基于 React + Tailwind CSS + shadcn/ui 的用户交互和基础身份/角色管理后台。
- 采用 MySQL 8.0（去外键设计）和 Redis 7.0（含 Singleflight 防击穿方案）作为底层高可用存储支撑。

## Capabilities

### New Capabilities

- `oidc-core`: OIDC 协议核心实现，包含授权码流程、Token 签发、JWKS 端点及单点登出支持。
- `tenant-management`: 多租户系统的初始化与生命周期管理。
- `user-identity`: 支持多租户维度的用户CRUD注册登录、基础角色（Groups/Roles）定义与授权分配。
- `session-security`: 分布式环境下的会话维持、单点设备踢出（Blacklist）机制与 Refresh Token 防重放防盗用。

### Modified Capabilities



## Impact

- **业务应用客户端 (Client Applications)**: 需调整鉴权架构，将自身的登录逻辑切换为 OIDC 交互流程获取 Token，并在本地完成基于 JWT 下发角色的资源鉴权。
- **全栈项目结构**: 引入 Full-Stack Monorepo (nexus-id) 来同时管理 Go 后端、React 前端以及 Docker 运维脚本。
- **公共服务资源**: 需要申请并部署独立的 MySQL 和 Redis 实例专门用于 NexusID 核心存储。
