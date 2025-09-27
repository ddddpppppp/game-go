// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentChannel is the golang structure for table payment_channel.
type PaymentChannel struct {
	Id            int         `json:"id"            orm:"id"              description:""`         //
	Name          string      `json:"name"          orm:"name"            description:"名称"`       // 名称
	BelongAdminId string      `json:"belongAdminId" orm:"belong_admin_id" description:"所属管理者id"`  // 所属管理者id
	Type          string      `json:"type"          orm:"type"            description:"类型"`       // 类型
	Rate          float64     `json:"rate"          orm:"rate"            description:"费率"`       // 费率
	ChargeFee     float64     `json:"chargeFee"     orm:"charge_fee"      description:"单笔手续费"`    // 单笔手续费
	CountTime     string      `json:"countTime"     orm:"count_time"      description:"结算时间"`     // 结算时间
	Guarantee     string      `json:"guarantee"     orm:"guarantee"       description:"保证金"`      // 保证金
	FreezeTime    string      `json:"freezeTime"    orm:"freeze_time"     description:"冻结时间"`     // 冻结时间
	DayLimitMoney float64     `json:"dayLimitMoney" orm:"day_limit_money" description:"每日限额"`     // 每日限额
	DayLimitCount int         `json:"dayLimitCount" orm:"day_limit_count" description:"每日限额次数"`   // 每日限额次数
	Remark        string      `json:"remark"        orm:"remark"          description:"备注"`       // 备注
	Status        int         `json:"status"        orm:"status"          description:"-1停用，1开启"` // -1停用，1开启
	IsBackup      int         `json:"isBackup"      orm:"is_backup"       description:"是否备用渠道"`   // 是否备用渠道
	Params        string      `json:"params"        orm:"params"          description:"渠道参数"`     // 渠道参数
	Sort          int         `json:"sort"          orm:"sort"            description:""`         //
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"      description:""`         //
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"      description:""`         //
	DeletedAt     *gtime.Time `json:"deletedAt"     orm:"deleted_at"      description:""`         //
}
