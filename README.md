# GoPark API服务

GoPark是一个基于Go语言和Gin框架开发的RESTful API服务，提供用户管理和其他基础功能。

## 功能特点

- 用户管理API（创建、查询、更新、删除）
- 用户搜索和分页列表
- 健康检查接口
- 数据库自动迁移
- 统一错误处理
- API版本控制
- 中间件支持（CORS、日志、请求ID）
- 优雅关闭服务

## 技术栈

- Go 1.24
- Gin Web框架
- SQLite数据库
- go-sqlite3驱动
- Logrus日志库
- Viper配置管理

## 项目结构

```
.
├── cmd/                    # 应用入口
│   └── main.go             # 主程序
├── config/                 # 配置相关
│   ├── config.go           # 配置加载
│   └── config.yaml         # 配置文件
├── internal/               # 内部包
│   ├── db/                 # 数据库操作
│   │   ├── db.go           # 数据库连接
│   │   ├── migrations.go   # 数据库迁移
│   │   └── user.go         # 用户数据操作
│   ├── docs/               # API文档
│   │   └── swagger.go      # Swagger文档
│   ├── handlers/           # 请求处理器
│   │   ├── error.go        # 错误处理
│   │   ├── health_handler.go # 健康检查
│   │   ├── hello_handler.go # Hello处理器
│   │   └── user_handler.go  # 用户处理器
│   ├── hello/              # Hello服务
│   │   └── hello.go        # Hello实现
│   ├── middleware/         # 中间件
│   │   └── middleware.go   # 中间件实现
│   ├── migrations/         # SQL迁移文件
│   │   └── 001_create_users_table.sql # 用户表创建
│   ├── models/             # 数据模型
│   │   └── user.go         # 用户模型
│   ├── routes/             # 路由配置
│   │   └── routes.go       # 路由注册
│   └── server/             # 服务器
│       └── server.go       # HTTP服务器
├── go.mod                  # Go模块文件
└── go.sum                  # Go依赖校验
```

## 快速开始

### 前置条件

- Go 1.24或更高版本
- 无需安装额外的数据库（使用内置SQLite）

### 安装

1. 克隆仓库

```bash
git clone https://github.com/yourusername/gopark.git
cd gopark
```

2. 安装依赖

```bash
go mod download
```

3. 配置数据库

编辑`config/config.yaml`文件，设置SQLite数据库文件路径：

```yaml
database:
  type: sqlite
  path: ./gopark.db
```

4. 运行应用

```bash
go run cmd/main.go
```

应用将在配置的端口上启动（默认为8080）。

### API端点

#### 健康检查

```
GET /health
```

#### 用户管理

```
GET    /api/v1/users?id=1          # 获取用户
POST   /api/v1/users               # 创建用户
PUT    /api/v1/users/1             # 更新用户
DELETE /api/v1/users/1             # 删除用户
GET    /api/v1/users/search?name=test # 搜索用户
GET    /api/v1/users/list?limit=10&offset=0 # 列出用户
```

#### Hello服务

```
GET /api/v1/hello
```

### 示例请求

#### 创建用户

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","mail":"john@example.com"}'
```

#### 获取用户

```bash
curl http://localhost:8080/api/v1/users?id=1
```

#### 搜索用户

```bash
curl http://localhost:8080/api/v1/users/search?name=John
```

## 测试

运行单元测试：

```bash
go test ./...
```

## 部署

### Docker部署

1. 构建Docker镜像

```bash
docker build -t gopark .
```

2. 运行容器

```bash
docker run -p 8080:8080 gopark
```

## 贡献

欢迎提交问题和拉取请求。

## 许可证

本项目采用MIT许可证 - 详见LICENSE文件。