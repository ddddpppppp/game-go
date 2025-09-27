// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Role is the golang structure for table role.
type Role struct {
	Id         int         `json:"id"         orm:"id"          description:""`                        //
	Name       string      `json:"name"       orm:"name"        description:""`                        //
	Type       int         `json:"type"       orm:"type"        description:"1:管理员,2:商户,3:代理,4:个码管理者"` // 1:管理员,2:商户,3:代理,4:个码管理者
	MerchantId string      `json:"merchantId" orm:"merchant_id" description:"商户id"`                    // 商户id
	Access     string      `json:"access"     orm:"access"      description:""`                        //
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:""`                        //
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  description:""`                        //
	DeletedAt  *gtime.Time `json:"deletedAt"  orm:"deleted_at"  description:""`                        //
}
