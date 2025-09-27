// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Transactions is the golang structure of table game_transactions for DAO operations like Where/Data.
type Transactions struct {
	g.Meta       `orm:"table:game_transactions, do:true"`
	Id           interface{} // ID
	UserId       interface{} // 用户ID
	Type         interface{} // 交易类型: deposit-充值, withdraw-提现
	ChannelId    interface{} // 支付渠道ID
	Amount       interface{} // 交易金额
	ActualAmount interface{} // 实际金额
	Account      interface{} // 账户
	OrderNo      interface{} // 订单号
	Fee          interface{} // 手续费
	Gift         interface{} // 赠送金额
	Status       interface{} // 状态: pending-待处理, completed-已完成, failed-失败, expired-已过期
	CreatedAt    *gtime.Time // 创建时间
	CompletedAt  *gtime.Time // 完成时间
	ExpiredAt    *gtime.Time // 过期时间
	DeletedAt    *gtime.Time // 删除时间
}
