package boot

import (
	"context"
	aliyunlog "demo/internal/logger"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
)

// InitLogger 初始化日志系统
func InitLogger(ctx context.Context) {
	// 从配置中读取阿里云SLS配置
	v, err := g.Cfg().Get(ctx, "logger.aliyun")
	if err != nil {
		g.Log().Error(ctx, "Failed to get aliyun config:", err)
		return
	}

	// 在 GoFrame v2 中，Map() 不返回错误
	slsConfig := v.Map()
	if len(slsConfig) == 0 {
		g.Log().Error(ctx, "Empty aliyun config")
		return
	}

	// 获取压缩类型配置（默认使用LZ4压缩 - 值为1）
	compressType := 1
	if slsConfig["compress_type"] != nil {
		compressType = gconv.Int(slsConfig["compress_type"])
	}

	// 创建阿里云日志处理器
	config := aliyunlog.Config{
		Endpoint:        gconv.String(slsConfig["endpoint"]),          // 使用gconv防止类型转换错误
		AccessKeyID:     gconv.String(slsConfig["access_key_id"]),     // 使用gconv防止类型转换错误
		AccessKeySecret: gconv.String(slsConfig["access_key_secret"]), // 使用gconv防止类型转换错误
		Project:         gconv.String(slsConfig["project"]),           // 使用gconv防止类型转换错误
		Logstore:        gconv.String(slsConfig["logstore"]),          // 使用gconv防止类型转换错误
		Topic:           gconv.String(slsConfig["topic"]),             // 使用gconv防止类型转换错误
		Source:          gconv.String(slsConfig["source"]),            // 使用gconv防止类型转换错误
		BatchSize:       100,
		FlushInterval:   5,
		MaxRetries:      3,
		CompressType:    compressType, // 添加压缩类型配置
	}

	// 创建日志处理器
	aliyunHandler := aliyunlog.NewHandler(config)

	// 配置全局日志
	// 从配置里读取是否开启阿里云日志
	aliyunEnable, _ := g.Cfg().Get(ctx, "logger.aliyun.enable")
	enable := aliyunEnable.Bool()
	if enable {
		g.Log().SetHandlers(func(ctx context.Context, in *glog.HandlerInput) {
			aliyunHandler.Log(ctx, in)
		})
	}

	// 设置日志级别，可以从配置读取
	logLevel, _ := g.Cfg().Get(ctx, "logger.level")
	if !logLevel.IsEmpty() {
		levelStr := logLevel.String()
		// 转换日志级别字符串到整数常量
		switch levelStr {
		case "all":
			g.Log().SetLevel(glog.LEVEL_ALL)
		case "dev":
			g.Log().SetLevel(glog.LEVEL_DEV)
		case "prod":
			g.Log().SetLevel(glog.LEVEL_PROD)
		// 可以添加更多级别
		default:
			g.Log().SetLevel(glog.LEVEL_ALL)
		}
	} else {
		g.Log().SetLevel(glog.LEVEL_ALL)
	}

	// 注册应用程序退出时关闭处理器
	// GoFrame v2 中没有全局的 AddShutdownHandler
	// 这里我们可以依赖应用程序在结束时正常关闭
	// 或者使用信号处理等方式管理关闭
}
