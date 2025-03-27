# 项目名称

这是一个微服务架构的项目，使用 Gin、fx、Elasticsearch、Redis、RabbitMQ、PostgreSQL、Consul 等技术栈进行开发。项目主要由多个服务模块组成，每个服务模块分别处理不同的业务逻辑。

## 技术栈

- **Gin**：用于构建高性能的 Web 服务。
- **fx**：用于依赖注入和服务的生命周期管理。
- **Elasticsearch**：用于全文搜索和分析。
- **Redis**：用于缓存、消息队列和数据存储。
- **RabbitMQ**：用于异步消息传递。
- **PostgreSQL**：关系型数据库，用于存储应用数据。
- **Consul**：服务发现和配置管理。

## 技术栈

- **Gin**: 轻量级的Web框架，用于处理HTTP请求和路由。
- **fx**: 用于构建Go应用程序的依赖注入框架，简化了项目的初始化和依赖管理。
- **Elasticsearch**: 用于全文搜索和日志分析。
- **Redis**: 用作缓存存储，提高查询速度。
- **RabbitMQ**: 消息队列，用于实现服务间异步消息传递。
- **PostgreSQL**: 高效的关系型数据库，用于持久化数据。
- **Consul**: 服务发现和配置管理工具，确保服务之间的高可用性。

## 服务介绍

本项目包含多个微服务，每个服务都由独立的模块和控制器组成。各个服务包括但不限于：

- **身份验证服务** (`auth`): 提供用户认证、JWT生成、身份验证等功能。
- **订单服务** (`order`): 处理用户订单相关的逻辑。
- **系统服务** (`system`): 处理系统级的操作，如角色管理、用户管理等。

### 微服务功能

1. **身份验证服务**:
    - 提供用户注册、登录、身份验证等功能。
    - 使用JWT进行安全认证。
    - 支持与数据库和Redis缓存的交互。

2. **订单服务**:
    - 管理用户的订单数据。
    - 提供创建、查询订单等功能。

3. **系统服务**:
    - 提供用户、角色、菜单等系统级数据管理功能。

## 启动和部署

### 本地开发

1. **安装依赖**:

    ```bash
    go mod tidy
    ```

2. **启动服务**:

   使用 `fx` 启动微服务：

    ```bash
    go run services/auth/cmd/main.go
    go run services/order/cmd/main.go
    go run services/system/cmd/main.go
    ```

3. **启动网关**:

   启动网关服务，它将路由请求到相应的微服务：

    ```bash
    go run ./main.go
    ```

### 部署到生产环境

1. 配置服务的环境变量，如数据库连接、Redis和RabbitMQ配置等。
2. 使用Docker进行容器化部署，或根据具体的部署需求进行部署。

## 配置文件

- **config.yml**: 配置文件包括了服务的基本配置，例如数据库、缓存、消息队列等。
- **swagger.yaml**: 自动生成的API文档，帮助开发人员理解接口。

## 服务注册与发现

所有微服务都通过 **Consul** 进行注册与发现，确保服务的高可用性。在服务启动时，它会将自己注册到Consul中，供其他服务查询和发现。

## 使用的中间件

1. **JWT验证** (`jwt_middleware.go`): 所有需要认证的请求将通过此中间件进行验证。
2. **熔断器** (`circuit_middleware.go`): 保护服务免受不可控故障影响。
3. **日志** (`logger_middleware.go`): 所有请求都会记录日志，帮助进行故障排查。
4. **数据库相关中间件** (`db_middleware.go`): 处理数据库连接的管理。

## 扩展与自定义

根据需要，你可以轻松地扩展现有服务，或者添加新的微服务模块。每个模块都具备良好的解耦性和可扩展性，支持高并发场景。

## 参考文档

