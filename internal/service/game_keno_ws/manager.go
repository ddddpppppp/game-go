package game_keno_ws

import (
	"context"
	"demo/internal/dao"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gorilla/websocket"
)

// Connection 代表单个WebSocket连接
type Connection struct {
	*ghttp.WebSocket
	UserID     string
	ConnID     string // 新增: 每个连接的唯一标识符
	LastActive time.Time
	mu         sync.RWMutex // 读写锁，用于保护连接状态
	writeMu    sync.Mutex   // 专门的写入锁，确保WebSocket写入操作的串行化
	closed     int32
	closeMutex sync.Mutex
	Ctx        context.Context
	cancel     context.CancelFunc
}

// UserConnections 用于存储单个用户的所有连接
type UserConnections struct {
	Connections map[string]*Connection // map[connID]*Connection
	mu          sync.RWMutex
}

var (
	// connections 改为存储每个用户的多个连接
	connections     = sync.Map{}       // map[userID]*UserConnections
	connectionIDGen int64              // 用于生成唯一的连接ID
	cleanupInterval = 1 * time.Minute  // 更频繁的清理检查
	inactiveTimeout = 5 * time.Minute  // 更短的超时时间
	pingInterval    = 30 * time.Second // ping间隔
)

// 初始化函数，启动清理任务
func init() {
	go cleanupInactiveConnections()
}

// 清理不活跃的连接
func cleanupInactiveConnections() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		connections.Range(func(userIDKey, userConnsValue interface{}) bool {
			userID := userIDKey.(string)
			userConns := userConnsValue.(*UserConnections)

			userConns.mu.Lock()
			defer userConns.mu.Unlock()

			// 检查该用户的每个连接
			for connID, conn := range userConns.Connections {
				conn.mu.RLock()
				inactive := now.Sub(conn.LastActive) > inactiveTimeout
				isClosed := conn.IsClosed()
				conn.mu.RUnlock()

				// 检查不活跃或已关闭的连接
				if inactive || isClosed {
					// 关闭不活跃的连接
					if !isClosed {
						// 安全地关闭连接
						go SafeCloseConnection(conn)
					}
					delete(userConns.Connections, connID)
				}
			}

			// 如果用户没有任何连接了，则从全局map中删除
			if len(userConns.Connections) == 0 {
				connections.Delete(userID)
			}

			return true
		})
	}
}

// IsClosed 检查连接状态
func (c *Connection) IsClosed() bool {
	return atomic.LoadInt32(&c.closed) == 1
}

// IsHealthy 检查连接是否健康（移除ping测试，避免并发写入）
func (c *Connection) IsHealthy() bool {
	if c.IsClosed() {
		return false
	}

	// 检查context是否被取消
	select {
	case <-c.Ctx.Done():
		return false
	default:
		return true
	}
}

// SafeWriteMessage 安全的写入消息方法，确保串行化
func (c *Connection) SafeWriteMessage(messageType int, data []byte) error {
	if c.IsClosed() {
		return fmt.Errorf("连接已关闭")
	}

	// 使用专门的写入锁确保串行化
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	// 再次检查连接状态
	if c.IsClosed() {
		return fmt.Errorf("连接已关闭")
	}

	err := c.WriteMessage(messageType, data)
	if err != nil {
		// 写入失败时标记连接为已关闭
		atomic.StoreInt32(&c.closed, 1)
		return err
	}

	// 更新活跃时间
	c.mu.Lock()
	c.LastActive = time.Now()
	c.mu.Unlock()

	return nil
}

// GetConnections 获取用户的所有活跃连接
func GetConnections(userID string) ([]*Connection, bool) {
	v, ok := connections.Load(userID)
	if !ok {
		return nil, false
	}

	userConns := v.(*UserConnections)
	userConns.mu.RLock()
	defer userConns.mu.RUnlock()

	result := make([]*Connection, 0, len(userConns.Connections))
	for _, conn := range userConns.Connections {
		if !conn.IsClosed() {
			result = append(result, conn)
		}
	}

	return result, len(result) > 0
}

// GetConnection 获取单个连接（保留向后兼容性，返回第一个活跃连接）
func GetConnection(userID string) (*Connection, bool) {
	conns, ok := GetConnections(userID)
	if !ok || len(conns) == 0 {
		return nil, false
	}
	return conns[0], true
}

// UpdateActiveTime 更新连接活跃时间
func (c *Connection) UpdateActiveTime() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.LastActive = time.Now()
}

// StoreConnection 存储新连接（不再删除旧连接）
func StoreConnection(userID string, conn *ghttp.WebSocket) *Connection {
	// 生成唯一的连接ID
	connID := fmt.Sprintf("%s_%d", userID, atomic.AddInt64(&connectionIDGen, 1))

	// 创建带取消功能的新连接
	ctx, cancel := context.WithCancel(context.Background())
	newConn := &Connection{
		WebSocket:  conn,
		UserID:     userID,
		ConnID:     connID,
		LastActive: time.Now(),
		Ctx:        ctx,
		cancel:     cancel,
	}

	// 获取或创建用户的连接集合
	v, _ := connections.LoadOrStore(userID, &UserConnections{
		Connections: make(map[string]*Connection),
	})
	userConns := v.(*UserConnections)

	// 添加新连接
	userConns.mu.Lock()
	userConns.Connections[connID] = newConn
	userConns.mu.Unlock()

	return newConn
}

