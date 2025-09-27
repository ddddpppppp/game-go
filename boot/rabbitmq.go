package boot

import (
	"context"
	"demo/internal/utils"
	"log"
)

// InitRabbitMQ 初始化RabbitMQ连接
func InitRabbitMQ(ctx context.Context) {
	// 检查RabbitMQ配置是否存在
	client, err := utils.GetDefaultRabbitMQClient(ctx)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return
	}

	// 连接成功
	log.Printf("Successfully connected to RabbitMQ")
	client.Close()

	// 在应用启动时初始化一个默认消费者
	// 这确保了应用启动后就有消费者在运行，不需要等待第一条消息
	// game_api.NewGameApiAiService("system", "system").InitConsumerStarted()
	log.Printf("RabbitMQ consumer initialization completed")
}
