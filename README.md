# Godemo - 基于 Gin 的用户角色权限管理系统

基于 Go + Gin + GORM + SQLite3 开发的 RESTful API 服务，实现了完整的 RBAC（基于角色的访问控制）权限系统。

## 功能特性

- ✅ 用户管理（CRUD）
- ✅ 角色管理（CRUD）
- ✅ 权限管理
- ✅ JWT 认证
- ✅ 基于角色的权限控制
- ✅ CORS 跨域支持（可配置）
- ✅ 请求日志记录
- ✅ IP 地理位置解析
- ✅ 日志查询和统计
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
- IP2Region (lionsoul2014/ip2region)
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

# 日志配置
LOG_SAVE_TO_DB=true
LOG_TO_CONSOLE=true
GEOIP_DB_PATH=./ip2region.xdb
```

### 3. （可选）配置 IP2Region 地理位置功能

如果需要 IP 地理位置解析功能，下载 IP2Region 数据库（免费，无需注册）：

```bash
curl -L -o ip2region.xdb https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region.xdb
```

然后在 `.env` 文件中配置：

```bash
GEOIP_DB_PATH=./ip2region.xdb
```

IP2Region 特点：

- 完全免费，无需注册
- 数据库文件小（约 11MB）
- 查询速度快（微秒级）
- 支持国内 IP 精确到市级

### 4. 运行服务

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

### 日志管理

- `GET /api/logs` - 获取请求日志列表（支持分页和筛选）
- `GET /api/logs/:id` - 获取日志详情
- `GET /api/logs/statistics` - 获取日志统计信息

#### 日志查询参数

支持的筛选参数：

- `page`: 页码（默认 1）
- `page_size`: 每页数量（默认 20，最大 100）
- `method`: HTTP 方法（GET, POST, PUT, DELETE）
- `status`: 状态码（200, 404, 500 等）
- `ip`: 客户端 IP
- `username`: 用户名
- `country`: 国家
- `city`: 城市

示例：

```bash
# 获取状态码为 500 的错误日志
curl -X GET "http://localhost:8080/api/logs?status=500" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 获取来自中国的请求
curl -X GET "http://localhost:8080/api/logs?country=中国" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

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
│   ├── user.go         # 用户、角色、权限模型
│   └── request_log.go  # 请求日志模型
├── handler/             # 处理器
│   ├── auth.go         # 认证处理器
│   ├── user.go         # 用户处理器
│   ├── role.go         # 角色处理器
│   ├── permission.go   # 权限处理器
│   └── log.go          # 日志处理器
├── service/             # 业务逻辑
│   ├── auth.go         # 认证服务
│   └── geoip.go        # GeoIP 服务
├── middleware/          # 中间件
│   ├── auth.go         # JWT 认证中间件
│   ├── permission.go   # 权限控制中间件
│   ├── cors.go         # CORS 中间件
│   └── logging.go      # 日志中间件
├── database/            # 数据库
│   └── database.go     # 数据库初始化
├── docs/                # Swagger 文档
├── .env.example         # 环境变量示例
├── ip2region.xdb       # IP2Region 数据库（可选，需自行下载）
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

## 环境变量配置

### 服务器配置

```bash
SERVER_PORT=8080              # 服务端口
GIN_MODE=debug                # 运行模式: debug/release/test
```

### 数据库配置

```bash
DB_DRIVER=sqlite              # 数据库驱动
DB_DSN=godemo.db             # 数据库文件路径
```

### JWT 配置

```bash
JWT_SECRET=your-secret-key    # JWT 签名密钥（生产环境必须修改）
JWT_EXPIRE_HOUR=24           # Token 有效期（小时）
```

### CORS 配置

```bash
CORS_ALLOW_ORIGINS=*                              # 允许的源（生产环境指定具体域名）
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS   # 允许的方法
CORS_ALLOW_HEADERS=Origin,Content-Type,Authorization  # 允许的请求头
CORS_ALLOW_CREDENTIALS=true                       # 是否允许携带凭证
CORS_MAX_AGE=43200                               # 预检请求缓存时间（秒）
```

### 日志配置

```bash
LOG_SAVE_TO_DB=true          # 是否保存到数据库
LOG_TO_CONSOLE=true          # 是否输出到控制台
LOG_REQUEST_BODY=false       # 是否记录请求体（谨慎开启）
LOG_RESPONSE_BODY=false      # 是否记录响应体（谨慎开启）
LOG_MAX_BODY_SIZE=1024       # 请求/响应体最大记录长度（字节）
GEOIP_DB_PATH=./ip2region.xdb  # IP2Region 数据库路径
```

### 日志功能

- ✅ 自动记录所有 HTTP 请求
- ✅ IP 地理位置解析（国家、城市）
- ✅ 用户信息关联
- ✅ 异步数据库存储（不阻塞请求）
- ✅ 控制台彩色日志输出
- ✅ 支持请求/响应体记录（可配置）

### 性能优化

- 日志写入数据库是异步的，不影响请求响应速度
- IP2Region 使用内存搜索，查询速度微秒级
- 数据库表已添加索引优化查询性能
- 建议定期清理旧日志（30天以上）

### 安全建议

⚠️ **不要记录敏感信息**：

- 登录请求的密码
- JWT Token
- 个人隐私数据

**生产环境建议**：

- 关闭 `LOG_REQUEST_BODY` 和 `LOG_RESPONSE_BODY`
- 设置强 JWT 密钥（至少 32 字符）
- CORS 配置指定具体域名，不使用 `*`
- 定期更新 IP2Region 数据库
- 启用 HTTPS

## 注意事项

1. **生产环境安全**：请修改 `JWT_SECRET` 环境变量为强密钥
2. **敏感配置**：建议使用环境变量管理敏感配置
3. **管理员账户**：默认管理员账户（ID=1）不能被删除
4. **软删除**：所有删除操作都是软删除，数据不会真正从数据库中移除
5. **Token 有效期**：默认 24 小时，可通过 `JWT_EXPIRE_HOUR` 环境变量配置
6. **日志清理**：建议定期清理旧日志数据（30天以上）
7. **CORS 配置**：生产环境请指定具体的允许域名，不要使用 `*`

## 开发

### 重新生成 Swagger 文档

```bash
swag init
```

### 构建

```bash
go build -o godemo
```

### 清理旧日志

```bash
# SQLite 清理 30 天前的日志
sqlite3 godemo.db "DELETE FROM request_logs WHERE created_at < datetime('now', '-30 days');"
```

## 参考资料

- [Gin 框架文档](https://gin-gonic.com/)
- [GORM 文档](https://gorm.io/)
- [IP2Region](https://github.com/lionsoul2014/ip2region) - IP 地理位置库
- [JWT 介绍](https://jwt.io/)
- [Swagger 文档](https://swagger.io/)

## License

MIT
