// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ErrorLog is the golang structure of table game_error_log for DAO operations like Where/Data.
type ErrorLog struct {
	g.Meta    `orm:"table:game_error_log, do:true"`
	Id        interface{} //
	Content   interface{} //
	CreatedAt *gtime.Time //
	DeletedAt *gtime.Time //
}
