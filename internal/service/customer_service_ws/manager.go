package customer_service_ws

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
	"github.com/gorilla/websocket"
)

// Connection represents a single WebSocket connection
type Connection struct {
	*ghttp.WebSocket
	UserID     string
	ConnID     string // Unique identifier for each connection
	IsAdmin    bool   // Whether this is an admin connection
	LastActive time.Time
	mu         sync.RWMutex // Read-write lock to protect connection state
	writeMu    sync.Mutex   // Write lock to ensure serial WebSocket write operations
	closed     int32
	closeMutex sync.Mutex
	Ctx        context.Context
	cancel     context.CancelFunc
}

// UserConnections stores all connections for a single user
type UserConnections struct {
	Connections map[string]*Connection // map[connID]*Connection
	mu          sync.RWMutex
}

var (
	// connections changed to store multiple connections per user
	connections     = sync.Map{}       // map[userID]*UserConnections
	connectionIDGen int64              // For generating unique connection IDs
	cleanupInterval = 1 * time.Minute  // More frequent cleanup checks
	inactiveTimeout = 5 * time.Minute  // Shorter timeout
	pingInterval    = 30 * time.Second // ping interval
)

// init function to start cleanup task
func init() {
	go cleanupInactiveConnections()
}

// cleanupInactiveConnections cleans up inactive connections
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

			// Check each connection for this user
			for connID, conn := range userConns.Connections {
				conn.mu.RLock()
				inactive := now.Sub(conn.LastActive) > inactiveTimeout
				isClosed := conn.IsClosed()
				conn.mu.RUnlock()

				// Check inactive or closed connections
				if inactive || isClosed {
					// Close inactive connections
					if !isClosed {
						// Safely close the connection
						go SafeCloseConnection(conn)
					}
					delete(userConns.Connections, connID)
				}
			}

			// If user has no connections left, delete from global map
			if len(userConns.Connections) == 0 {
				connections.Delete(userID)
			}

			return true
		})
	}
}

// IsClosed checks connection state
func (c *Connection) IsClosed() bool {
	return atomic.LoadInt32(&c.closed) == 1
}

// IsHealthy checks if connection is healthy (removed ping test to avoid concurrent writes)
func (c *Connection) IsHealthy() bool {
	if c.IsClosed() {
		return false
	}

	// Check if context is cancelled
	select {
	case <-c.Ctx.Done():
		return false
	default:
		return true
	}
}

// SafeWriteMessage safely writes message, ensuring serialization
func (c *Connection) SafeWriteMessage(messageType int, data []byte) error {
	if c.IsClosed() {
		return fmt.Errorf("connection closed")
	}

	// Use dedicated write lock to ensure serialization
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	// Check connection state again
	if c.IsClosed() {
		return fmt.Errorf("connection closed")
	}

	err := c.WriteMessage(messageType, data)
	if err != nil {
		// Mark connection as closed on write failure
		atomic.StoreInt32(&c.closed, 1)
		return err
	}

	// Update active time
	c.mu.Lock()
	c.LastActive = time.Now()
	c.mu.Unlock()

	return nil
}

// GetConnections gets all active connections for a user
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

// GetConnection gets a single connection (maintains backward compatibility, returns first active connection)
func GetConnection(userID string) (*Connection, bool) {
	conns, ok := GetConnections(userID)
	if !ok || len(conns) == 0 {
		return nil, false
	}
	return conns[0], true
}

// UpdateActiveTime updates connection active time
func (c *Connection) UpdateActiveTime() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.LastActive = time.Now()
}

// StoreConnection stores new connection (no longer deletes old connections)
func StoreConnection(userID string, conn *ghttp.WebSocket, isAdmin bool) *Connection {
	// Generate unique connection ID
	connID := fmt.Sprintf("%s_%d", userID, atomic.AddInt64(&connectionIDGen, 1))

	// Create new connection with cancel functionality
	ctx, cancel := context.WithCancel(context.Background())
	newConn := &Connection{
		WebSocket:  conn,
		UserID:     userID,
		ConnID:     connID,
		IsAdmin:    isAdmin,
		LastActive: time.Now(),
		Ctx:        ctx,
		cancel:     cancel,
	}

	// Get or create user's connection set
	v, _ := connections.LoadOrStore(userID, &UserConnections{
		Connections: make(map[string]*Connection),
	})
	userConns := v.(*UserConnections)

	// Add new connection
	userConns.mu.Lock()
	userConns.Connections[connID] = newConn
	userConns.mu.Unlock()

	return newConn
}

// SafeCloseConnection safely closes a single connection
func SafeCloseConnection(conn *Connection) {
	conn.closeMutex.Lock()
	defer conn.closeMutex.Unlock()

	if atomic.CompareAndSwapInt32(&conn.closed, 0, 1) {
		conn.cancel()
		conn.Close()
	}

	// Remove this connection from user's connection set
	if v, ok := connections.Load(conn.UserID); ok {
		userConns := v.(*UserConnections)
		userConns.mu.Lock()
		delete(userConns.Connections, conn.ConnID)

		// If user has no connections left, delete from global map
		if len(userConns.Connections) == 0 {
			connections.Delete(conn.UserID)
		}
		userConns.mu.Unlock()
	}
}

