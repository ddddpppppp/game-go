// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// Bingo28DrawsDao is the data access object for the table game_bingo28_draws.
type Bingo28DrawsDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  Bingo28DrawsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// Bingo28DrawsColumns defines and stores column names for the table game_bingo28_draws.
type Bingo28DrawsColumns struct {
	Id            string // 主键ID
	PeriodNumber  string // 期号，如：3333197
	Status        string // 状态：0-等待开奖，1-开奖中，2-已开奖，3-已结算
	StartAt       string // 开始投注时间
	EndAt         string // 停止投注时间
	DrawAt        string // 开奖时间
	ResultNumbers string // 开奖号码，JSON格式存储三个数字
	ResultSum     string // 开奖结果总和(0-27)
	CreatedAt     string // 创建时间
	UpdatedAt     string // 更新时间
	DeletedAt     string // 删除时间
}

// bingo28DrawsColumns holds the columns for the table game_bingo28_draws.
var bingo28DrawsColumns = Bingo28DrawsColumns{
	Id:            "id",
	PeriodNumber:  "period_number",
	Status:        "status",
	StartAt:       "start_at",
	EndAt:         "end_at",
	DrawAt:        "draw_at",
	ResultNumbers: "result_numbers",
	ResultSum:     "result_sum",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
}

// NewBingo28DrawsDao creates and returns a new DAO object for table data access.
func NewBingo28DrawsDao(handlers ...gdb.ModelHandler) *Bingo28DrawsDao {
	return &Bingo28DrawsDao{
		group:    "default",
		table:    "game_bingo28_draws",
		columns:  bingo28DrawsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *Bingo28DrawsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *Bingo28DrawsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *Bingo28DrawsDao) Columns() Bingo28DrawsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *Bingo28DrawsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *Bingo28DrawsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *Bingo28DrawsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
