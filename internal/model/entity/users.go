// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure for table users.
type Users struct {
	Id            int         `json:"id"            orm:"id"             description:""`                //
	Uuid          string      `json:"uuid"          orm:"uuid"           description:"用户名/手机号/邮箱"`      // 用户名/手机号/邮箱
	Username      string      `json:"username"      orm:"username"       description:"用户名/手机号/邮箱"`      // 用户名/手机号/邮箱
	Type          string      `json:"type"          orm:"type"           description:"bot/user"`        // bot/user
	Password      string      `json:"password"      orm:"password"       description:"加密后的密码"`          // 加密后的密码
	Balance       float64     `json:"balance"       orm:"balance"        description:"余额"`              // 余额
	BalanceFrozen float64     `json:"balanceFrozen" orm:"balance_frozen" description:"冻结余额"`            // 冻结余额
	Nickname      string      `json:"nickname"      orm:"nickname"       description:"昵称"`              // 昵称
	Avatar        string      `json:"avatar"        orm:"avatar"         description:"头像URL"`           // 头像URL
	MerchantId    string      `json:"merchantId"    orm:"merchant_id"    description:"商户ID"`            // 商户ID
	ParentId      int         `json:"parentId"      orm:"parent_id"      description:"邀请人ID"`           // 邀请人ID
	Status        int         `json:"status"        orm:"status"         description:"状态 (1:正常, 0:禁用)"` // 状态 (1:正常, 0:禁用)
	Salt          string      `json:"salt"          orm:"salt"           description:"密码盐值"`            // 密码盐值
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:""`                //
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     description:""`                //
	DeletedAt     *gtime.Time `json:"deletedAt"     orm:"deleted_at"     description:""`                //
}
