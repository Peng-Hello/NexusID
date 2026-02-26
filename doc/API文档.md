---
tags:
  - api
  - documentation
  - sso
  - oidc
created: 2026-02-25
version: v1.0.0
---

# API文档

## 📡 基础信息

**Base URL**: `http://localhost:8080`
**API Version**: v1.0.0
**协议**: HTTP/HTTPS
**数据格式**: JSON

---

## 🔐 认证方式

### Bearer Token
大多数API需要在Header中携带Access Token:

```
Authorization: Bearer <access_token>
```

---

## 🌐 OIDC协议端点

### 1. 发现端点

**GET** `/.well-known/openid-configuration`

获取OIDC Provider的配置信息。

**响应示例**:
```json
{
  "issuer": "http://localhost:8080",
  "authorization_endpoint": "http://localhost:8080/oauth/authorize",
  "token_endpoint": "http://localhost:8080/oauth/token",
  "userinfo_endpoint": "http://localhost:8080/oauth/userinfo",
  "jwks_uri": "http://localhost:8080/.well-known/jwks.json",
  "scopes_supported": ["openid", "profile", "email"],
  "response_types_supported": ["code"],
  "grant_types_supported": ["authorization_code", "refresh_token"],
  "subject_types_supported": ["public"],
  "id_token_signing_alg_values_supported": ["RS256"]
}
```

---

### 2. 授权端点

**GET** `/oauth/authorize`

发起授权请求，用户登录并授权。

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| client_id | string | 是 | 客户端ID |
| redirect_uri | string | 是 | 回调URI |
| response_type | string | 是 | 固定值: code |
| scope | string | 是 | 请求的权限范围 |
| state | string | 是 | CSRF防护 |
| code_challenge | string | 否 | PKCE code challenge |
| code_challenge_method | string | 否 | 固定值: S256 |

**请求示例**:
```
GET /oauth/authorize?client_id=abc123&redirect_uri=http://localhost:3000/callback&response_type=code&scope=openid profile email&state=xyz&code_challenge=E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM&code_challenge_method=S256
```

**响应**:
- 返回登录页面 (HTML)
- 用户登录后返回授权确认页面

---

### 3. Token端点

**POST** `/oauth/token`

用授权码换取Access Token。

**请求参数** (application/x-www-form-urlencoded):
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| grant_type | string | 是 | 固定值: authorization_code |
| code | string | 是 | 授权码 |
| redirect_uri | string | 是 | 回调URI |
| client_id | string | 是 | 客户端ID |
| client_secret | string | 否 | 客户端密钥 (非公开客户端) |
| code_verifier | string | 否 | PKCE code verifier |

**请求示例**:
```bash
curl -X POST http://localhost:8080/oauth/token \
  -d "grant_type=authorization_code" \
  -d "code=SplxlOBeZQQYbYS6WxSbIA" \
  -d "redirect_uri=http://localhost:3000/callback" \
  -d "client_id=abc123" \
  -d "client_secret=secret"
```

**响应示例**:
```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "tGzv3JOkF0XG5Qx2TlKWIA",
  "id_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

### 4. 刷新Token端点

**POST** `/oauth/token/refresh`

使用Refresh Token获取新的Access Token。

**请求参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| grant_type | string | 是 | 固定值: refresh_token |
| refresh_token | string | 是 | 刷新令牌 |
| client_id | string | 是 | 客户端ID |
| client_secret | string | 否 | 客户端密钥 |

**响应示例**:
```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "tGzv3JOkF0XG5Qx2TlKWIA"
}
```

---

### 5. 用户信息端点

**GET** `/oauth/userinfo`

获取当前用户信息。

**请求头**:
```
Authorization: Bearer <access_token>
```

**响应示例**:
```json
{
  "sub": "1234567890",
  "username": "john_doe",
  "email": "john@example.com",
  "email_verified": true,
  "name": "John Doe",
  "picture": "http://localhost:8080/avatars/john.jpg",
  "roles": ["user"]
}
```

---

### 6. 撤销Token端点

**POST** `/oauth/revoke`

撤销Access Token或Refresh Token。

**请求参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| token | string | 是 | 要撤销的token |
| token_type_hint | string | 否 | access_token 或 refresh_token |

**响应**:
```json
{
  "success": true
}
```

---

## 👤 用户管理API

### 1. 用户注册

**POST** `/api/v1/auth/register`

**请求体**:
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "phone": "+1234567890"
}
```

**响应**:
```json
{
  "success": true,
  "message": "注册成功，请查收验证邮件",
  "data": {
    "user_id": 123,
    "username": "john_doe",
    "email": "john@example.com"
  }
}
```

