// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KenoBetTypes is the golang structure for table keno_bet_types.
type KenoBetTypes struct {
	Id          int         `json:"id"          orm:"id"          description:""`           //
	MerchantId  string      `json:"merchantId"  orm:"merchant_id" description:"商户ID"`       // 商户ID
	TypeName    string      `json:"typeName"    orm:"type_name"   description:"玩法名称"`       // 玩法名称
	TypeKey     string      `json:"typeKey"     orm:"type_key"    description:"玩法标识"`       // 玩法标识
	Description string      `json:"description" orm:"description" description:"玩法描述"`       // 玩法描述
	Odds        float64     `json:"odds"        orm:"odds"        description:"赔率倍数"`       // 赔率倍数
	Status      int         `json:"status"      orm:"status"      description:"状态：1启用，0禁用"` // 状态：1启用，0禁用
	Sort        int         `json:"sort"        orm:"sort"        description:"排序"`         // 排序
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"  description:""`           //
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"  description:""`           //
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"  description:""`           //
}
