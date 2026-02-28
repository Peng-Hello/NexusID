# NexusID 开发者接入指南

> 本文档面向**开发者**，介绍如何将你的应用系统接入 NexusID，实现用户统一登录（单点登录 SSO）。

---

## 接入前提

在开始之前，你需要联系 NexusID 管理员完成以下注册：

1. **创建租户**：让管理员为你的组织创建一个租户。
2. **注册客户端**：让管理员在该租户下注册你的应用系统，你会拿到：
   - `client_id`：客户端标识，公开的
   - `client_secret`：客户端密钥，**请妥善保管，不要泄露**

---

## 核心概念：OIDC 登录流程

NexusID 基于标准 **OIDC（OpenID Connect）** 协议。接入后，登录流程如下：

```
1. 用户点击你系统里的「登录」按钮
2. 你的系统把用户重定向到 NexusID 的授权页
3. 用户在 NexusID 输入账号密码完成认证
4. NexusID 把用户重定向回你系统的「回调地址」，并携带一个临时 code
5. 你的后端用这个 code 换取 access_token 和 id_token
6. 用 access_token 调用 NexusID 的 userinfo API 获取用户信息
7. 登录成功，建立你系统的会话
```

---

## 第一步：构造登录跳转 URL

当用户点击「登录」时，把用户重定向到以下地址：

```
GET http://{NexusID地址}/api/v1/oauth/authorize
  ?response_type=code
  &client_id={你的client_id}
  &redirect_uri={你的回调地址}
  &scope=openid profile email
  &state={随机字符串，用来防 CSRF}
```

**参数说明：**

| 参数 | 说明 | 示例 |
|------|------|------|
| `response_type` | 固定填 `code` | `code` |
| `client_id` | 你的客户端 ID | `abc123` |
| `redirect_uri` | 登录完跳回你系统的地址（必须和注册时填写的一致） | `https://yourapp.com/callback` |
| `scope` | 申请的权限范围 | `openid profile email` |
| `state` | 随机字符串，你自己生成，登录完会原样返回，用于验证防伪 | `xK9f2mZ1` |

**示例：**
```
https://nexusid.example.com/api/v1/oauth/authorize?response_type=code&client_id=abc123&redirect_uri=https://yourapp.com/callback&scope=openid%20profile%20email&state=xK9f2mZ1
```

---

## 第二步：处理回调，用 code 换 token

用户登录成功后，NexusID 会把用户重定向到你的回调地址，并附上参数：

```
GET https://yourapp.com/callback?code=AUTH_CODE&state=xK9f2mZ1
```

**你的后端需要做：**

1. 验证 `state` 是否和你之前发出去的一致（防 CSRF）
2. 用 `code` 换取 token：

```http
POST http://{NexusID地址}/api/v1/oauth/token
Content-Type: application/x-www-form-urlencoded

grant_type=authorization_code
&code={上一步拿到的code}
&redirect_uri={和上一步一样的回调地址}
&client_id={你的client_id}
&client_secret={你的client_secret}
```

**成功响应：**
```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4...",
  "id_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

## 第三步：获取用户信息

用 `access_token` 调用 userinfo 端点：

```http
GET http://{NexusID地址}/api/v1/oauth/userinfo
Authorization: Bearer {access_token}
```

**成功响应：**
```json
{
  "sub": "123",
  "email": "zhangsan@yourcompany.com",
  "name": "张三",
  "email_verified": true
}
```

拿到这些信息后，在你自己的系统里创建/更新用户并建立登录会话，接入完成！

---

## 第四步：刷新 token（可选）

`access_token` 默认 1 小时过期，过期前可以用 `refresh_token` 换新的：

```http
POST http://{NexusID地址}/api/v1/oauth/token
Content-Type: application/x-www-form-urlencoded

grant_type=refresh_token
&refresh_token={你的refresh_token}
&client_id={你的client_id}
&client_secret={你的client_secret}
```

---

## 第五步：退出登录

用户退出时，调用两个接口：

**① 吊销 refresh_token（防止被继续使用）：**
```http
POST http://{NexusID地址}/api/v1/oauth/revoke
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "token": "{refresh_token}"
}
```

**② 清除 NexusID 会话（可选，用于单点登出）：**
```http
GET http://{NexusID地址}/api/v1/oauth/logout
  ?post_logout_redirect_uri=https://yourapp.com/logged-out
```

---

## API 地址汇总

| 用途 | 方法 | 地址 |
|------|------|------|
| 发起授权 | `GET` | `/api/v1/oauth/authorize` |
| 换取 token | `POST` | `/api/v1/oauth/token` |
| 获取用户信息 | `GET` | `/api/v1/oauth/userinfo` |
| 吊销 token | `POST` | `/api/v1/oauth/revoke` |
| 退出登录 | `GET` | `/api/v1/oauth/logout` |
| 获取公钥（验签用）| `GET` | `/api/v1/oauth/jwks` |

---

## 验证 id_token 的方式（进阶）

`id_token` 是一个 JWT，你可以本地验证它的合法性，不需要每次都请求 NexusID：

1. 从 `/api/v1/oauth/jwks` 获取公钥
2. 用公钥验证 `id_token` 的签名（RS256 算法）
3. 检查 `iss`（颁发者）、`aud`（受众=你的client_id）、`exp`（过期时间）

推荐使用各语言现成的 JWT 库：

| 语言 | 推荐库 |
|------|--------|
| Go | `github.com/golang-jwt/jwt` |
| Python | `python-jose` |
| Node.js | `jose` |
| Java | `nimbus-jose-jwt` |

---

## 常见问题

**Q：redirect_uri 不匹配怎么办？**  
A：你的回调地址必须与注册客户端时填写的完全一致（包括协议、端口、路径），否则 NexusID 会拒绝请求。让管理员在客户端配置里确认或修改。

**Q：客户端密钥泄露了怎么办？**  
A：立刻让管理员在 NexusID 后台点击该客户端的「轮换密钥」，旧密钥会立即失效，拿到新密钥后更新你的应用配置。

**Q：公开客户端（SPA/移动端）怎么接入？**  
A：注册时勾选「公开客户端」，换 token 时不需要传 `client_secret`，改用 PKCE 流程（后续版本支持）。

**Q：同一个用户在多个应用间如何实现无感 SSO？**  
A：只要用户已经在 NexusID 登录过（Session 未过期），在访问其他接入系统时，步骤 3（用户输入密码）会被自动跳过，直接返回 code，达到无感单点登录效果。

---

*如有接入问题，请联系 NexusID 管理员。*
