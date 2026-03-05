# Godemo - 基于 Gin 的用户角色权限管理系统

基于 Go + Gin + GORM + SQLite3 开发的 RESTful API 服务，实现了完整的 RBAC（基于角色的访问控制）权限系统。

## 功能特性

- ✅ 用户管理（CRUD）
- ✅ 角色管理（CRUD）
- ✅ 权限管理
- ✅ JWT 认证
- ✅ 基于角色的权限控制
- ✅ CORS 跨域支持（可配置）
- ✅ 环境变量配置
- ✅ SQLite3 数据库
- ✅ Swagger API 文档
- ✅ 软删除支持

## 技术栈

- Go 1.25
- Gin Web Framework
- GORM ORM
- SQLite3
- JWT (golang-jwt/jwt)
- CORS (gin-contrib/cors)
- Swagger (swaggo)
- godotenv (环境变量管理)

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置环境变量

复制示例配置文件：

```bash
cp .env.example .env
```

编辑 `.env` 文件（可选，使用默认值也可以）：

```bash
# 服务器配置
SERVER_PORT=8080
GIN_MODE=debug

# JWT 配置
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRE_HOUR=24

# CORS 配置（开发环境允许所有来源）
CORS_ALLOW_ORIGINS=*
```

详细配置说明：

- [配置管理指南](CONFIG.md)
- [CORS 配置指南](CORS.md)

### 3. 运行服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动

### 3. 访问 Swagger 文档

打开浏览器访问：`http://localhost:8080/swagger/index.html`

## 默认账户

系统会自动创建默认管理员账户：

- 用户名：`admin`
- 密码：`admin123`
- 角色：管理员（拥有所有权限）

## API 接口

### 认证接口

- `POST /api/login` - 用户登录

### 用户管理

- `GET /api/users` - 获取用户列表（需要 `user:view` 权限）
- `GET /api/users/:id` - 获取用户详情（需要 `user:view` 权限）
- `POST /api/users` - 创建用户（需要 `user:create` 权限）
- `PUT /api/users/:id` - 更新用户（需要 `user:edit` 权限）
- `DELETE /api/users/:id` - 删除用户（需要 `user:delete` 权限）

### 角色管理

- `GET /api/roles` - 获取角色列表（需要 `role:view` 权限）
- `GET /api/roles/:id` - 获取角色详情（需要 `role:view` 权限）
- `POST /api/roles` - 创建角色（需要 `role:create` 权限）
- `PUT /api/roles/:id` - 更新角色（需要 `role:edit` 权限）
- `DELETE /api/roles/:id` - 删除角色（需要 `role:delete` 权限）

### 权限管理

- `GET /api/permissions` - 获取权限列表

## 使用示例

### 1. 登录获取 Token

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

响应：

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "type": "Bearer"
}
```

### 2. 使用 Token 访问受保护的接口

```bash
curl -X GET http://localhost:8080/api/users \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 3. 创建新用户

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "email": "test@example.com",
    "role_ids": [1]
  }'
```

## 项目结构

```
godemo/
├── main.go              # 主入口
├── config/              # 配置管理
│   └── config.go       # 配置加载
├── models/              # 数据模型
│   └── user.go         # 用户、角色、权限模型
├── handler/             # 处理器
│   ├── auth.go         # 认证处理器
│   ├── user.go         # 用户处理器
│   ├── role.go         # 角色处理器
│   └── permission.go   # 权限处理器
├── service/             # 业务逻辑
│   └── auth.go         # 认证服务
├── middleware/          # 中间件
│   ├── auth.go         # 认证和权限中间件
│   └── cors.go         # CORS 中间件
├── database/            # 数据库
│   └── database.go     # 数据库初始化
├── docs/                # Swagger 文档
├── .env.example         # 环境变量示例
└── godemo.db           # SQLite 数据库文件（运行后生成）
```

## 默认权限列表

系统预置以下权限：

- `user:view` - 查看用户
- `user:create` - 创建用户
- `user:edit` - 编辑用户
- `user:delete` - 删除用户
- `role:view` - 查看角色
- `role:create` - 创建角色
- `role:edit` - 编辑角色
- `role:delete` - 删除角色

## 注意事项

1. 生产环境请修改 `JWT_SECRET` 环境变量（详见 [CONFIG.md](CONFIG.md)）
2. 建议使用环境变量管理敏感配置
3. 默认管理员账户（ID=1）不能被删除
4. 所有删除操作都是软删除，数据不会真正从数据库中移除
5. Token 默认有效期为 24 小时，可通过 `JWT_EXPIRE_HOUR` 环境变量配置

## 开发

### 重新生成 Swagger 文档

```bash
swag init
```

### 构建

```bash
go build -o godemo
```

## License

MIT
