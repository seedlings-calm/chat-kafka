# chat-kafka

项目目的： 使用grpc+kafka 来实现 私聊，群聊，广播等功能

使用 采用 Clean Architecture 模式构建项目结构 

根据Clean Architecture 模式 创建目录结构和测试代码


已实现：私聊，群聊，广播
未实现：群组逻辑，信息存储，消息已读未读

api操作流程：
http -> internal.api ->internal.api.router ->internal.api.dto(处理请求参数)->internal.usecase(处理逻辑)->internal.domain(处理存储)

