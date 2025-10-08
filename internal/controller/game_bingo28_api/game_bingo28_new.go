package game_bingo28_api

import game_bingo28_api "demo/api/game_bingo28_api"

type ControllerV1 struct{}

func NewV1() game_bingo28_api.IGameBingo28ApiV1 {
	return &ControllerV1{}
}
