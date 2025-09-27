// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SystemSetting is the golang structure for table system_setting.
type SystemSetting struct {
	Id          int         `json:"id"          orm:"id"          description:""`                                         //
	Name        string      `json:"name"        orm:"name"        description:"设置名称，如recharge_setting、withdraw_setting等"` // 设置名称，如recharge_setting、withdraw_setting等
	Title       string      `json:"title"       orm:"title"       description:"设置标题"`                                     // 设置标题
	Description string      `json:"description" orm:"description" description:"设置描述"`                                     // 设置描述
	Config      string      `json:"config"      orm:"config"      description:"设置配置JSON数据"`                               // 设置配置JSON数据
	Status      int         `json:"status"      orm:"status"      description:"状态：1启用，0禁用"`                               // 状态：1启用，0禁用
	Sort        int         `json:"sort"        orm:"sort"        description:"排序"`                                       // 排序
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"  description:""`                                         //
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"  description:""`                                         //
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"  description:""`                                         //
}
