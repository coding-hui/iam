# 用户创建邮件通知功能

## 概述

WeCoding IAM 现在支持在用户创建成功后自动发送欢迎邮件通知的功能。该功能基于事件驱动架构实现，当用户创建成功后会自动触发邮件发送流程。

## 功能特性

- **欢迎邮件**: 在新用户创建成功后，自动发送包含用户名和初始密码的欢迎邮件
- **密码重置邮件**: 支持发送密码重置邮件（预留功能）
- **配置灵活**: 支持通过配置文件启用/禁用邮件服务
- **事件驱动**: 基于事件系统，不阻塞主业务流程

## 配置说明

在 `configs/iam-apiserver.yaml` 配置文件中添加邮件配置：

```yaml
# 邮件配置
mail:
  enabled: false # 是否启用邮件服务
  host: smtp.example.com # SMTP服务器地址
  port: 587 # SMTP端口
  username: username@example.com # SMTP用户名
  password: your-smtp-password # SMTP密码
  from: no-reply@example.com # 发件人邮箱
  fromName: WeCoding IAM System # 发件人名称
```

### 配置参数说明

- `enabled`: 是否启用邮件服务，默认为 `false`
- `host`: SMTP服务器地址
- `port`: SMTP端口，默认为587
- `username`: SMTP认证用户名
- `password`: SMTP认证密码
- `from`: 发件人邮箱地址
- `fromName`: 发件人显示名称

## 使用示例

### 1. 启用邮件服务

编辑 `configs/iam-apiserver.yaml` 文件，启用邮件服务并配置正确的SMTP服务器信息：

```yaml
mail:
  enabled: true
  host: smtp.gmail.com
  port: 587
  username: your-email@gmail.com
  password: your-app-password
  from: your-email@gmail.com
  fromName: WeCoding IAM
```

### 2. 重启服务

重启 IAM API 服务器以应用新的配置：

```bash
# 如果使用二进制文件
./_output/platforms/linux/amd64/iam-apiserver -c configs/iam-apiserver.yaml

# 如果使用 Docker Compose
cd installer/docker-compose
docker-compose restart api-server
```

### 3. 创建用户

通过 API 或命令行工具创建用户，系统会自动发送欢迎邮件：

```bash
# 使用 iamctl 创建用户
iamctl user create --name=newuser --email=newuser@example.com --password=InitialPassword123
```

## 邮件模板

### 欢迎邮件模板

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>欢迎加入 WeCoding IAM 系统</title>
</head>
<body>
    <div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
        <h2 style="color: #333;">欢迎加入 {{.System}}</h2>
        <p>尊敬的 {{.Username}}，</p>
        <p>您的账户已成功创建！以下是您的登录信息：</p>
        <div style="background-color: #f5f5f5; padding: 15px; border-radius: 5px; margin: 20px 0;">
            <p><strong>用户名：</strong> {{.Username}}</p>
            <p><strong>邮箱：</strong> {{.Email}}</p>
            <p><strong>初始密码：</strong> {{.Password}}</p>
        </div>
        <p>请尽快登录系统并修改您的密码。</p>
        <p>如果您有任何问题，请联系系统管理员。</p>
        <br>
        <p>感谢您的使用！</p>
        <p>{{.System}} 团队</p>
    </div>
</body>
</html>
```

## 技术实现

### 架构设计

1. **事件驱动**: 使用内置的事件总线系统
2. **异步处理**: 邮件发送在后台异步执行，不阻塞用户创建流程
3. **依赖注入**: 通过 IoC 容器管理服务依赖

### 核心组件

- `mail.Service`: 邮件服务接口
- `UserCreatedEvent`: 用户创建事件
- `userCreatedListener`: 用户创建事件监听器
- `MailOptions`: 邮件配置选项

### 事件流程

1. 用户创建成功
2. 发布 `UserCreatedEvent` 事件
3. `userCreatedListener` 监听事件
4. 调用邮件服务发送欢迎邮件
5. 记录发送结果

## 故障排除

### 常见问题

1. **邮件发送失败**
   - 检查 SMTP 配置是否正确
   - 验证网络连接和防火墙设置
   - 检查发件人邮箱的认证信息

2. **邮件服务未启用**
   - 确保 `mail.enabled` 设置为 `true`
   - 检查配置文件的语法和缩进

3. **用户未收到邮件**
   - 检查用户邮箱地址是否正确
   - 查看邮件是否被标记为垃圾邮件
   - 检查服务日志了解发送状态

### 日志查看

邮件发送相关的日志可以在以下位置查看：

```bash
# 查看 API 服务器日志
tail -f /var/log/iam/iam-apiserver.log
```

日志中会包含邮件发送成功或失败的信息。

## 扩展开发

### 添加新的邮件模板

可以在 `internal/apiserver/domain/mail/mail.go` 中添加新的邮件模板方法。

### 自定义事件监听器

可以通过实现 `event.Listener` 接口来添加新的邮件类型监听器。