// SafeCloseConnection 安全关闭单个连接
func SafeCloseConnection(conn *Connection) {
	conn.closeMutex.Lock()
	defer conn.closeMutex.Unlock()

	if atomic.CompareAndSwapInt32(&conn.closed, 0, 1) {
		conn.cancel()
		conn.Close()
	}

	// 从用户的连接集合中移除此连接
	if v, ok := connections.Load(conn.UserID); ok {
		userConns := v.(*UserConnections)
		userConns.mu.Lock()
		delete(userConns.Connections, conn.ConnID)

		// 如果用户没有任何连接了，则从全局map中删除
		if len(userConns.Connections) == 0 {
			connections.Delete(conn.UserID)
		}
		userConns.mu.Unlock()
	}
}

// SendToUser 向用户的所有连接发送消息（使用新的安全写入方法）
func SendToUser(userID string, message []byte) error {
	conns, ok := GetConnections(userID)
	if !ok || len(conns) == 0 {
		return nil // 用户没有活跃连接，但这不是错误
	}

	var lastErr error
	var successCount int

	for _, conn := range conns {
		// 检查连接是否健康
		if !conn.IsHealthy() {
			// 连接不健康，异步关闭并清理
			go SafeCloseConnection(conn)
			continue
		}

		// 使用单独的goroutine发送消息，避免阻塞
		go func(connection *Connection) {
			if err := connection.SafeWriteMessage(websocket.TextMessage, message); err != nil {
				// 发送失败时关闭连接
				go SafeCloseConnection(connection)
			}
		}(conn)
		successCount++
	}

	// 如果没有任何健康的连接，返回错误
	if successCount == 0 {
		return fmt.Errorf("no healthy connections for user %s", userID)
	}

	return lastErr
}

// SendUserMessage 根据管理员UUID推送WebSocket消息
func SendUserMessage(ctx context.Context, userUuid string, message interface{}) error {
	if userUuid == "" {
		return gerror.New("管理员UUID不能为空")
	}

	// 构建用户ID格式
	userID := fmt.Sprintf("game_conversation_%s", userUuid)

	// 将消息序列化为JSON
	jsonMessage, err := gjson.Encode(message)
	if err != nil {
		return gerror.Wrap(err, "消息序列化失败")
	}

	// 发送消息到该管理员的所有连接
	if err := SendToUser(userID, jsonMessage); err != nil {
		return gerror.Wrapf(err, "向管理员 %s 发送消息失败", userUuid)
	}

	return nil
}

// GetActiveCount 获取活跃连接数
func GetActiveCount() int {
	count := 0
	connections.Range(func(_, userConnsValue interface{}) bool {
		userConns := userConnsValue.(*UserConnections)
		userConns.mu.RLock()
		count += len(userConns.Connections)
		userConns.mu.RUnlock()
		return true
	})
	return count
}

// GetActiveUserCount 获取活跃用户数
func GetActiveUserCount() int {
	count := 0
	connections.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

// BroadcastMessage 向所有活跃用户广播消息
func BroadcastMessage(message []byte) error {
	var totalConnections int
	var successCount int
	var lastErr error

	connections.Range(func(userIDKey, userConnsValue interface{}) bool {
		userConns := userConnsValue.(*UserConnections)

		userConns.mu.RLock()
		conns := make([]*Connection, 0, len(userConns.Connections))
		for _, conn := range userConns.Connections {
			if !conn.IsClosed() {
				conns = append(conns, conn)
				totalConnections++
			}
		}
		userConns.mu.RUnlock()

		// 向该用户的所有连接发送消息
		for _, conn := range conns {
			if !conn.IsHealthy() {
				go SafeCloseConnection(conn)
				continue
			}

			// 使用goroutine异步发送，避免阻塞
			go func(connection *Connection) {
				if err := connection.SafeWriteMessage(websocket.TextMessage, message); err != nil {
					go SafeCloseConnection(connection)
				}
			}(conn)
			successCount++
		}

		return true
	})

	if totalConnections == 0 {
		return fmt.Errorf("没有活跃的连接")
	}

	if successCount == 0 {
		return fmt.Errorf("没有健康的连接可以发送消息")
	}

	return lastErr
}

// BroadcastToAllUsers 根据消息内容向所有用户广播WebSocket消息
func BroadcastToAllUsers(ctx context.Context, message g.Map) error {
	// 将消息序列化为JSON
	jsonMessage, err := gjson.Encode(message)
	if err != nil {
		return gerror.Wrap(err, "消息序列化失败")
	}

	action := gconv.String(message["action"])

	if action == "new_message" {
		var data = gconv.Map(message["data"])
		dao.GroupMessage.Ctx(ctx).Insert(g.Map{
			"user_id":  gconv.String(data["user_id"]),
			"group_id": "keno",
			"message":  gconv.String(data["message"]),
			"type":     gconv.String(data["type"]),
		})
	}

	// 广播消息到所有活跃连接
	if err := BroadcastMessage(jsonMessage); err != nil {
		return gerror.Wrap(err, "广播消息失败")
	}

	return nil
}
