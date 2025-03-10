# Auth 子服务需求文档与接口设计

### 需求背景：
Auth 子服务主要负责用户身份验证、Token 生成与管理以及用户信息的存储与查询。在这个子服务中，权限管理的逻辑将不涉及，因为权限管理将在 System 子服务中完成。Auth 子服务将包括以下功能

### 功能模块：
1. 注册：用户可以通过用户名、密码、邮箱、电话等信息注册账户。
2. 登录：用户使用已注册的用户名和密码登录，系统返回一个有效的 Token。
3. 登出：用户注销登录状态，Token 失效。
4. 身份辨别：根据用户类型（管理员或客户）进行身份辨别。
5. Token 验证：用于检查当前请求是否含有有效的 Token。
6. 邮箱验证码：发送邮箱地址，获取验证码

### 系统架构：
1. Elasticsearch：用于搜索用户数据或记录（例如快速查找用户状态）。
2. Redis：缓存用户的认证信息（如 Token），提高认证性能。
3. RabbitMQ：用于异步处理任务，如发送邮件验证码等。
4. PostgreSQL：存储用户数据（如 sky_auth_users 和 sky_auth_tokens 表）。
5. gRPC：提供高效的跨服务通信，允许不同微服务之间进行身份验证。
6. Consul：服务发现，确保服务之间的高可用性。

### 接口设计：
#### 用户注册：
* URL: /auth/register   
* 方法: POST   
* 请求体：
```json

{
  "username": "string",
  "password": "string",
  "email": "string",
  "phone": "string",
  "user_type": 2
}
```
* 响应：

成功：
```json
{
  "message": "注册成功"
}

```
失败：
```json
{
  "message": "用户名已存在"
}
```

#### 用户登陆：
* URL: /auth/login
* 方法: POST
* 请求体: email、code 仅限后台登陆时使用
```json
{
  "username": "string",
  "password": "string",
  "email": "string",
  "code": "string"
}
```
* 响应:

成功：
```json
{
  "token": "JWT Token",
  "user_type": 2
}
```
失败：
```json
{
  "message": "无效的用户名或密码"
}
```

#### 用户登出
* URL: /auth/logout
* 方法: POST
* 请求体:
```json
{
  "token": "JWT Token"
}
```
* 响应:

成功：
```json
{
  "message": "退出成功"
}
```
失败：
```json
{
  "message": "无效的令牌"
}
```

#### Token 验证:
* URL: /auth/verify-token
* 方法: GET
* 请求体:
```json
{
  "token": "JWT Token"
}
```
* 响应：
成功：
```json
{
  "message": "令牌有效"
}
```
失败：
```json
{
  "message": "令牌无效或已过期"
}
```

#### 发送邮箱验证码
* URL: /auth/send-code
* 方法: POST
* 请求体:
```json
{
  "email": "user@example.com"
}
```
* 响应

成功：
```json
{
  "message": "验证码已发送"
}
```
失败：
```json
{
  "message": "请稍候再试"
}
```

## 技术选型与示例
### Elasticsearch
用于存储用户相关的日志或搜索用户信息：
用途：在用户注册时，可以将一些日志信息（如注册日期）同步到 Elasticsearch 中进行快速检索。
示例：当用户登录失败时，可以在 Elasticsearch 中记录日志，并供管理者查看。

### Redis
用途：缓存用户的 Token 或会话信息，提高性能。
示例：用户登录后，生成的 Token 存储在 Redis 中，并设置过期时间。每次请求时，检查 Redis 中的 Token 是否有效。
```go
// 使用 Redis 验证 Token
func validateToken(token string) bool {
    val, err := redisClient.Get(token).Result()
    if err != nil {
        return false
    }
    return val == "valid"
}
```

### RabbitMQ
用途：处理异步任务，例如发送邮件验证码。
示例：在用户注册时，系统通过 RabbitMQ 发送异步请求给邮件服务来发送验证码。
```go
// 使用 RabbitMQ 发送验证码请求
func sendVerificationEmail(email string) {
    ch := rabbitMQChannel
    msg := fmt.Sprintf("Send verification email to %s", email)
    ch.Publish(
        "",       // exchange
        "queue",  // routing key
        false,    // mandatory
        false,    // immediate
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(msg),
        })
}
```

### PostgreSQL
用途：存储用户的基本信息及其认证信息。
示例：在用户登录时，使用 PostgreSQL 校验用户名和密码。
```go
// 使用 PostgreSQL 查询用户信息
func findUserByUsername(username string) (*User, error) {
    var user User
    err := db.QueryRow("SELECT id, username, password FROM sky_auth_users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password)
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

###  gRPC
用途：与其他微服务（如 system 服务）进行身份验证交互。
示例：使用 gRPC 向 system 子服务请求用户角色信息。
```proto
// gRPC 服务定义
service AuthService {
    rpc GetUserRole (UserRequest) returns (UserResponse);
}
```

### Consul
用途：服务发现，确保在多个服务实例间能够高效访问 auth 子服务。
示例：通过 Consul 查询可用的 auth 服务实例。
```go
// 使用 Consul 查找服务
func getServiceAddress(serviceName string) string {
    addr, err := consulClient.Agent().ServiceAddress(serviceName)
    if err != nil {
        log.Fatal(err)
    }
    return addr
}
```