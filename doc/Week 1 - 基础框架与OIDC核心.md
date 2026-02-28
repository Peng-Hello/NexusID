---
tags:
  - tasklist
  - sso
  - week1
created: 2026-02-25
week: 1
timeline: 2026-02-25 --- 2026-03-03
---

# Week 1 - 基础框架与OIDC核心

## 📅 时间安排
2026-02-25 --- 2026-03-03 (7天)

---

## Day 1-2: 项目初始化 + 基础设施

### 任务清单
- [ ] 项目脚手架搭建
- [ ] 依赖管理配置 (go.mod)
- [ ] 配置管理实现 (Viper)
- [ ] 数据库连接与ORM配置 (GORM)
- [ ] Redis连接配置
- [ ] 日志系统 (zap)
- [ ] 数据库迁移脚本
- [ ] Docker环境配置

### 验收标准
- [x] 项目可以成功启动
- [ ] 数据库连接正常
- [ ] Redis连接正常
- [ ] 日志正常输出

### 预计工时
16小时

---

## Day 3-4: 用户管理模块

### 任务清单
- [ ] 用户领域模型定义
- [ ] 用户Repository实现
- [ ] 用户Service实现
  - [ ] 用户注册 (密码加密、邮箱验证)
  - [ ] 用户登录 (密码验证)
  - [ ] 用户信息查询/更新
  - [ ] 密码重置
- [ ] 用户管理API
  - [ ] POST /api/v1/users/register
  - [ ] POST /api/v1/users/login
  - [ ] GET /api/v1/users/profile
  - [ ] PUT /api/v1/users/profile
  - [ ] POST /api/v1/users/password/reset
- [ ] 密码加密工具 (bcrypt)
- [ ] JWT Token生成/验证工具

### 验收标准
- [ ] 用户可以成功注册
- [ ] 用户可以成功登录并获取token
- [ ] JWT token验证正常工作

### 预计工时
16小时

---

## Day 5-7: OIDC核心实现

### 任务清单
- [ ] 客户端应用管理
  - [ ] 客户端注册/管理
  - [ ] Client ID/Secret生成
- [ ] OIDC授权端点实现
  - [ ] GET /oauth/authorize - 授权页面
  - [ ] POST /oauth/authorize - 用户确认授权
- [ ] OIDC Token端点实现
  - [ ] POST /oauth/token - 获取token (authorization code)
  - [ ] POST /oauth/token/refresh - 刷新token
  - [ ] POST /oauth/token/revoke - 撤销token
- [ ] OIDC用户信息端点
  - [ ] GET /oauth/userinfo - 获取用户信息
- [ ] OIDC配置端点
  - [ ] GET /.well-known/openid-configuration
- [ ] 授权码生成与验证 (支持PKCE)
- [ ] Token生成与验证 (JWT格式)
  - [ ] Access Token
  - [ ] Refresh Token
  - [ ] ID Token
- [ ] Redis缓存授权码 (过期控制)

### 验收标准
- [ ] OIDC授权流程完整可用
- [ ] 可以成功获取token
- [ ] ID Token符合OIDC规范
- [ ] PKCE正常工作

### 预计工时
24小时

---

## 📊 本周总结

### 完成情况
- [ ] Day 1-2: 项目初始化
- [ ] Day 3-4: 用户管理模块
- [ ] Day 5-7: OIDC核心实现

### 遇到的问题
<!-- 记录遇到的问题和解决方案 -->

### 下周计划
- [ ] RBAC权限系统
- [ ] 管理后台
- [ ] 集成测试
