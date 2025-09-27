// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Transactions is the golang structure for table transactions.
type Transactions struct {
	Id           uint64      `json:"id"           orm:"id"            description:"ID"`                                                     // ID
	UserId       uint64      `json:"userId"       orm:"user_id"       description:"用户ID"`                                                   // 用户ID
	Type         string      `json:"type"         orm:"type"          description:"交易类型: deposit-充值, withdraw-提现"`                          // 交易类型: deposit-充值, withdraw-提现
	ChannelId    string      `json:"channelId"    orm:"channel_id"    description:"支付渠道ID"`                                                 // 支付渠道ID
	Amount       float64     `json:"amount"       orm:"amount"        description:"交易金额"`                                                   // 交易金额
	ActualAmount float64     `json:"actualAmount" orm:"actual_amount" description:"实际金额"`                                                   // 实际金额
	Account      string      `json:"account"      orm:"account"       description:"账户"`                                                     // 账户
	OrderNo      string      `json:"orderNo"      orm:"order_no"      description:"订单号"`                                                    // 订单号
	Fee          float64     `json:"fee"          orm:"fee"           description:"手续费"`                                                    // 手续费
	Gift         float64     `json:"gift"         orm:"gift"          description:"赠送金额"`                                                   // 赠送金额
	Status       string      `json:"status"       orm:"status"        description:"状态: pending-待处理, completed-已完成, failed-失败, expired-已过期"` // 状态: pending-待处理, completed-已完成, failed-失败, expired-已过期
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:"创建时间"`                                                   // 创建时间
	CompletedAt  *gtime.Time `json:"completedAt"  orm:"completed_at"  description:"完成时间"`                                                   // 完成时间
	ExpiredAt    *gtime.Time `json:"expiredAt"    orm:"expired_at"    description:"过期时间"`                                                   // 过期时间
	DeletedAt    *gtime.Time `json:"deletedAt"    orm:"deleted_at"    description:"删除时间"`                                                   // 删除时间
}
