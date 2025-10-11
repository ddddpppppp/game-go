// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KenoBetsDao is the data access object for the table game_keno_bets.
type KenoBetsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  KenoBetsColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// KenoBetsColumns defines and stores column names for the table game_keno_bets.
type KenoBetsColumns struct {
	Id              string // 主键ID
	MerchantId      string // 商户ID
	UserId          string // 用户ID
	PeriodNumber    string // 期号
	SelectedNumbers string // 玩家选择的10个号码，JSON格式: [1,5,12,...]
	DrawnNumbers    string // 开出的20个号码，JSON格式: [3,7,12,...]
	MatchedNumbers  string // 匹配的号码，JSON格式: [12,...]
	MatchCount      string // 匹配数量：0-10
	Amount          string // 投注金额
	Multiplier      string // 投注时的赔率（基于匹配数量）
	WinAmount       string // 中奖金额
	Status          string // 状态：pending-等待开奖，win-已中奖，lose-未中奖，cancel-已取消
	Ip              string // 投注IP地址
	CreatedAt       string // 创建时间
	UpdatedAt       string // 更新时间
	SettledAt       string // 结算时间
	DeletedAt       string // 删除时间
}

// kenoBetsColumns holds the columns for the table game_keno_bets.
var kenoBetsColumns = KenoBetsColumns{
	Id:              "id",
	MerchantId:      "merchant_id",
	UserId:          "user_id",
	PeriodNumber:    "period_number",
	SelectedNumbers: "selected_numbers",
	DrawnNumbers:    "drawn_numbers",
	MatchedNumbers:  "matched_numbers",
	MatchCount:      "match_count",
	Amount:          "amount",
	Multiplier:      "multiplier",
	WinAmount:       "win_amount",
	Status:          "status",
	Ip:              "ip",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
	SettledAt:       "settled_at",
	DeletedAt:       "deleted_at",
}

// NewKenoBetsDao creates and returns a new DAO object for table data access.
func NewKenoBetsDao(handlers ...gdb.ModelHandler) *KenoBetsDao {
	return &KenoBetsDao{
		group:    "default",
		table:    "game_keno_bets",
		columns:  kenoBetsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KenoBetsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KenoBetsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KenoBetsDao) Columns() KenoBetsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KenoBetsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KenoBetsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KenoBetsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
