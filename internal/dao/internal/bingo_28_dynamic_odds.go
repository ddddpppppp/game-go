// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// Bingo28DynamicOddsDao is the data access object for the table game_bingo28_dynamic_odds.
type Bingo28DynamicOddsDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  Bingo28DynamicOddsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// Bingo28DynamicOddsColumns defines and stores column names for the table game_bingo28_dynamic_odds.
type Bingo28DynamicOddsColumns struct {
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

// bingo28DynamicOddsColumns holds the columns for the table game_bingo28_dynamic_odds.
var bingo28DynamicOddsColumns = Bingo28DynamicOddsColumns{
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

// NewBingo28DynamicOddsDao creates and returns a new DAO object for table data access.
func NewBingo28DynamicOddsDao(handlers ...gdb.ModelHandler) *Bingo28DynamicOddsDao {
	return &Bingo28DynamicOddsDao{
		group:    "default",
		table:    "game_bingo28_dynamic_odds",
		columns:  bingo28DynamicOddsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *Bingo28DynamicOddsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *Bingo28DynamicOddsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *Bingo28DynamicOddsDao) Columns() Bingo28DynamicOddsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *Bingo28DynamicOddsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *Bingo28DynamicOddsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *Bingo28DynamicOddsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