// SendToUser sends message to all user connections (using new safe write method)
func SendToUser(userID string, message []byte) error {
	conns, ok := GetConnections(userID)
	if !ok || len(conns) == 0 {
		return nil // User has no active connections, but this is not an error
	}

	var lastErr error
	var successCount int

	for _, conn := range conns {
		// Check if connection is healthy
		if !conn.IsHealthy() {
			// Connection is not healthy, close and clean up asynchronously
			go SafeCloseConnection(conn)
			continue
		}

		// Use separate goroutine to send message to avoid blocking
		go func(connection *Connection) {
			if err := connection.SafeWriteMessage(websocket.TextMessage, message); err != nil {
				// Close connection on send failure
				go SafeCloseConnection(connection)
			}
		}(conn)
		successCount++
	}

	// If no healthy connections, return error
	if successCount == 0 {
		return fmt.Errorf("no healthy connections for user %s", userID)
	}

	return lastErr
}

// SendMessageToUser sends WebSocket message based on user UUID
func SendMessageToUser(ctx context.Context, userUuid string, message interface{}) error {
	if userUuid == "" {
		return gerror.New("User UUID cannot be empty")
	}

	// Build user ID format
	userID := fmt.Sprintf("customer_service_%s", userUuid)

	// Serialize message to JSON
	jsonMessage, err := gjson.Encode(message)
	if err != nil {
		return gerror.Wrap(err, "Message serialization failed")
	}

	// Send message to all connections of this user
	if err := SendToUser(userID, jsonMessage); err != nil {
		return gerror.Wrapf(err, "Failed to send message to user %s", userUuid)
	}

	return nil
}

// GetActiveCount gets active connection count
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

// GetActiveUserCount gets active user count
func GetActiveUserCount() int {
	count := 0
	connections.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

// BroadcastMessage broadcasts message to all active users
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

		// Send message to all connections for this user
		for _, conn := range conns {
			if !conn.IsHealthy() {
				go SafeCloseConnection(conn)
				continue
			}

			// Use goroutine to send asynchronously to avoid blocking
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
		return fmt.Errorf("no active connections")
	}

	if successCount == 0 {
		return fmt.Errorf("no healthy connections to send message")
	}

	return lastErr
}

// BroadcastToAllAdmins broadcasts message to all admin connections
func BroadcastToAllAdmins(ctx context.Context, message g.Map) error {
	// Serialize message to JSON
	jsonMessage, err := gjson.Encode(message)
	if err != nil {
		return gerror.Wrap(err, "Message serialization failed")
	}

	var totalAdminConnections int
	var successCount int

	connections.Range(func(userIDKey, userConnsValue interface{}) bool {
		userConns := userConnsValue.(*UserConnections)

		userConns.mu.RLock()
		for _, conn := range userConns.Connections {
			if conn.IsAdmin && !conn.IsClosed() {
				totalAdminConnections++
				if conn.IsHealthy() {
					go func(connection *Connection) {
						if err := connection.SafeWriteMessage(websocket.TextMessage, jsonMessage); err != nil {
							go SafeCloseConnection(connection)
						}
					}(conn)
					successCount++
				} else {
					go SafeCloseConnection(conn)
				}
			}
		}
		userConns.mu.RUnlock()

		return true
	})

	if totalAdminConnections == 0 {
		return fmt.Errorf("no active admin connections")
	}

	if successCount == 0 {
		return fmt.Errorf("no healthy admin connections to send message")
	}

	return nil
}

// SaveMessageToDB saves message to database
func SaveMessageToDB(ctx context.Context, userID string, adminID *string, message string, msgType string) error {
	data := g.Map{
		"user_id": userID,
		"message": message,
		"type":    msgType,
	}

	if adminID != nil {
		data["admin_id"] = *adminID
	}

	_, err := dao.CustomerServiceMessage.Ctx(ctx).Insert(data)
	if err != nil {
		return gerror.Wrap(err, "Failed to save message to database")
	}

	// Update session
	session := g.Map{
		"user_id":         userID,
		"last_message":    message,
		"last_message_at": time.Now(),
	}

	// If admin sent message, update admin_id
	if adminID != nil {
		session["admin_id"] = *adminID
	}

	// Check if session exists
	count, err := dao.CustomerServiceSession.Ctx(ctx).Where("user_id", userID).Count()
	if err != nil {
		return gerror.Wrap(err, "Failed to check session")
	}

	if count > 0 {
		// Update existing session
		// 如果是admin发送消息，增加unread_count（用户未读）
		// 如果是用户发送消息，不增加unread_count（admin未读由前端处理）
		if adminID != nil {
			// Admin发送消息给用户，增加用户的未读数
			_, err = dao.CustomerServiceSession.Ctx(ctx).
				Where("user_id", userID).
				Data(session).
				Increment("unread_count", 1)
		} else {
			// 用户发送消息，只更新session信息，不增加未读数
			_, err = dao.CustomerServiceSession.Ctx(ctx).
				Where("user_id", userID).
				Update(session)
		}
	} else {
		// Create new session
		session["status"] = 1
		// 如果是admin发送的第一条消息，设置unread_count为1
		if adminID != nil {
			session["unread_count"] = 1
		} else {
			session["unread_count"] = 0
		}
		_, err = dao.CustomerServiceSession.Ctx(ctx).Insert(session)
	}

	if err != nil {
		return gerror.Wrap(err, "Failed to update session")
	}

	return nil
}

// NotifyNewMessage notifies admin of new message from user
func NotifyNewMessage(ctx context.Context, userID string, userName string, message string) error {
	notification := g.Map{
		"action": "new_user_message",
		"data": g.Map{
			"user_id":   userID,
			"user_name": userName,
			"message":   message,
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}

	return BroadcastToAllAdmins(ctx, notification)
}
