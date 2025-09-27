package aliyunlog

import (
	"context"
	"sync"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// AliyunLogHandler 阿里云SLS日志处理器
type AliyunLogHandler struct {
	client        *sls.Client
	project       string
	logstore      string
	topic         string
	source        string
	batchSize     int
	flushInterval time.Duration
	maxRetries    int
	compressType  int // 压缩类型
	buffer        []*sls.Log
	bufferLock    sync.Mutex
	timer         *time.Timer
	closeChan     chan struct{}
}

// Config 配置结构
type Config struct {
	Endpoint        string `json:"endpoint"`        // SLS服务入口（如 "cn-hangzhou.log.aliyuncs.com"）
	AccessKeyID     string `json:"accessKeyId"`     // 阿里云AK ID
	AccessKeySecret string `json:"accessKeySecret"` // 阿里云AK Secret
	Project         string `json:"project"`         // SLS项目名称
	Logstore        string `json:"logstore"`        // 日志库名称
	Topic           string `json:"topic"`           // 日志主题（可选）
	Source          string `json:"source"`          // 日志来源标识（可选）
	BatchSize       int    `json:"batchSize"`       // 批量上传条数（默认100）
	FlushInterval   int    `json:"flushInterval"`   // 自动刷新间隔秒数（默认5）
	MaxRetries      int    `json:"maxRetries"`      // 最大重试次数（默认3）
	CompressType    int    `json:"compressType"`    // 压缩类型（默认0，不压缩; 1：lz4; 2：zlib; 3：zstd）
}

// NewAliyunLogHandler 创建处理器
func NewHandler(config Config) *AliyunLogHandler {
	if config.BatchSize <= 0 {
		config.BatchSize = 100
	}
	if config.FlushInterval <= 0 {
		config.FlushInterval = 5
	}
	if config.MaxRetries <= 0 {
		config.MaxRetries = 3
	}
	// 默认使用LZ4压缩
	if config.CompressType <= 0 {
		config.CompressType = 1 // 默认LZ4压缩
	}

	client := &sls.Client{
		Endpoint:        config.Endpoint,
		AccessKeyID:     config.AccessKeyID,
		AccessKeySecret: config.AccessKeySecret,
	}

	handler := &AliyunLogHandler{
		client:        client,
		project:       config.Project,
		logstore:      config.Logstore,
		topic:         config.Topic,
		source:        config.Source,
		batchSize:     config.BatchSize,
		flushInterval: time.Duration(config.FlushInterval) * time.Second,
		maxRetries:    config.MaxRetries,
		compressType:  config.CompressType,
		buffer:        make([]*sls.Log, 0, config.BatchSize),
		closeChan:     make(chan struct{}),
	}

	// 启动定时刷新
	handler.timer = time.AfterFunc(handler.flushInterval, handler.flushPeriodically)
	return handler
}

// Log 实现glog.Handler接口
func (h *AliyunLogHandler) Log(ctx context.Context, in *glog.HandlerInput) {
	timestamp := uint32(in.Time.Unix())
	logItem := &sls.Log{
		Time:     &timestamp,
		Contents: h.buildContents(in),
	}

	h.bufferLock.Lock()
	h.buffer = append(h.buffer, logItem)
	if len(h.buffer) >= h.batchSize {
		go h.flush() // 异步触发上传
	}
	h.bufferLock.Unlock()
}

// Close 关闭处理器
func (h *AliyunLogHandler) Close(ctx context.Context) error {
	close(h.closeChan)
	h.timer.Stop()
	return h.flush() // 最终强制上传
}

// 构建日志内容
func (h *AliyunLogHandler) buildContents(in *glog.HandlerInput) []*sls.LogContent {
	// 将日志级别转换为字符串
	levelStr := in.LevelFormat

	// 创建指针类型
	level := new(string)
	*level = levelStr

	message := new(string)
	*message = in.Content

	traceId := new(string)
	*traceId = in.TraceId

	if *traceId == "" {
		*traceId = gmd5.MustEncryptString(gtime.Now().String())
	}

	// 添加当前时间，精确到毫秒
	timeStr := new(string)
	*timeStr = time.Now().Format("2006-01-02 15:04:05.000")

	contents := []*sls.LogContent{
		{Key: stringPtr("level"), Value: level},
		{Key: stringPtr("traceId"), Value: traceId},
		{Key: stringPtr("timestamp"), Value: timeStr},
	}

	// 添加额外字段
	if len(in.Values) > 0 {
		for _, v := range in.Values {
			key := new(string)
			*key = gconv.String("content")

			value := new(string)
			*value = gconv.String(v)

			contents = append(contents, &sls.LogContent{
				Key:   key,
				Value: value,
			})
		}
	}

	return contents
}

// 辅助函数：创建字符串指针
func stringPtr(s string) *string {
	return &s
}

// 定时刷新
func (h *AliyunLogHandler) flushPeriodically() {
	select {
	case <-h.closeChan:
		return
	default:
		h.flush()
		h.timer.Reset(h.flushInterval)
	}
}

// 执行上传
func (h *AliyunLogHandler) flush() error {
	h.bufferLock.Lock()
	if len(h.buffer) == 0 {
		h.bufferLock.Unlock()
		return nil
	}

	logs := h.buffer
	h.buffer = make([]*sls.Log, 0, h.batchSize)
	h.bufferLock.Unlock()

	// 创建日志组
	logGroup := &sls.LogGroup{
		Logs:   logs,
		Topic:  &h.topic,
		Source: &h.source,
	}

	// 带重试的上传
	var err error
	for i := 0; i < h.maxRetries; i++ {
		// 使用带压缩的方法上传日志
		if err = h.client.PutLogsWithCompressType(h.project, h.logstore, logGroup, h.compressType); err == nil {
			return nil
		}
		time.Sleep(time.Second * time.Duration(i+1)) // 指数退避
	}

	// 使用一个空的上下文，因为这里我们没有可用的上下文
	backgroundCtx := context.Background()
	glog.Errorf(backgroundCtx, "SLS upload failed after %d retries: %v", h.maxRetries, err)
	return err
}
