// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package game_canada28_api

import (
	"context"

	v1 "demo/api/game_canada28_api/v1"
)

type IGameApiV1 interface {
	Canada28Bet(ctx context.Context, req *v1.Canada28BetReq) (res *v1.GameCommonRes, err error)
}
