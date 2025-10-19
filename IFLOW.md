# IFLOW.md - WeCoding IAM 项目指南

## 项目概述

WeCoding IAM 是一个基于 Go 语言开发的身份识别与访问管理系统，用于对资源访问进行授权。项目采用现代化的微服务架构，包含多个独立的服务组件。

### 核心特性

- **RBAC 访问控制模型**：基于 Casbin 实现精细到按钮的权限控制
- **多种认证方式**：支持 JWT、Basic、SecretKey、OAuth、LDAP、微信小程序认证
- **多协议支持**：RESTful API + gRPC 支持
- **RESTful API 设计**：基于 Gin Web API 框架，提供丰富的中间件支持
- **多数据库支持**：基于 GORM 支持 MySQL 等多种数据库
- **Swagger 文档**：自动生成 API 文档
- **多组件架构**：包含 API Server、Authz Server 和命令行工具
- **事件驱动架构**：内置事件总线系统
- **依赖注入容器**：基于 IoC 容器的组件管理

### 技术栈

- **编程语言**: Go 1.24.0
- **Web 框架**: Gin
- **数据库**: MySQL (GORM)
- **缓存**: Redis
- **权限控制**: Casbin v2.105.0
- **API 文档**: Swagger (swaggo)
- **配置管理**: Viper
- **日志系统**: Logrus/Zap
- **gRPC 支持**: Google gRPC
- **OAuth 支持**: golang.org/x/oauth2
- **LDAP 支持**: go-ldap/ldap
- **微信集成**: silenceper/wechat
- **代码质量**: golangci-lint
- **测试框架**: Ginkgo + Gomega
- **构建工具**: Make + Docker
- **容器编排**: Docker Compose + Kubernetes + Helm

## 项目结构

```
├── cmd/                 # 命令行入口
│   ├── iam-apiserver/   # API 服务器主程序
│   ├── iam-authzserver/ # 授权服务器
│   └── iamctl/          # 命令行工具
├── internal/            # 内部包（不对外暴露）
│   ├── apiserver/       # API 服务器核心逻辑
│   │   ├── config/      # 配置管理
│   │   ├── domain/      # 领域模型
│   │   ├── event/       # 事件系统
│   │   ├── infrastructure/ # 基础设施层
│   │   ├── interfaces/  # 接口层
│   │   └── utils/       # 工具类
│   ├── authzserver/     # 授权服务器核心逻辑
│   │   ├── adapter/     # 适配器
│   │   ├── authorization/ # 授权逻辑
│   │   ├── config/      # 配置管理
│   │   ├── interfaces/  # 接口层
│   │   └── store/       # 存储层
│   ├── iamctl/          # 命令行工具逻辑
│   └── pkg/             # 内部公共包
├── pkg/                 # 可复用的公共包
│   ├── api/             # API 相关定义
│   ├── app/             # 应用框架
│   ├── code/            # 错误码定义
│   ├── container/       # 依赖注入容器
│   ├── db/              # 数据库组件
│   ├── log/             # 日志组件
│   ├── middleware/      # 中间件
│   ├── options/         # 配置选项
│   ├── server/          # 服务器组件
│   └── shutdown/        # 优雅停机
├── configs/             # 配置文件
├── hack/                # 构建和部署脚本
├── installer/           # 部署配置
└── api/                 # API 文档和测试
```

## 开发环境设置

### 环境要求

- **Go 1.24.0** (支持 1.13-1.24 版本)
- **MySQL 数据库** (推荐 8.0+)
- **Redis 缓存** (推荐 6.0+)
- **Node.js v16.19.1** (前端)
- **pnpm 8.5.1** (前端包管理)
- **Docker 和 Docker Compose** (可选，用于容器化部署)

### 快速开始

1. **克隆项目**
   ```bash
   git clone https://github.com/coding-hui/iam.git
   cd iam
   ```

2. **安装依赖**
   ```bash
   make tidy
   ```

3. **构建项目**
   ```bash
   make build
   ```

4. **配置数据库**
   - 修改 `configs/iam-apiserver.yaml` 中的数据库配置
   - 确保 MySQL 和 Redis 服务已启动
   - 或者使用 Docker Compose 快速启动所有服务

