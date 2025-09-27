// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentChannel is the golang structure of table game_payment_channel for DAO operations like Where/Data.
type PaymentChannel struct {
	g.Meta        `orm:"table:game_payment_channel, do:true"`
	Id            interface{} //
	Name          interface{} // 名称
	BelongAdminId interface{} // 所属管理者id
	Type          interface{} // 类型
	Rate          interface{} // 费率
	ChargeFee     interface{} // 单笔手续费
	CountTime     interface{} // 结算时间
	Guarantee     interface{} // 保证金
	FreezeTime    interface{} // 冻结时间
	DayLimitMoney interface{} // 每日限额
	DayLimitCount interface{} // 每日限额次数
	Remark        interface{} // 备注
	Status        interface{} // -1停用，1开启
	IsBackup      interface{} // 是否备用渠道
	Params        interface{} // 渠道参数
	Sort          interface{} //
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
	DeletedAt     *gtime.Time //
}
