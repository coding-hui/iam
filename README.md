# WeCoding IAM

IAM = **I**dentity and **A**ccess **M**anagement

IAM 是一个基于 Go 语言开发的身份识别与访问管理系统，用于对资源访问进行授权。

## ✨ 特性

- 基于Casbin的 RBAC 访问控制模型，提供精细到按钮的权限控制

- 多种认证方式 JWT、Basic、SecretKey

- 基于 GIN WEB API 框架，提供了丰富的中间件支持（用户认证、跨域、访问日志、追踪ID等）

- 遵循 RESTful API 设计规范

- 支持 Swagger 文档(基于swaggo)

- 基于 GORM 的数据库存储，可扩展多种类型数据库

- 支持动态加载多种配置文件

- 多指令模式，提供 iamctl 命令行工具

- TODO: 多租户的支持

- TODO: 单元测试

## 🎁 内置

- 用户管理：用户是系统操作者，该功能主要完成系统用户配置。

- 组织管理：配置系统组织机构（公司、部门、小组）。

- 资源管理：资源是业务系统中具体资源的标识符，可以是一个实体，如用户，也可以是一个菜单、按钮、API。

- 权限策略 权限策略将多个资源、操作、授权作用组合在一起，为应用程序提供灵活的访问权限管理和控制功能。

- 角色管理：角色是一组权限资源的集合，可以为角色授权某些资源与操作权限。当角色授予给用户之后，该用户将会继承这个角色中的所有权限。

## 📦 本地开发

### 环境要求

- go 1.19
- node v16.19.1
- pnpm 8.5.1

### 开发目录创建

```bash
mkdir wecoding
cd wecoding
```

### 获取代码

```bash
# 获取后端代码
git clone https://github.com/coding-hui/iam.git

# 获取前端代码
git clone https://github.com/coding-hui/iam-frontend.git
```

### 启动说明

#### 服务端启动

```bash
# 进入 iam 后端项目
cd ./iam

# 构建后的二进制文件保存在 _output/platforms/linux/amd64/ 目录下。
make build

# 修改配置 
# 文件路径  iam/configs/iam-apiserver.yaml
vi ./configs/iam-apiserver.yaml

# 启动服务
# macOS or linux 下使用
go run ./cmd/iam-apiserver/main.go -c configs/iam-apiserver.yaml

# windows 下使用
go run .\cmd\iam-apiserver\main.go -c configs\iam-apiserver.yaml
```

构建后的二进制文件保存在 `_output/platforms/linux/amd64/` 目录下。

### Console UI 启动

```bash
# 安装依赖
cd iam-frontend

npm install -g pnpm

pnpm install

pnpm start
```

## 使用指南

[IAM Documentation](docs/guide/zh-CN)
