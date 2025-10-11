package customer_service_ws

import (
	"context"
	"demo/internal/dao"
	"demo/internal/model"
	"demo/internal/model/do"
	"demo/internal/service/customer_service_ws"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gorilla/websocket"
)

type WsController struct{}

type Message struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

func NewWsController() *WsController {
	return &WsController{}
}

func (c *WsController) Connect(r *ghttp.Request) {
	ctx := r.Context()

	// Authenticate user
	user := gconv.Map(r.GetCtxVar(model.CtxUserKey))
	var userID string
	var userName string
	var isAdmin bool = false

	if user != nil {
		userID = gconv.String(user["Uuid"])
		userName = gconv.String(user["Nickname"])
		if userName == "" {
			userName = gconv.String(user["Username"])
		}

		// 检查是否为admin - 直接查询game_admin表
		var adminData *do.Admin
		err := dao.Admin.Ctx(ctx).Where("uuid = ?", userID).Fields("uuid").Scan(&adminData)
		if err == nil && adminData != nil {
			// 找到admin记录，说明是管理员
			isAdmin = true
		}
	} else {
		// Reject connection if not authenticated
		r.Response.WriteStatus(401)
		r.Response.WriteJson(g.Map{
			"error": "Authentication required",
		})
		return
	}

	// Build connection ID
	connUserID := fmt.Sprintf("customer_service_%s", userID)

	// Upgrade connection and store
	ws, _ := r.WebSocket()
	conn := customer_service_ws.StoreConnection(connUserID, ws, isAdmin)

	readChan := make(chan []byte)
	errorChan := make(chan error)

	// Independent read goroutine
	go func(conn *customer_service_ws.Connection) {
		for {
			select {
			case <-conn.Ctx.Done():
				return
			default:
				_, msg, err := conn.ReadMessage()
				if err != nil {
					errorChan <- err
					return
				}
				readChan <- msg
				conn.SetReadDeadline(time.Now().Add(30 * time.Second))
			}
		}
	}(conn)

	for {
		select {
		case msg := <-readChan:
			processMessage(ctx, conn, msg, userID, userName, isAdmin)
		case err := <-errorChan:
			customer_service_ws.SafeCloseConnection(conn) // Close current instance
			handleCloseError(ctx, err, connUserID)
			return
		case <-ctx.Done():
			customer_service_ws.SafeCloseConnection(conn)
			return
		}
	}
}

func processMessage(ctx context.Context, conn *customer_service_ws.Connection, data []byte, userID string, userName string, isAdmin bool) error {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		g.Log().Error(ctx, g.Map{
			"error": "Message parsing failed",
			"data":  g.Map{"error": err.Error()},
		})
		return gerror.Wrap(err, "Message parsing failed")
	}

	switch msg.Action {
	case "heartbeat":
		responseBytes, _ := json.Marshal(g.Map{"action": "heartbeat_ack"})
		return conn.SafeWriteMessage(websocket.TextMessage, responseBytes)

	case "send_message":
		// User or admin sends message
		msgData := gconv.Map(msg.Data)
		message := gconv.String(msgData["message"])
		msgType := gconv.String(msgData["type"])

		if msgType == "" {
			msgType = "text"
		}

		if message == "" {
			return gerror.New("Message cannot be empty")
		}

		// Save message to database
		var adminIDPtr *string
		var targetUserID string

		if isAdmin {
			// Admin sending message to user
			targetUserID = gconv.String(msgData["user_id"])
			if targetUserID == "" {
				return gerror.New("Target user_id required for admin")
			}
			adminIDPtr = &userID

			// Save to database
			if err := customer_service_ws.SaveMessageToDB(ctx, targetUserID, adminIDPtr, message, msgType); err != nil {
				g.Log().Error(ctx, g.Map{
					"error": "Failed to save admin message",
					"data":  err.Error(),
				})
				return err
			}

			// Send message to target user
			targetConnID := fmt.Sprintf("customer_service_%s", targetUserID)
			responseMsg := g.Map{
				"action": "new_message",
				"data": g.Map{
					"id":         time.Now().UnixNano(),
					"user_id":    targetUserID,
					"admin_id":   userID,
					"user_name":  userName,
					"message":    message,
					"type":       msgType,
					"is_admin":   true,
					"is_read":    0,
					"created_at": time.Now().Format(time.RFC3339),
				},
			}

			responseBytes, _ := json.Marshal(responseMsg)
			if err := customer_service_ws.SendToUser(targetConnID, responseBytes); err != nil {
				g.Log().Warning(ctx, g.Map{
					"message": "Failed to send message to user, user may be offline",
					"user_id": targetUserID,
				})
			}

		} else {
			// User sending message to support
			targetUserID = userID

			// Save to database
			if err := customer_service_ws.SaveMessageToDB(ctx, targetUserID, nil, message, msgType); err != nil {
				g.Log().Error(ctx, g.Map{
					"error": "Failed to save user message",
					"data":  err.Error(),
				})
				return err
			}

			// Notify all admins
			if err := customer_service_ws.NotifyNewMessage(ctx, userID, userName, message); err != nil {
				g.Log().Warning(ctx, g.Map{
					"message": "Failed to notify admins",
					"error":   err.Error(),
				})
			}
		}

		// Send confirmation to sender
		successResponse, _ := json.Marshal(g.Map{
			"action": "message_sent",
			"data": g.Map{
				"message":    message,
				"type":       msgType,
				"created_at": time.Now().Format(time.RFC3339),
			},
		})
		return conn.SafeWriteMessage(websocket.TextMessage, successResponse)

	default:
		g.Log().Debugf(ctx, "Unknown websocket action type: %s", msg.Action)
		return nil
	}
}

func handleCloseError(ctx context.Context, err error, userID string) {
	if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
		g.Log().Info(ctx, g.Map{
			"event":   "connection_closed",
			"client":  userID,
			"type":    "websocket",
			"status":  "normal",
			"message": "Client connection closed normally",
		})
	} else if websocket.IsCloseError(err, websocket.CloseNoStatusReceived) {
		g.Log().Info(ctx, g.Map{
			"event":   "connection_closed",
			"client":  userID,
			"type":    "websocket",
			"status":  "normal",
			"message": "Client closed connection without status",
		})
	} else if websocket.IsCloseError(err, websocket.CloseGoingAway) {
		g.Log().Info(ctx, g.Map{
			"event":   "connection_closed",
			"client":  userID,
			"type":    "websocket",
			"status":  "normal",
			"message": "Client actively left",
		})
	} else if err != nil {
		g.Log().Debugf(ctx, "WebSocket read error: %v", err)
	}
}
