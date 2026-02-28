## Context

随着多套业务系统（Client Applications）的建立，现有的分散式身份验证导致了重复开发、数据孤岛以及安全隐患等问题。我们计划提取认证与授权的公共部分，构建 NexusID（SSO单点登录平台）。该平台作为集中式 Identity Provider (IdP)，专注 Authentication 和部分宏观角色（Role/Group）的下发，而不介入各业务系统的细粒度鉴权（Resource/Action）。技术中台后端采用 Go 语言实现，前端基于 React + Tailwind CSS，底层存储及并发抗压采用 MySQL 8.0 主数据库与 Redis 7.0 高可用缓存。

## Goals / Non-Goals

**Goals:**
- 实现基于 OpenID Connect (OIDC) 授权码模式和 JWT 令牌机制的统一登录服务。
- 使用 RS256 非对称加密算法并暴露 JWKS 端点进行安全的无状态 Token 验证。
- 实现防缓存穿透/击穿（基于 Singleflight 机制）的高并发可靠防线。
- 提供多租户架构，严格隔离不同实体维度的用户信息与客户端配置。
- 提供前后端分离的现代化 UI，包含统一登录注册页、用户中心及超级管理后台。

**Non-Goals:**
- 不接管子系统内部的细粒度操作级权限控制。下游子系统获取 Token 解码获得宏观 Role 后需自行实现内部的业务级鉴权。
- 坚决杜绝依赖基于对称加密（如 HS256）的 Token 签名模式，这在多子系统协同中是不安全的。

## Decisions

1. **认证协议选择：OpenID Connect (OIDC)**
   - *Rationale*: 作为行业标准的基于 OAuth2.0 之上的身份认证层，能将标准用户信息（ID Token）与授权访问区分开来，支持极其广泛的多语言生态安全接入。
   - *Alternatives Considered*: ① 纯 OAuth2.0（缺少原生的标准身份信息封装规范）；② SAML 2.0（基于 XML，过于沉重，与现代云原生和移动端互通不友好）。

2. **签名算法与密钥下发：RS256 非对称加密与 JWKS**
   - *Rationale*: SSO 作为唯一的 Auth Core，持私钥生成签名；所有外部子系统只需通过公开的 `/.well-known/jwks.json` 端点定期拉取公钥进行本地极速验签。这避免了对称加密模式中共享 Secret 密钥一旦泄漏导致的核弹级安全危机。
   - *Alternatives Considered*: HS256 对称加密（需向所有子系统分发同一串 Secret，维护难度大且极其危险）。

3. **数据库存储架构：去外键设计**
   - *Rationale*: 提高单库读写分离时的并发性能，避免高频操作引发的级联复杂锁表；方便后期应对数据暴增时的水平 Sharding（分库分表）；在 Go 的 GORM 业务代码层面上约束保障一致性。
   - *Alternatives Considered*: 遵循传统三范式的严格外键约束设计（不利于高并发场景和微服务演进）。

4. **缓存高可用防线：Singleflight 与 防穿透 NULL 占位**
   - *Rationale*: 引入 `golang.org/x/sync/singleflight` 和为不存在的记录短暂缓存字面量 `NULL`。合并防抖和应对恶意伪造 ID 的突发穿透，使得绝大部分高频恶意请求止步于 Redis，防止数据库被瞬间击穿。

## Risks / Trade-offs

- **[Risk] 非对称加密(RS256)签名带来的 CPU 开销增加**
  - *Mitigation*: 考虑到 Go 核心服务优异的并发性能及协程表现，该开销完全可接受。子系统在验证层面获取公钥后，由于验签也是纯本地无 I/O 的 CPU 运算，对整体延迟影响甚微。
- **[Risk] Redis 宕机导致并发击挂 DB 以及 SSO 会话雪崩**
  - *Mitigation*: 生产环境的 Redis 必须作为一等公民投入维护（推荐集群/哨兵模式加载持久化）。SSO 服务端配置合理的熔断机制，若 Redis 长时间不响应，认证发行新票据的操作可主动限流阻断，保障核心库 MySQL 存活。
