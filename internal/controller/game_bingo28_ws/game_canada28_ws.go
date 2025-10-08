package game_bingo28_ws

import (
	"context"
	"demo/internal/model"
	"demo/internal/service/game_bingo28_ws"
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

	// 认证用户
	// 客户端标识
	user := gconv.Map(r.GetCtxVar(model.CtxUserKey))
	var userID string
	if user != nil {
		userID = fmt.Sprintf("game_conversation_user_%s", gconv.String(user["Uuid"]))
	} else {
		userID = fmt.Sprintf("game_conversation_guest_%s", r.GetClientIp())
	}

	// 升级连接并存储
	ws, _ := r.WebSocket()
	conn := game_bingo28_ws.StoreConnection(userID, ws)

	readChan := make(chan []byte)
	errorChan := make(chan error)

	// 独立读取协程
	go func(conn *game_bingo28_ws.Connection) {
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
			processMessage(ctx, conn, msg)
		case err := <-errorChan:
			game_bingo28_ws.SafeCloseConnection(conn) // 关闭当前实例
			handleCloseError(ctx, err, userID)
			return
		case <-ctx.Done():
			game_bingo28_ws.SafeCloseConnection(conn)
			return
		}
	}
}

func processMessage(ctx context.Context, conn *game_bingo28_ws.Connection, data []byte) error {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		g.Log().Error(ctx, g.Map{
			"error": "消息解析失败",
			"data":  g.Map{"error": err.Error()},
		})
		return gerror.Wrap(err, "消息解析失败")
	}

	switch msg.Action {
	case "heartbeat":
		responseBytes, _ := json.Marshal(g.Map{"action": "heartbeat_ack"})
		return conn.SafeWriteMessage(websocket.TextMessage, responseBytes)
	case "send_msg_to_user":
		g.Log().Info(ctx, g.Map{
			"message": "接收ws消息",
			"data":    msg.Data,
		})
		successResponse, _ := json.Marshal(g.Map{"action": "success", "data": msg.Data})
		return conn.SafeWriteMessage(websocket.TextMessage, successResponse)
	default:
		g.Log().Debugf(ctx, "websocket未知操作类型: %s", msg.Action)
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
			"message": "客户端连接正常关闭",
		})
	} else if websocket.IsCloseError(err, websocket.CloseNoStatusReceived) {
		g.Log().Info(ctx, g.Map{
			"event":   "connection_closed",
			"client":  userID,
			"type":    "websocket",
			"status":  "normal",
			"message": "客户端无状态关闭连接",
		})
	} else if websocket.IsCloseError(err, websocket.CloseGoingAway) {
		g.Log().Info(ctx, g.Map{
			"event":   "connection_closed",
			"client":  userID,
			"type":    "websocket",
			"status":  "normal",
			"message": "客户端主动离开",
		})
	} else if err != nil {
		g.Log().Debugf(ctx, "WebSocket read error: %v", err)
	}
}
