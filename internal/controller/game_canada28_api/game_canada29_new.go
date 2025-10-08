package game_canada28_api

import "demo/api/game_canada28_api"

type ControllerV1 struct{}

func NewV1() game_canada28_api.IGameApiV1 {
	return &ControllerV1{}
}
