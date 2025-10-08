// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package game_bingo28_api

import (
	"context"

	v1 "demo/api/game_bingo28_api/v1"
)

type IGameBingo28ApiV1 interface {
	Bingo28Bet(ctx context.Context, req *v1.Bingo28BetReq) (res *v1.GameCommonRes, err error)
}