---

### 2. 用户登录

**POST** `/api/v1/auth/login`

**请求体**:
```json
{
  "username": "john_doe",
  "password": "SecurePass123!"
}
```

**响应**:
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "tGzv3JOkF0XG5Qx2TlKWIA",
    "expires_in": 3600,
    "user": {
      "id": 123,
      "username": "john_doe",
      "email": "john@example.com"
    }
  }
}
```

---

### 3. 获取用户信息

**GET** `/api/v1/users/profile`

**请求头**:
```
Authorization: Bearer <access_token>
```

**响应**:
```json
{
  "success": true,
  "data": {
    "id": 123,
    "username": "john_doe",
    "email": "john@example.com",
    "phone": "+1234567890",
    "avatar": "http://localhost:8080/avatars/john.jpg",
    "email_verified": true,
    "roles": ["user"]
  }
}
```

---

### 4. 更新用户信息

**PUT** `/api/v1/users/profile`

**请求头**:
```
Authorization: Bearer <access_token>
```

**请求体**:
```json
{
  "phone": "+9876543210",
  "avatar": "http://localhost:8080/avatars/new-avatar.jpg"
}
```

**响应**:
```json
{
  "success": true,
  "message": "用户信息更新成功"
}
```

---

### 5. 修改密码

**POST** `/api/v1/users/password/change`

**请求头**:
```
Authorization: Bearer <access_token>
```

**请求体**:
```json
{
  "old_password": "OldPass123!",
  "new_password": "NewPass456!"
}
```

**响应**:
```json
{
  "success": true,
  "message": "密码修改成功"
}
```

---

## 🔧 管理后台API

### 1. 用户列表

**GET** `/api/v1/admin/users`

**请求头**:
```
Authorization: Bearer <admin_access_token>
```

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认20 |
| status | int | 否 | 用户状态筛选 |
| keyword | string | 否 | 关键词搜索 |

**响应**:
```json
{
  "success": true,
  "data": {
    "total": 100,
    "page": 1,
    "page_size": 20,
    "users": [
      {
        "id": 123,
        "username": "john_doe",
        "email": "john@example.com",
        "status": 1,
        "created_at": "2026-02-25T10:00:00Z"
      }
    ]
  }
}
```

---

### 2. 创建客户端

**POST** `/api/v1/admin/clients`

**请求头**:
```
Authorization: Bearer <admin_access_token>
```

**请求体**:
```json
{
  "name": "My App",
  "redirect_uris": ["http://localhost:3000/callback"],
  "scopes": ["openid", "profile", "email"],
  "logo_url": "http://example.com/logo.png",
  "description": "My application description"
}
```

**响应**:
```json
{
  "success": true,
  "data": {
    "client_id": "abc123",
    "client_secret": "secret123",
    "name": "My App"
  }
}
```

---

### 3. 创建角色

**POST** `/api/v1/admin/roles`

**请求头**:
```
Authorization: Bearer <admin_access_token>
```

**请求体**:
```json
{
  "name": "编辑",
  "code": "editor",
  "description": "内容编辑角色"
}
```

**响应**:
```json
{
  "success": true,
  "data": {
    "id": 3,
    "name": "编辑",
    "code": "editor"
  }
}
```

---

### 4. 分配角色权限

**POST** `/api/v1/admin/roles/:id/permissions`

**请求头**:
```
Authorization: Bearer <admin_access_token>
```

**路径参数**:
- `id`: 角色ID

**请求体**:
```json
{
  "permission_ids": [1, 2, 3]
}
```

**响应**:
```json
{
  "success": true,
  "message": "权限分配成功"
}
```

---

## 📨 通用响应格式

### 成功响应
```json
{
  "success": true,
  "data": { ... },
  "message": "操作成功"
}
```

### 错误响应
```json
{
  "success": false,
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "用户名或密码错误",
    "details": { ... }
  }
}
```

### 常见错误码
| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| INVALID_REQUEST | 400 | 请求参数错误 |
| UNAUTHORIZED | 401 | 未授权 |
| FORBIDDEN | 403 | 禁止访问 |
| NOT_FOUND | 404 | 资源不存在 |
| INVALID_CREDENTIALS | 401 | 认证失败 |
| INVALID_TOKEN | 401 | Token无效 |
| EXPIRED_TOKEN | 401 | Token已过期 |
| INTERNAL_ERROR | 500 | 服务器内部错误 |

---

## 🔍 OpenAPI文档

完整的OpenAPI 3.0规范文档: [[OpenAPI规范]]
