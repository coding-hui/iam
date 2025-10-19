# API Key 管理接口文档

本文档详细说明 IAM 系统中 API Key 管理功能的 REST API 接口，供前端开发使用。

## 时间格式说明

所有时间字段必须使用 **RFC3339** 格式，例如：`"2025-12-31T23:59:59Z"`。

## REST API 接口

API Key 管理功能为前端提供了完整的 CRUD 操作接口，支持细粒度的权限控制和 IP 限制。

## 🔗 API 端点

### 1. 创建 API Key
- **方法**: `POST /api/v1/apikeys`
- **认证**: Bearer Token
- **权限**: `apikeys` 权限

**请求体**:
```json
{
  "name": "API Key 名称",
  "description": "描述信息",
  "expiresAt": "2025-12-31T23:59:59Z",
  "permissions": {
    "roles": ["admin", "user"],
    "actions": ["read", "write"],
    "scopes": ["api"],
    "resources": [
      {
        "resourceType": "user",
        "resourceIds": ["user-1", "user-2"],
        "actions": ["read"]
      }
    ]
  },
  "allowedIps": {
    "ips": ["192.168.1.1", "10.0.0.1"],
    "cidrs": ["10.0.0.0/8", "192.168.0.0/16"]
  }
}
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "instanceId": "apikey-1234567890",
    "name": "API Key 名称",
    "key": "sk-ec3336a7a7548c6344a34a0fdf1a6b54",
    "secret": "64位十六进制密钥（仅创建时返回一次）",
    "status": 1,
    "expiresAt": "2025-12-31T23:59:59Z",
    "permissions": {...},
    "allowedIps": {...}
  }
}
```

### 2. 获取 API Key 列表
- **方法**: `GET /api/v1/apikeys`
- **查询参数**:
  - `page` - 页码（默认 1）
  - `limit` - 每页数量（默认 10）
  - `status` - 状态筛选（1=激活，0=禁用）

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "totalCount": 15,
    "items": [
      {
        "instanceId": "apikey-1234567890",
        "name": "API Key 名称",
        "key": "sk-ec3336a7a7548c6344a34a0fdf1a6b54",
        "status": 1,
        "expiresAt": "2025-12-31T23:59:59Z",
        "lastUsedAt": "2025-10-19T10:30:00Z",
        "usageCount": 15
      }
    ]
  }
}
```

### 3. 获取单个 API Key
- **方法**: `GET /api/v1/apikeys/{instanceId}`

**响应**: 与创建接口响应类似，但不包含 secret

### 4. 更新 API Key
- **方法**: `PUT /api/v1/apikeys/{instanceId}`
- **请求体**: 与创建接口相同，可更新任意字段

### 5. 删除 API Key
- **方法**: `DELETE /api/v1/apikeys/{instanceId}`

### 6. 重新生成密钥
- **方法**: `POST /api/v1/apikeys/{instanceId}/regenerate`
- **响应**: 返回新的密钥（仅返回一次）

### 7. 启用/禁用 API Key
- **启用**: `PUT /api/v1/apikeys/{instanceId}/enable`
- **禁用**: `PUT /api/v1/apikeys/{instanceId}/disable`

## 🔑 API Key 格式

### 密钥格式
- **API Key**: `sk-{32位十六进制字符}`
- **示例**: `sk-ec3336a7a7548c6344a34a0fdf1a6b54`

### 认证方式
前端可以使用两种方式认证：

**方式1: Authorization Header**
```
Authorization: Bearer sk-ec3336a7a7548c6344a34a0fdf1a6b54:secret_key_here
```

**方式2: 自定义 Header**
```
X-API-Key: sk-ec3336a7a7548c6344a34a0fdf1a6b54
X-API-Secret: secret_key_here
```

## ⚙️ 权限控制

### 权限结构
```json
{
  "permissions": {
    "roles": ["admin", "user"],        // 可承担的角色
    "actions": ["read", "write"],      // 可执行的操作
    "scopes": ["api"],                 // 作用域
    "resources": [                      // 资源权限
      {
        "resourceType": "user",
        "resourceIds": ["user-1"],     // 空数组表示所有资源
        "actions": ["read", "update"]
      }
    ]
  }
}
```

### IP 限制
```json
{
  "allowedIps": {
    "ips": ["192.168.1.1"],            // 单个 IP 地址
    "cidrs": ["10.0.0.0/8"]            // CIDR 网段
  }
}
```

## 📊 状态管理

### API Key 状态
- `1` - 激活 (Active)
- `0` - 禁用 (Inactive)
- `2` - 已过期 (Expired)

### 使用统计
每个 API Key 包含使用统计信息：
- `lastUsedAt` - 最后使用时间
- `usageCount` - 使用次数

## 🚨 错误码

前端需要处理的错误码：
- `110801` - API Key 不存在
- `110802` - API Key 已存在
- `110803` - 无效的 API Key 或密钥
- `110804` - API Key 未激活
- `110805` - API Key 已过期
- `110806` - API Key 已启用
- `110807` - API Key 已禁用

## 💡 使用建议

1. **密钥安全**: 创建时务必保存好密钥，系统不会再次显示
2. **权限最小化**: 为每个 API Key 分配最小必要权限
3. **IP 限制**: 强烈建议设置 IP 白名单
4. **定期轮换**: 定期重新生成密钥
5. **监控使用**: 关注使用统计，及时发现异常

## 🔗 集成示例

```javascript
// 创建 API Key
const createApiKey = async (apiKeyData) => {
  const response = await fetch('/api/v1/apikeys', {
    method: 'POST',
    headers: {
      'Authorization': 'Bearer user_jwt_token',
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(apiKeyData)
  });
  return await response.json();
};

// 使用 API Key 认证
const apiCallWithApiKey = async (endpoint, data) => {
  const response = await fetch(endpoint, {
    method: 'POST',
    headers: {
      'X-API-Key': 'sk-ec3336a7a7548c6344a34a0fdf1a6b54',
      'X-API-Secret': 'secret_key_here',
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  });
  return await response.json();
};
```

这个接口为前端提供了完整的 API Key 管理能力，支持细粒度的权限控制和安全的程序化访问。