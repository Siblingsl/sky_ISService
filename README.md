# 项目目录结构

## `sky_ISservice/`
这是微服务项目的根目录。

### `cmd/`                # 启动入口
- `main.go`               # 微服务的启动文件

### `services/`           # 微服务目录
每个模块代表一个微服务的功能模块，包括控制器、服务层、模块化封装、数据访问层等。

#### `auth/`              # 登录注册模块
- `controller/`           # Controller 层
    - `auth_controller.go`  # 登录注册相关的控制器
- `service/`              # Service 层
    - `auth_service.go`     # 登录注册的业务逻辑
- `module/`               # 业务模块（封装 DI、模块初始化）
    - `auth_module.go`      # 模块初始化与依赖注入
- `repository/`           # 数据访问层
    - `auth_repository.go`  # 登录注册相关数据访问逻辑

#### `system/`             # 系统设置模块
- `controller/`           # 系统设置控制器
- `service/`              # 系统设置业务层
- `module/`               # 系统设置模块初始化
- `repository/`           # 系统设置数据访问层

#### `user/`               # 用户服务模块
- `controller/`           # 用户服务控制器
- `service/`              # 用户服务业务层
- `module/`               # 用户服务模块初始化
- `repository/`           # 用户服务数据访问层

#### `order/`              # 订单服务模块
- `controller/`           # 订单服务控制器
- `service/`              # 订单服务业务层
- `module/`               # 订单服务模块初始化
- `repository/`           # 订单服务数据访问层

#### `payment/`            # 支付整合模块
- `controller/`           # 支付整合控制器
- `service/`              # 支付整合服务层
- `module/`               # 支付整合模块初始化
- `repository/`           # 支付整合数据访问层

#### `logistics/`          # 物流服务模块
- `controller/`           # 物流服务控制器
- `service/`              # 物流服务业务层
- `module/`               # 物流服务模块初始化
- `repository/`           # 物流服务数据访问层

#### `inventory/`          # 库存服务模块
- `controller/`           # 库存服务控制器
- `service/`              # 库存服务业务层
- `module/`               # 库存服务模块初始化
- `repository/`           # 库存服务数据访问层

#### `sync/`               # 数据同步模块
- `controller/`           # 数据同步控制器
- `service/`              # 数据同步业务层
- `module/`               # 数据同步模块初始化
- `repository/`           # 数据同步数据访问层

#### `notification/`       # 消息通知模块
- `controller/`           # 消息通知控制器
- `service/`              # 消息通知业务层
- `module/`               # 消息通知模块初始化
- `repository/`           # 消息通知数据访问层

#### `log/`                # 日志服务模块
- `controller/`           # 日志服务控制器
- `service/`              # 日志服务业务层
- `module/`               # 日志服务模块初始化
- `repository/`           # 日志服务数据访问层

### `shared/`             # 共享模块
- `cache/`                # Redis 缓存封装
- `logger/`               # Zap 日志封装
- `mq/`                   # 消息队列封装

### `proto/`              # gRPC proto 定义
- `auth.proto`            # 例如：登录认证服务的 proto 文件
- `order.proto`           # 订单服务的 proto 文件
- `payment.proto`         # 支付整合服务的 proto 文件

### `config/`             # 配置文件
- `config.yaml`           # 配置文件（如数据库连接、服务端口等）

### `pkg/`                # 公共库
- `database/`             # 数据库连接与模型
- `middleware/`           # 中间件
- `utils/`                # 工具函数
- `errors/`               # 错误处理

### `deployment/`         # K8s / Docker 部署
- `docker-compose.yml`    # Docker Compose 配置
- `k8s-deployment.yaml`   # Kubernetes 部署配置文件

### `go.mod`              # Go 依赖管理

### `README.md`           # 项目文档


### 目前尚未配置完成的有 统一错误处理、服务注册与发现、Redis、Swagger、CI/CD 和自动化部署

### 不要删除auth
用户登录、注册：这些是身份验证服务的核心职责，API Gateway 不应该直接处理这些操作。
Token 生成与管理：API Gateway 只负责验证 Token，而不生成或管理它们。

### 网关设置 黑白名单、访问策略、流量策略、开启 CORS、同步 swagger3.0