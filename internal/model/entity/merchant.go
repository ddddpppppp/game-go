// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Merchant is the golang structure for table merchant.
type Merchant struct {
	Id        int         `json:"id"        orm:"id"         description:""`          //
	AdminId   string      `json:"adminId"   orm:"admin_id"   description:"admin表的id"` // admin表的id
	Balance   float64     `json:"balance"   orm:"balance"    description:"余额"`        // 余额
	Uuid      string      `json:"uuid"      orm:"uuid"       description:""`          //
	Name      string      `json:"name"      orm:"name"       description:""`          //
	Logo      string      `json:"logo"      orm:"logo"       description:""`          //
	AppKey    string      `json:"appKey"    orm:"app_key"    description:""`          //
	Type      int         `json:"type"      orm:"type"       description:""`          //
	Status    int         `json:"status"    orm:"status"     description:""`          //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`          //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`          //
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" description:""`          //
}