5. **启动服务**
   ```bash
   # 启动 API 服务器
   go run ./cmd/iam-apiserver/main.go -c configs/iam-apiserver.yaml
   
   # 或使用构建的二进制文件
   ./_output/platforms/linux/amd64/iam-apiserver -c configs/iam-apiserver.yaml
   
   # 使用 Docker Compose 启动所有服务
   cd installer/docker-compose
   docker-compose up -d
   ```

## 常用命令

### 构建相关

```bash
# 完整构建流程（包含代码生成、格式化、检查、构建）
make all

# 仅构建
make build

# 多平台构建（支持 linux/amd64, linux/arm64 等平台）
make build.multiarch

# 指定构建特定组件
make build BINS="iam-apiserver iam-authz-server"

# 清理构建输出
make clean
```

### 代码质量

```bash
# 代码格式化（包含 gofmt, goimports, golines）
make format

# 代码检查（使用 golangci-lint）
make lint

# 运行测试
make test

# 测试覆盖率并生成报告
make cover
```

### 开发工具

```bash
# 生成 Swagger 文档
make swag

# 启动 Swagger UI
make serve-swagger

# 生成错误码文件
make gen

# 添加版权头
make add-copyright

# 验证版权头
make verify-copyright

# 检查依赖更新
make check-updates

# 安装开发工具
make tools

# Go mod tidy
make tidy
```

### 部署相关

```bash
# 构建 Docker 镜像
make image

# 推送镜像到注册表
make push

# Kubernetes 部署（使用 Helm）
make deploy

# 卸载部署
make undeploy

# 发布构建
make release.build

# 完整发布流程
make release

# 多平台镜像构建和推送
make push.multiarch
```

### 证书管理

```bash
# 生成 CA 证书文件
make ca
```

## 配置说明

### 主要配置文件

- `configs/iam-apiserver.yaml` - API 服务器配置
- `configs/iam-apiserver-docker.yaml` - Docker 环境的 API 服务器配置
- `configs/iam-authzserver.yaml` - 授权服务器配置  
- `configs/iamctl.yaml` - 命令行工具配置

### 关键配置项

- **服务器配置**: HTTP/HTTPS/gRPC 服务端口和模式
- **数据库连接**: MySQL 主机、端口、用户名、密码
- **缓存配置**: Redis 连接信息
- **日志配置**: 日志级别、输出格式、文件路径
- **认证配置**: JWT 密钥、域名设置、OAuth 配置
- **中间件**: 启用的中间件列表
- **监控配置**: 性能分析、metrics 指标
- **特性开关**: 是否启用特定功能

## 新功能特性

### 多协议支持
- **RESTful API**: 完整的 REST API 接口
- **gRPC 服务**: 高性能的 gRPC 接口支持
- **HTTPS**: 支持 TLS/SSL 加密传输

### 身份认证扩展
- **OAuth 2.0**: 支持标准 OAuth 2.0 认证流程
- **LDAP**: 企业级 LDAP 身份验证集成
- **微信小程序**: 微信小程序身份认证支持
- **JWT Token**: JWT 令牌认证
- **Basic Auth**: 基础认证方式
- **SecretKey**: API 密钥认证

### 架构改进
- **依赖注入容器**: 基于 IoC 容器的组件管理
- **事件总线**: 事件驱动的架构设计
- **优雅停机**: 完善的进程优雅关闭机制

## 开发规范

### 代码风格

- 使用 `make format` 自动格式化代码（包含 gofmt、goimports、golines）
- 遵循 Go 语言官方代码规范
- 使用 `golangci-lint` 进行代码检查，配置见 `.golangci.yaml`
- 遵循项目版权声明规范

### 提交规范

- 提交前运行 `make lint` 确保代码质量
- 重要修改需要添加或更新测试用例
- 遵循语义化版本控制
- 确保所有代码文件包含正确的版权头

### API 设计

- 遵循 RESTful API 设计规范
- 使用统一的错误码和响应格式
- 为所有 API 端点添加 Swagger 注解
- 使用统一的认证和授权中间件

### 测试规范

- 使用 Ginkgo + Gomega 进行单元测试
- 为关键业务逻辑编写测试用例
- 确保测试覆盖率达到要求
- 使用 `make test` 运行测试，`make cover` 查看覆盖率

## 测试

### 运行测试

