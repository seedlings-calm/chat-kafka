# 项目目录结构
```
chat-app/
├── cmd/
│   └── main.go              # 启动程序的入口
│
├── config/
│   └── config.yaml          # 配置文件
│
├── internal/
│   ├── api/
│   │   ├── dto/             # 数据传输对象（请求和响应模型）
│   │   │   ├── user_dto.go     # 用户相关请求和响应模型
│   │   │   ├── group_dto.go    # 群组相关请求和响应模型
│   │   │   └── chat_dto.go     # 聊天室相关请求和响应模型
│   │   ├── handler/         # API 处理器
│   │   │   ├── user_handler.go   # 处理用户相关请求（如添加好友）
│   │   │   ├── group_handler.go  # 处理群组相关请求（如创建群）
│   │   │   └── chat_handler.go   # 处理聊天室相关请求（如进入聊天、退出聊天室）
│   │   └── router/          # 路由配置
│   │       └── router.go
│   │
│   ├── domain/
│   │   ├── user/            # 用户领域
│   │   │   ├── user.go     # 用户领域模型
│   │   │   └── user_service.go # 用户领域服务接口
│   │   ├── group/           # 群组领域
│   │   │   ├── group.go    # 群组领域模型
│   │   │   └── group_service.go # 群组领域服务接口
│   │   └── chat/            # 聊天室领域
│   │       ├── chat.go     # 聊天室领域模型
│   │       └── chat_service.go # 聊天室领域服务接口
│   │
│   ├── usecase/
│   │   ├── user_usecase.go  # 用户相关业务逻辑
│   │   ├── group_usecase.go # 群组相关业务逻辑
│   │   └── chat_usecase.go  # 聊天室相关业务逻辑
│   │
│   ├── infrastructure/
│   │   ├── kafka/           # Kafka 相关实现
│   │   │   ├── producer.go
│   │   │   └── consumer.go
│   │   ├── grpc/            # gRPC 相关实现
│   │   │   ├── server.go    # gRPC 服务器实现
│   │   │   ├── client.go    # gRPC 客户端实现
│   │   │   └── proto/       # Proto 文件及生成的代码
│   │   │       ├── chat.proto   # Proto 文件
│   │   │       └── chat.pb.go   # 生成的 Proto 代码
│   │   ├── mysql/           # MySQL 数据库操作
│   │   │   ├── mysql.go
│   │   │   └── user_repository.go
│   │   ├── redis/           # Redis 数据库操作
│   │   │   ├── redis.go
│   │   │   └── cache.go
│   │   └── config/          # Viper 配置操作
│   │       └── viper.go
│   │
│  

```



