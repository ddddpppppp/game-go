# RabbitMQ Client

这个包提供了一个简单易用的RabbitMQ客户端，用于连接RabbitMQ服务器、创建交换机和队列，以及发布和消费消息。

## 功能

- 自动连接和重连机制
- 交换机声明和管理
- 队列声明和绑定
- 消息发布
- 消息消费
- 支持持久化消息

## 安装

确保已经安装了RabbitMQ的Go客户端库：

```bash
go get github.com/rabbitmq/amqp091-go
```

## 使用示例

### 创建客户端并连接

```go
import (
    "demo/internal/utils"
    "log"
    "time"
)

// 创建RabbitMQ客户端
client := utils.NewRabbitMQClient("amqp://guest:guest@localhost:5672/", 5*time.Second)

// 连接到RabbitMQ服务器
err := client.Connect()
if err != nil {
    log.Fatalf("无法连接到RabbitMQ: %v", err)
}
defer client.Close()
```

### 声明交换机和队列

```go
// 声明交换机
err = client.DeclareExchange(
    "my_exchange",  // 交换机名称
    "direct",       // 交换机类型
    true,           // 持久化
    false,          // 不自动删除
    false,          // 不是内部交换机
    false,          // 等待确认
    nil,            // 无额外参数
)

// 声明队列
queue, err := client.DeclareQueue(
    "my_queue",     // 队列名称
    true,           // 持久化
    false,          // 不自动删除
    false,          // 不独占
    false,          // 等待确认
    nil,            // 无额外参数
)

// 绑定队列到交换机
err = client.BindQueue(
    "my_queue",     // 队列名称
    "my_routing_key", // 路由键
    "my_exchange",  // 交换机名称
    false,          // 等待确认
    nil,            // 无额外参数
)
```

### 发送消息

```go
import "context"

// 创建上下文
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// 创建消息
message := "Hello RabbitMQ!"
msg := utils.CreateMessage(
    "text/plain",       // 内容类型
    []byte(message),    // 消息体
    true,               // 持久化消息
)

// 发送消息
err = client.PublishMessage(
    ctx,
    "my_exchange",      // 交换机名称
    "my_routing_key",   // 路由键
    false,              // 不强制投递
    false,              // 不立即投递
    msg,                // 消息
)
```

### 接收消息

要接收消息，您需要在 RabbitMQClient 结构中添加以下方法：

```go
// ConsumeMessages 消费队列中的消息
func (c *RabbitMQClient) ConsumeMessages(queueName, consumerName string, autoAck bool) (<-chan amqp.Delivery, error) {
    if !c.isConnected {
        return nil, fmt.Errorf("not connected to RabbitMQ")
    }
    
    return c.channel.Consume(
        queueName,     // 队列名称
        consumerName,  // 消费者名称（标识）
        autoAck,       // 自动确认
        false,         // 是否独占
        false,         // 阻止其他消费者消费此队列
        false,         // 等待服务器确认
        nil,           // 额外参数
    )
}
```

然后可以这样使用消费者：

```go
// 开始消费消息
msgs, err := client.ConsumeMessages(
    "my_queue",   // 队列名称
    "consumer_1", // 消费者标识
    false,        // 不自动确认，需要手动确认
)
if err != nil {
    log.Fatalf("无法开始消费消息: %v", err)
}

// 创建一个永不停止的 goroutine 来处理消息
go func() {
    for msg := range msgs {
        // 处理消息
        log.Printf("收到消息: %s", string(msg.Body))
        
        // 处理完成后手动确认消息
        err := msg.Ack(false) // false 表示只确认当前消息
        if err != nil {
            log.Printf("确认消息失败: %v", err)
        }
    }
    
    log.Println("消费者通道已关闭")
}()

// 保持应用程序运行
select {} // 永久阻塞，使消费者 goroutine 能够持续工作
```

### 完整的消费者示例

```go
package main

import (
    "log"
    "time"
    
    "demo/internal/utils"
)

func main() {
    // 创建RabbitMQ客户端
    client := utils.NewRabbitMQClient("amqp://guest:guest@localhost:5672/", 5*time.Second)
    
    // 连接到RabbitMQ服务器
    err := client.Connect()
    if err != nil {
        log.Fatalf("无法连接到RabbitMQ: %v", err)
    }
    defer client.Close()
    
    // 声明队列，确保队列存在
    queue, err := client.DeclareQueue(
        "my_queue", // 队列名称
        true,       // 持久化
        false,      // 不自动删除
        false,      // 不独占
        false,      // 等待确认
        nil,        // 无额外参数
    )
    if err != nil {
        log.Fatalf("无法声明队列: %v", err)
    }
    log.Printf("准备从队列消费消息: %s", queue.Name)
    
    // 开始消费消息
    msgs, err := client.ConsumeMessages(
        "my_queue",   // 队列名称
        "consumer_1", // 消费者标识
        false,        // 不自动确认
    )
    if err != nil {
        log.Fatalf("无法开始消费消息: %v", err)
    }
    
    // 设置消息处理函数
    go func() {
        for msg := range msgs {
            log.Printf("收到消息: %s", string(msg.Body))
            
            // 模拟消息处理
            time.Sleep(100 * time.Millisecond)
            
            // 处理完成后确认消息
            err := msg.Ack(false)
            if err != nil {
                log.Printf("确认消息失败: %v", err)
            }
        }
    }()
    
    log.Println("消费者已启动，等待消息...")
    
    // 保持应用程序运行
    select {}
}
```

### 消息确认和拒绝

在消费消息时，有几种方式处理消息：

1. **确认消息 (Ack)**：表示消息已成功处理
   ```go
   msg.Ack(false) // false 表示只确认当前消息
   ```

2. **拒绝消息 (Nack)**：表示消息处理失败
   ```go
   msg.Nack(false, true) // 第一个参数: false 表示只拒绝当前消息
                         // 第二个参数: true 表示将消息重新放入队列
   ```

3. **拒绝消息 (Reject)**：与 Nack 类似，但功能较简单
   ```go
   msg.Reject(true) // true 表示将消息重新放入队列
   ```

### 消息处理的最佳实践

1. **适当设置预取计数**：控制一次性获取的消息数量
   ```go
   err := client.channel.Qos(
       10,    // 预取计数 - 一次最多处理10条消息
       0,     // 预取大小 - 不限制消息大小
       false, // 全局设置 - false表示只应用于当前通道
   )
   ```

2. **实现幂等性**：确保消息可以被多次处理而不产生副作用

3. **使用死信队列**：处理无法成功处理的消息
   ```go
   args := amqp.Table{
       "x-dead-letter-exchange":    "dead_letter_exchange",
       "x-dead-letter-routing-key": "dead_letter_key",
   }
   
   queue, err := client.DeclareQueue("my_queue", true, false, false, false, args)
   ```

## 注意事项

1. 确保在使用完毕后调用 `Close()` 方法关闭连接。
2. 客户端会自动处理断线重连，但是在重连期间发送的消息可能会失败。
3. 处理消息时应考虑异常情况，实现适当的重试策略。
4. 建议在生产环境中增加适当的错误处理和重试机制。
5. 消费者应当优雅地处理连接断开的情况，可能需要重新设置消费者。 