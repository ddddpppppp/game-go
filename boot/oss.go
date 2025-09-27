package boot

import (
	"context"
	"demo/internal/utils"

	"github.com/gogf/gf/v2/frame/g"
)

// InitializeOSS 初始化阿里云OSS
func InitializeOSS() {
	ctx := context.Background()

	// 检查OSS配置是否存在
	ossConfig, err := g.Cfg().Get(ctx, "oss")
	if err != nil {
		g.Log().Error(ctx, "Failed to get OSS config:", err)
		return
	}

	if ossConfig.IsEmpty() {
		g.Log().Warning(ctx, "OSS configuration is empty, OSS upload will not be available")
		return
	}

	// 检查OSS是否启用
	enabled, err := g.Cfg().Get(ctx, "oss.enable")
	if err != nil {
		g.Log().Error(ctx, "Failed to get OSS enable config:", err)
	}

	// 初始化OSS客户端
	if err := utils.InitOSS(); err != nil {
		g.Log().Error(ctx, "Failed to initialize OSS client:", err)
		return
	}

	if enabled.Bool() {
		g.Log().Info(ctx, "OSS client initialized successfully and is enabled")
	} else {
		g.Log().Info(ctx, "OSS client initialized successfully but is disabled")
	}
}
