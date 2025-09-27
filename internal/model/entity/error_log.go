// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ErrorLog is the golang structure for table error_log.
type ErrorLog struct {
	Id        int         `json:"id"        orm:"id"         description:""` //
	Content   string      `json:"content"   orm:"content"    description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" description:""` //
}
