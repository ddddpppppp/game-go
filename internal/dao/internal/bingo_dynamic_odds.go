// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// BingoDynamicOddsDao is the data access object for the table game_bingo_dynamic_odds.
type BingoDynamicOddsDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  BingoDynamicOddsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// BingoDynamicOddsColumns defines and stores column names for the table game_bingo_dynamic_odds.
type BingoDynamicOddsColumns struct {
	Id                 string //
	MerchantId         string // 商户ID
	RuleName           string // 规则名称
	TriggerCondition   string // 触发条件：sum_range, sum_exact, sum_in
	TriggerValues      string // 触发条件值（JSON格式）
	BetTypeAdjustments string // 投注类型赔率调整（JSON格式）
	Status             string // 状态：1启用，0禁用
	Priority           string // 优先级
	CreatedAt          string //
	UpdatedAt          string //
}

// bingoDynamicOddsColumns holds the columns for the table game_bingo_dynamic_odds.
var bingoDynamicOddsColumns = BingoDynamicOddsColumns{
	Id:                 "id",
	MerchantId:         "merchant_id",
	RuleName:           "rule_name",
	TriggerCondition:   "trigger_condition",
	TriggerValues:      "trigger_values",
	BetTypeAdjustments: "bet_type_adjustments",
	Status:             "status",
	Priority:           "priority",
	CreatedAt:          "created_at",
	UpdatedAt:          "updated_at",
}

// NewBingoDynamicOddsDao creates and returns a new DAO object for table data access.
func NewBingoDynamicOddsDao(handlers ...gdb.ModelHandler) *BingoDynamicOddsDao {
	return &BingoDynamicOddsDao{
		group:    "default",
		table:    "game_bingo_dynamic_odds",
		columns:  bingoDynamicOddsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *BingoDynamicOddsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *BingoDynamicOddsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *BingoDynamicOddsDao) Columns() BingoDynamicOddsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *BingoDynamicOddsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *BingoDynamicOddsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *BingoDynamicOddsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
