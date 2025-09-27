package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQClient 封装了与RabbitMQ连接的客户端
type RabbitMQClient struct {
	conn              *amqp.Connection
	channel           *amqp.Channel
	connectionString  string
	reconnectInterval time.Duration
	isConnected       bool
}

// NewRabbitMQClient 从配置文件创建一个新的RabbitMQ客户端
func NewRabbitMQClient(ctx context.Context, reconnectInterval time.Duration) (*RabbitMQClient, error) {
	// 从配置文件获取RabbitMQ配置
	host, err := g.Cfg().Get(ctx, "rabbitmq.host")
	if err != nil {
		return nil, fmt.Errorf("failed to get rabbitmq host config: %w", err)
	}

	port, err := g.Cfg().Get(ctx, "rabbitmq.port")
	if err != nil {
		return nil, fmt.Errorf("failed to get rabbitmq port config: %w", err)
	}

	user, err := g.Cfg().Get(ctx, "rabbitmq.user")
	if err != nil {
		return nil, fmt.Errorf("failed to get rabbitmq user config: %w", err)
	}

	password, err := g.Cfg().Get(ctx, "rabbitmq.password")
	if err != nil {
		return nil, fmt.Errorf("failed to get rabbitmq password config: %w", err)
	}

	vhost, err := g.Cfg().Get(ctx, "rabbitmq.vhost")
	if err != nil {
		return nil, fmt.Errorf("failed to get rabbitmq vhost config: %w", err)
	}

	// 构建RabbitMQ连接字符串
	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		user.String(),
		password.String(),
		host.String(),
		port.String(),
		vhost.String(),
	)

	// 使用默认的重连间隔（如果没有提供）
	if reconnectInterval == 0 {
		reconnectInterval = 5 * time.Second
	}

	client := &RabbitMQClient{
		connectionString:  connectionString,
		reconnectInterval: reconnectInterval,
		isConnected:       false,
	}

	return client, nil
}

// GetDefaultRabbitMQClient 获取一个默认的RabbitMQ客户端实例
func GetDefaultRabbitMQClient(ctx context.Context) (*RabbitMQClient, error) {
	// 使用默认的5秒重连间隔
	client, err := NewRabbitMQClient(ctx, 5*time.Second)
	if err != nil {
		return nil, err
	}

	// 自动连接到RabbitMQ
	if err := client.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return client, nil
}

// Connect 建立与RabbitMQ的连接
func (c *RabbitMQClient) Connect() error {
	var err error

	c.conn, err = amqp.Dial(c.connectionString)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		c.conn.Close()
		return fmt.Errorf("failed to open a channel: %w", err)
	}

	c.isConnected = true

	// 监听连接关闭
	go c.handleReconnect()

	return nil
}

// handleReconnect 处理断线重连
func (c *RabbitMQClient) handleReconnect() {
	connClose := c.conn.NotifyClose(make(chan *amqp.Error))
	for {
		reason, ok := <-connClose
		if !ok {
			// 通道已关闭，客户端可能已经被关闭
			break
		}
		log.Printf("Connection closed, reason: %v", reason)

		c.isConnected = false
		for !c.isConnected {
			log.Printf("Attempting to reconnect to RabbitMQ in %s", c.reconnectInterval)
			time.Sleep(c.reconnectInterval)

			if err := c.Connect(); err != nil {
				log.Printf("Failed to reconnect: %v", err)
				continue
			}

			log.Println("Reconnected to RabbitMQ")
		}
	}
}

// DeclareExchange 声明一个交换机
func (c *RabbitMQClient) DeclareExchange(name, kind string, durable, autoDelete, internal, noWait bool, args map[string]interface{}) error {
	if !c.isConnected {
		return fmt.Errorf("not connected to RabbitMQ")
	}

	return c.channel.ExchangeDeclare(
		name,       // 交换机名称
		kind,       // 交换机类型 (direct, fanout, topic, headers)
		durable,    // 持久化
		autoDelete, // 自动删除
		internal,   // 内部交换机
		noWait,     // 无需等待确认
		args,       // 参数
	)
}

// DeclareQueue 声明一个队列
func (c *RabbitMQClient) DeclareQueue(name string, durable, autoDelete, exclusive, noWait bool, args map[string]interface{}) (amqp.Queue, error) {
	if !c.isConnected {
		return amqp.Queue{}, fmt.Errorf("not connected to RabbitMQ")
	}

	return c.channel.QueueDeclare(
		name,       // 队列名称
		durable,    // 持久化
		autoDelete, // 自动删除
		exclusive,  // 独占
		noWait,     // 无需等待确认
		args,       // 参数
	)
}

// BindQueue 将队列绑定到交换机
func (c *RabbitMQClient) BindQueue(queueName, key, exchangeName string, noWait bool, args map[string]interface{}) error {
	if !c.isConnected {
		return fmt.Errorf("not connected to RabbitMQ")
	}

	return c.channel.QueueBind(
		queueName,    // 队列名称
		key,          // 路由键
		exchangeName, // 交换机名称
		noWait,       // 无需等待确认
		args,         // 参数
	)
}

// PublishMessage 发布消息到指定的交换机和路由键
func (c *RabbitMQClient) PublishMessage(ctx context.Context, exchangeName, routingKey string, mandatory, immediate bool, msg amqp.Publishing) error {
	if !c.isConnected {
		return fmt.Errorf("not connected to RabbitMQ")
	}

	return c.channel.PublishWithContext(
		ctx,
		exchangeName, // 交换机名称
		routingKey,   // 路由键
		mandatory,    // 强制投递
		immediate,    // 立即投递
		msg,          // 消息
	)
}

// ConsumeMessages 消费队列中的消息
func (c *RabbitMQClient) ConsumeMessages(queueName, consumerName string, autoAck bool) (<-chan amqp.Delivery, error) {
	if !c.isConnected {
		return nil, fmt.Errorf("not connected to RabbitMQ")
	}

	return c.channel.Consume(
		queueName,    // 队列名称
		consumerName, // 消费者名称（标识）
		autoAck,      // 自动确认
		false,        // 是否独占
		false,        // 阻止其他消费者消费此队列
		false,        // 等待服务器确认
		nil,          // 额外参数
	)
}

// SetQoS 设置通道的QoS（服务质量）参数
func (c *RabbitMQClient) SetQoS(prefetchCount, prefetchSize int, global bool) error {
	if !c.isConnected {
		return fmt.Errorf("not connected to RabbitMQ")
	}

	return c.channel.Qos(
		prefetchCount, // 预取计数 - 一次最多处理的消息数
		prefetchSize,  // 预取大小 - 一次传递的最大消息体大小（0表示无限制）
		global,        // 是否将设置应用于整个连接（false表示只应用于当前通道）
	)
}

// CreateMessage 创建一个消息
func CreateMessage(contentType string, body []byte, persistent bool) amqp.Publishing {
	headers := make(amqp.Table)
	headers["timestamp"] = time.Now().Unix()

	msg := amqp.Publishing{
		ContentType: contentType,
		Body:        body,
		Headers:     headers,
		Timestamp:   time.Now(),
	}

	if persistent {
		msg.DeliveryMode = amqp.Persistent
	}

	return msg
}

// Close 关闭连接和通道
func (c *RabbitMQClient) Close() {
	if c.channel != nil {
		c.channel.Close()
	}

	if c.conn != nil {
		c.conn.Close()
	}

	c.isConnected = false
}
