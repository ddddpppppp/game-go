// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Admin is the golang structure for table admin.
type Admin struct {
	Id         int         `json:"id"         orm:"id"          description:""`                          //
	Uuid       string      `json:"uuid"       orm:"uuid"        description:""`                          //
	Nickname   string      `json:"nickname"   orm:"nickname"    description:""`                          //
	Avatar     string      `json:"avatar"     orm:"avatar"      description:""`                          //
	Username   string      `json:"username"   orm:"username"    description:""`                          //
	Balance    float64     `json:"balance"    orm:"balance"     description:"余额"`                        // 余额
	Password   string      `json:"password"   orm:"password"    description:""`                          //
	Salt       string      `json:"salt"       orm:"salt"        description:""`                          //
	MerchantId string      `json:"merchantId" orm:"merchant_id" description:"对应merchant表"`               // 对应merchant表
	RoleId     int         `json:"roleId"     orm:"role_id"     description:"对应role表"`                   // 对应role表
	ParentId   string      `json:"parentId"   orm:"parent_id"   description:"推荐人ID"`                     // 推荐人ID
	Path       string      `json:"path"       orm:"path"        description:"路径（记录所有上级ID包括自己，如0:1:2:3）"` // 路径（记录所有上级ID包括自己，如0:1:2:3）
	Depth      int         `json:"depth"      orm:"depth"       description:"层级深度"`                      // 层级深度
	Status     int         `json:"status"     orm:"status"      description:"-1冻结，1开启"`                  // -1冻结，1开启
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:""`                          //
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  description:""`                          //
	DeletedAt  *gtime.Time `json:"deletedAt"  orm:"deleted_at"  description:""`                          //
}