- [Gin框架](https://github.com/gin-gonic/gin)
- [fx框架](https://github.com/uber-go/fx)
- [Elasticsearch](https://www.elastic.co/)
- [Redis](https://redis.io/)
- [RabbitMQ](https://www.rabbitmq.com/)
- [PostgreSQL](https://www.postgresql.org/)
- [Consul](https://www.consul.io/)

## 许可协议

本项目遵循 [MIT License](LICENSE) 许可证。


├── build                           # 构建相关目录
│   └── sky                         # 构建相关文件
├── config                          # 配置文件目录
│   ├── config.go                   # 配置文件Go代码
│   └── config.yml                  # 配置文件YAML格式
├── deployment                      # 部署相关文档
├── dom                              # 文档目录
│   ├── auth.md                     # 身份验证相关文档
│   └── system.md                   # 系统相关文档
├── gateway                         # 网关相关目录
│   ├── build                       # 网关构建目录
│   ├── docs                        # 网关文档
│   │   ├── docs.go                 # 网关文档Go代码
│   │   ├── swagger.json            # Swagger文档JSON格式
│   │   └── swagger.yaml            # Swagger文档YAML格式
│   ├── middlewares                 # 中间件目录
│   │   └── auth_middleware.go      # 身份验证中间件
│   ├── proxy                       # 代理相关
│   │   └── proxy.go                # 代理功能实现
│   ├── router                      # 路由相关
│   │   └── router.go               # 路由配置文件
│   └── swagger                     # Swagger代码生成目录
│       └── swagger.go              # Swagger代码
├── pkg                              # 公共包目录
│   ├── grpc                         # gRPC相关功能
│   │   ├── client.go               # gRPC客户端实现
│   │   └── server.go               # gRPC服务器实现
│   ├── initialize                  # 初始化相关
│   │   └── init.go                 # 初始化文件
│   ├── middleware                  # 中间件实现
│   │   ├── circuit_middleware.go   # 熔断器中间件
│   │   ├── error_handler_middleware.go  # 错误处理中间件
│   │   ├── logger_middleware.go    # 日志中间件
│   │   ├── retry_middleware.go     # 重试中间件
│   │   ├── db_middleware.go        # 数据库相关中间件
│   │   ├── jwt_middleware.go       # JWT验证中间件
│   │   └── recovery_middleware.go  # 恢复中间件
│   └── shutdown                    # 关机/退出处理
│       └── shutdown.go             # 关机处理
├── proto                            # Proto文件目录
│   ├── auth                        # 身份验证相关Proto
│   │   └── auth.proto              # 身份验证服务的Proto文件
│   └── system                      # 系统相关Proto
│       ├── system_grpc.pb.go       # 生成的gRPC代码
│       ├── system.pb.go            # 生成的Protobuf代码
│       └── system.proto            # 系统服务Proto文件
├── services                         # 微服务目录
│   ├── auth                        # 身份验证服务
│   │   ├── cmd                     # 启动相关代码
│   │   │   └── main.go             # 启动入口文件
│   │   ├── controller              # 控制器相关
│   │   │   └── auth_controller.go  # 身份验证控制器
│   │   ├── dto                     # 数据传输对象
│   │   │   ├── req.go              # 请求DTO
│   │   │   └── res.go              # 响应DTO
│   │   ├── grpc                    # gRPC实现
│   │   │   ├── client.go           # gRPC客户端实现
│   │   │   └── server.go           # gRPC服务器实现
│   │   ├── module                  # 模块相关
│   │   │   └── auth_module.go      # 身份验证模块
│   │   ├── repository              # 数据库操作
│   │   │   ├── auth_repository.go  # 身份验证仓库
│   │   │   └── models              # 数据模型
│   │   │       ├── sky_auth_tokens.go  # 身份验证令牌模型
│   │   │       └── sky_auth_users.go   # 身份验证用户模型
│   │   ├── service                 # 服务逻辑
│   │   │   └── auth_service.go     # 身份验证服务逻辑
│   ├── order                       # 订单服务
│   │   └── cmd                     # 启动相关代码
│   │       └── main.go             # 启动入口文件
│   └── system                       # 系统服务
│       ├── cmd                      # 启动相关代码
│       │   └── main.go             # 启动入口文件
│       ├── controller               # 控制器相关
│       │   ├── menu_controller.go   # 菜单控制器
│       │   ├── role_controller.go   # 角色控制器
│       │   └── user_controller.go   # 用户控制器
│       ├── dto                      # 数据传输对象
│       │   ├── req.go               # 请求DTO
│       │   └── res.go               # 响应DTO
│       ├── grpc                     # gRPC实现
│       │   └── server.go            # gRPC服务器实现
│       ├── module                   # 模块相关
│       │   └── system_module.go     # 系统模块
│       ├── repository               # 数据库操作
│       │   ├── menu_repository.go   # 菜单仓库
│       │   ├── role_repository.go   # 角色仓库
│       │   ├── user_repository.go   # 用户仓库
│       │   └── models               # 数据模型
│       │       └── sky_system_user.go  # 用户数据模型
│       └── service                  # 服务逻辑
│           ├── menu_service.go      # 菜单服务逻辑
│           ├── role_service.go      # 角色服务逻辑
│           └── user_service.go      # 用户服务逻辑
├── shared                           # 公共服务包
│   ├── cache                        # 缓存相关
│   │   └── redis.go                # Redis操作
│   ├── elasticsearch                # Elasticsearch相关
│   │   └── elasticsearch.go         # Elasticsearch操作
│   ├── logger                       # 日志相关
│   │   └── logger.go                # 日志操作
│   ├── mq                           # 消息队列
│   │   └── rabbitmq.go             # RabbitMQ操作
│   ├── postgresql                   # PostgreSQL相关
│   │   └── postgres.go             # PostgreSQL操作
│   └── registerservice              # 服务注册
│       └── consul.go                # Consul服务注册
├── sql                              # SQL文件
│   ├── common_base.sql              # 公共数据库基础SQL
│   ├── sky_go_auth.sql              # 身份验证相关SQL
│   └── sky_go_system.sql            # 系统相关SQL
├── utils                            # 工具包
│   ├── crypto.go                    # 加密工具
│   ├── database                     # 数据库相关工具
│   │   ├── base.go                  # 基础数据库工具
│   │   └── query_helper.go          # 查询帮助函数
│   ├── errers.go                    # 错误处理工具
│   ├── jwt.go                       # JWT工具
│   ├── logs.go                      # 日志工具
│   ├── pagination.go                # 分页工具
│   ├── result.go                    # 结果处理工具
│   └── tool.go                      # 通用工具函数