// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// Canada28BetsDao is the data access object for the table game_canada28_bets.
type Canada28BetsDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  Canada28BetsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// Canada28BetsColumns defines and stores column names for the table game_canada28_bets.
type Canada28BetsColumns struct {
	Id           string // 主键ID
	MerchantId   string // 商户ID
	UserId       string // 用户ID
	PeriodNumber string // 期号
	BetType      string // 投注类型：high/low/odd/even/num_0等
	BetName      string // 投注名称：High/Low/Number 0等
	Amount       string // 投注金额
	Multiplier   string // 投注时的赔率
	Status       string // 状态：pending-等待开奖，win-已中奖，lose-未中奖，cancel-已取消
	Ip           string // 投注IP地址
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
	DeletedAt    string // 删除时间
}

// canada28BetsColumns holds the columns for the table game_canada28_bets.
var canada28BetsColumns = Canada28BetsColumns{
	Id:           "id",
	MerchantId:   "merchant_id",
	UserId:       "user_id",
	PeriodNumber: "period_number",
	BetType:      "bet_type",
	BetName:      "bet_name",
	Amount:       "amount",
	Multiplier:   "multiplier",
	Status:       "status",
	Ip:           "ip",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
}

// NewCanada28BetsDao creates and returns a new DAO object for table data access.
func NewCanada28BetsDao(handlers ...gdb.ModelHandler) *Canada28BetsDao {
	return &Canada28BetsDao{
		group:    "default",
		table:    "game_canada28_bets",
		columns:  canada28BetsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *Canada28BetsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *Canada28BetsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *Canada28BetsDao) Columns() Canada28BetsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *Canada28BetsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *Canada28BetsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *Canada28BetsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
