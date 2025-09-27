// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure of table game_users for DAO operations like Where/Data.
type Users struct {
	g.Meta        `orm:"table:game_users, do:true"`
	Id            interface{} //
	Uuid          interface{} // 用户名/手机号/邮箱
	Username      interface{} // 用户名/手机号/邮箱
	Type          interface{} // bot/user
	Password      interface{} // 加密后的密码
	Balance       interface{} // 余额
	BalanceFrozen interface{} // 冻结余额
	Nickname      interface{} // 昵称
	Avatar        interface{} // 头像URL
	MerchantId    interface{} // 商户ID
	ParentId      interface{} // 邀请人ID
	Status        interface{} // 状态 (1:正常, 0:禁用)
	Salt          interface{} // 密码盐值
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
	DeletedAt     *gtime.Time //
}