```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率报告
make cover
```

### 测试策略

- 单元测试覆盖核心业务逻辑
- 集成测试验证组件间交互
- API 测试确保接口正确性

## 部署

### 本地部署

1. **使用 Docker Compose (推荐)**
   ```bash
   cd installer/docker-compose
   docker-compose up -d
   ```
   服务启动后访问:
   - API 服务: http://localhost:8000
   - 前端界面: http://localhost:80
   - 数据库: localhost:13306
   - Redis: localhost:16379

2. **使用 Kubernetes**
   ```bash
   cd installer/kubernetes
   kubectl apply -f .
   ```

3. **使用 Helm Chart**
   ```bash
   cd installer/helm
   helm install iam iam/
   ```

### 生产部署

- 使用 Helm Chart 进行 Kubernetes 部署
- 配置 SSL/TLS 证书
- 设置监控和日志收集
- 配置高可用架构
- 使用外部数据库和缓存服务

## 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查 MySQL 服务状态
   - 验证配置文件中的连接信息
   - 确保数据库已创建且可访问

2. **权限验证失败**
   - 检查 JWT 密钥配置
   - 验证用户权限设置
   - 确认 RBAC 策略配置正确

3. **构建失败**
   - 确保 Go 版本符合要求 (1.13-1.24)
   - 运行 `make tidy` 更新依赖
   - 检查 `GOPATH` 和 `GOBIN` 环境变量

4. **代码检查失败**
   - 运行 `make format` 格式化代码
   - 使用 `make lint` 检查代码规范
   - 查看 `.golangci.yaml` 配置

### 日志查看

日志文件默认位置：
- API Server: `/var/log/iam/iam-apiserver.log`
- 错误日志: `/var/log/iam/iam-apiserver.error.log`
- 容器日志: `docker logs iam-apiserver` (Docker 环境)

## 贡献指南

### 提交问题

- 使用 GitHub Issues 报告 bug 或提出功能请求
- 提供详细的复现步骤和环境信息

### 提交代码

- Fork 项目并创建功能分支
- 遵循代码规范和提交规范
- 添加相应的测试用例
- 提交 Pull Request

## 相关链接

- **项目主页**: https://github.com/coding-hui/iam
- **在线演示**: http://iam.wecoding.top (默认账号: ADMIN/WECODING)
- **API 文档**: 启动服务后访问 `/swagger/index.html`
- **前端项目**: https://github.com/coding-hui/iam-frontend
- **问题反馈**: https://github.com/coding-hui/iam/issues

## 项目结构概览

```
├── cmd/                 # 命令行入口
│   ├── iam-apiserver/   # API 服务器主程序
│   ├── iam-authzserver/ # 授权服务器
│   └── iamctl/          # 命令行工具
├── internal/            # 内部包（不对外暴露）
│   ├── apiserver/       # API 服务器核心逻辑
│   │   ├── config/      # 配置管理
│   │   ├── domain/      # 领域模型
│   │   ├── event/       # 事件系统
│   │   ├── infrastructure/ # 基础设施层
│   │   ├── interfaces/  # 接口层
│   │   └── utils/       # 工具类
│   ├── authzserver/     # 授权服务器核心逻辑
│   │   ├── adapter/     # 适配器
│   │   ├── authorization/ # 授权逻辑
│   │   ├── config/      # 配置管理
│   │   ├── interfaces/  # 接口层
│   │   └── store/       # 存储层
│   ├── iamctl/          # 命令行工具逻辑
│   └── pkg/             # 内部公共包
├── pkg/                 # 可复用的公共包
│   ├── api/             # API 相关定义
│   ├── app/             # 应用框架
│   ├── code/            # 错误码定义
│   ├── container/       # 依赖注入容器
│   ├── db/              # 数据库组件
│   ├── log/             # 日志组件
│   ├── middleware/      # 中间件
│   ├── options/         # 配置选项
│   ├── server/          # 服务器组件
│   └── shutdown/        # 优雅停机
├── configs/             # 配置文件
├── hack/                # 构建和部署脚本
├── installer/           # 部署配置（Docker, Kubernetes, Helm）
├── api/                 # API 文档和测试
├── template/            # 模板文件
└── tools/               # 开发工具
```

---

*最后更新: 2025-10-19*