Language : [us](./README.md) | [🇨🇳](./README_zh.md)

<h1 align="center">WeCoding IAM</h1>

<div align="center">

IAM = **I**dentity and **A**ccess **M**anagement

An identity and access management system developed in Go, used for authorizing resource access.

</div>

- Preview: http://iam.wecoding.top (Default login: ADMIN/WECODING)
- FAQ: https://github.com/coding-hui/iam/issues

## ✨ Features

- RBAC access control model based on Casbin, providing fine-grained permission control down to buttons.

- Multiple authentication methods: JWT, Basic, SecretKey.

- Built on the GIN WEB API framework, providing rich middleware support (user authentication, CORS, access logs, trace ID, etc.).

- Follows RESTful API design specifications.

- Supports Swagger documentation (based on swaggo).

- Database storage based on GORM, extensible to multiple types of databases.

- Supports dynamically loading multiple configuration files.

- Multi-command mode, providing the iamctl command-line tool.

- TODO: Support for multi-tenancy.

- TODO: Unit tests.

## 🎁 Built-in

- User Management: Users are system operators, this function mainly completes system user configuration.

- Organization Management: Configures system organization (company, department, group).

- Resource Management: Resources are identifiers of specific resources in the business system, which can be an entity, such as a user, or a menu, button, API.

- Permission Policy: Permission policies combine multiple resources, operations, and authorization effects to provide flexible access permission management and control functions for applications.

- Role Management: A role is a collection of permission resources, which can authorize certain resources and operation permissions to the role. When a role is granted to a user, the user will inherit all permissions of this role.

## 📦 Local Development

### Environment Requirements

- go 1.19
- node v16.19.1
- pnpm 8.5.1

### Development Directory Creation

```bash
mkdir wecoding
cd wecoding
```

### Get Code

```bash
# Get backend code
git clone https://github.com/coding-hui/iam.git

# Get frontend code
git clone https://github.com/coding-hui/iam-frontend.git
```

### Start Instructions

#### Server Start

```bash
# Enter iam backend project
cd ./iam

# Build
make build

# Modify configuration
# File path iam/configs/iam-apiserver.yaml
vi ./configs/iam-apiserver.yaml

# Start service
# macOS or linux
go run ./cmd/iam-apiserver/main.go -c configs/iam-apiserver.yaml

# windows
go run .\cmd\iam-apiserver\main.go -c configs\iam-apiserver.yaml
```

The built binary file is saved in the `_output/platforms/linux/amd64/` directory.

#### Console UI Start

```bash
# Install dependencies
cd iam-frontend

npm install -g pnpm

pnpm install

pnpm start

```

Visit: http://localhost:8000 ADMIN/WECODING

## User Guide

[IAM Documentation](docs/guide/en)

## Contributing

We welcome contributions:

- Submit [issues](https://github.com/coding-hui/iam/issues) to report bugs or ask questions.
- Propose [pull requests](https://github.com/coding-hui/iam/pulls) to improve our code.
