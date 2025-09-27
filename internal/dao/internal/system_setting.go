// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SystemSettingDao is the data access object for the table game_system_setting.
type SystemSettingDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  SystemSettingColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// SystemSettingColumns defines and stores column names for the table game_system_setting.
type SystemSettingColumns struct {
	Id          string //
	Name        string // 设置名称，如recharge_setting、withdraw_setting等
	Title       string // 设置标题
	Description string // 设置描述
	Config      string // 设置配置JSON数据
	Status      string // 状态：1启用，0禁用
	Sort        string // 排序
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
}

// systemSettingColumns holds the columns for the table game_system_setting.
var systemSettingColumns = SystemSettingColumns{
	Id:          "id",
	Name:        "name",
	Title:       "title",
	Description: "description",
	Config:      "config",
	Status:      "status",
	Sort:        "sort",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewSystemSettingDao creates and returns a new DAO object for table data access.
func NewSystemSettingDao(handlers ...gdb.ModelHandler) *SystemSettingDao {
	return &SystemSettingDao{
		group:    "default",
		table:    "game_system_setting",
		columns:  systemSettingColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SystemSettingDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SystemSettingDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SystemSettingDao) Columns() SystemSettingColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SystemSettingDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SystemSettingDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *SystemSettingDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
