package utils

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// OSSClient 阿里云OSS客户端
type OSSClient struct {
	client     *oss.Client
	bucketName string
	bucket     *oss.Bucket
	domain     string        // OSS自定义域名
	timeout    time.Duration // 上传超时时间
	enabled    bool          // OSS是否启用
}

// OSSConfig OSS配置
type OSSConfig struct {
	Enable          bool   `json:"enable"`            // 是否启用OSS
	Endpoint        string `json:"endpoint"`          // OSS endpoint
	AccessKeyID     string `json:"access_key_id"`     // AccessKey ID
	AccessKeySecret string `json:"access_key_secret"` // AccessKey Secret
	BucketName      string `json:"bucket_name"`       // Bucket名称
	Domain          string `json:"domain"`            // 自定义域名（可选）
	Timeout         int    `json:"timeout"`           // 超时时间（秒），默认60秒
}

var OSSInstance *OSSClient

// InitOSS 初始化OSS客户端
func InitOSS() error {
	config := OSSConfig{}

	// 从配置中读取OSS设置
	ossConfig, err := g.Cfg().Get(context.Background(), "oss")
	if err != nil {
		return gerror.Wrap(err, "读取OSS配置失败")
	}

	// 使用gconv安全地转换配置
	config.Enable = gconv.Bool(ossConfig.Map()["enable"])
	config.Endpoint = gconv.String(ossConfig.Map()["endpoint"])
	config.AccessKeyID = gconv.String(ossConfig.Map()["access_key_id"])
	config.AccessKeySecret = gconv.String(ossConfig.Map()["access_key_secret"])
	config.BucketName = gconv.String(ossConfig.Map()["bucket_name"])
	config.Domain = gconv.String(ossConfig.Map()["domain"])
	config.Timeout = gconv.Int(ossConfig.Map()["timeout"])

	// 如果OSS未启用，则不初始化客户端，但创建一个标记为未启用的实例
	if !config.Enable {
		OSSInstance = &OSSClient{
			enabled: false,
		}
		return nil
	}

	timeout := 60 * time.Second
	if config.Timeout > 0 {
		timeout = time.Duration(config.Timeout) * time.Second
	}

	// 创建OSS客户端
	client, err := oss.New(config.Endpoint, config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		return gerror.Wrap(err, "创建OSS客户端失败")
	}

	// 获取存储空间
	bucket, err := client.Bucket(config.BucketName)
	if err != nil {
		return gerror.Wrap(err, "获取OSS Bucket失败")
	}

	OSSInstance = &OSSClient{
		client:     client,
		bucketName: config.BucketName,
		bucket:     bucket,
		domain:     config.Domain,
		timeout:    timeout,
		enabled:    config.Enable,
	}

	return nil
}

// UploadFile 上传文件到OSS
// data: 文件数据
// objectKey: 对象键（OSS中的文件路径）
// contentType: 内容类型
// 返回: 可访问的URL，错误
func (o *OSSClient) UploadFile(
	ctx context.Context,
	data []byte,
	objectKey string,
	contentType string,
) (string, error) {
	// 如果OSS未启用，返回错误
	if !o.enabled {
		return "", gerror.New("OSS未启用")
	}

	if o.bucket == nil {
		return "", gerror.New("OSS客户端未初始化")
	}

	// 使用带超时的上下文
	_, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	// 上传选项
	options := []oss.Option{
		oss.ContentType(contentType),
		oss.Progress(&uploadProgressListener{}),
	}
	// 使用uploadCtx上传文件
	if err := o.bucket.PutObject(objectKey, bytes.NewReader(data), options...); err != nil {
		return "", gerror.Wrap(err, "上传文件到OSS失败")
	}

	// 构建访问URL
	var fileURL string
	if o.domain != "" {
		// 使用自定义域名
		fileURL = fmt.Sprintf("%s/%s", o.domain, objectKey)
	} else {
		// 使用默认OSS域名
		fileURL = fmt.Sprintf("https://%s.%s/%s", o.bucketName, strings.TrimPrefix(o.client.Config.Endpoint, "https://"), objectKey)
	}

	return fileURL, nil
}

// GenerateObjectKey 生成OSS对象键
func (o *OSSClient) GenerateObjectKey(contentType, directory string) string {
	// 基础文件名：时间戳+随机数
	name := gtime.Now().Format("YmdHis") + "_" + strconv.Itoa(rand.Intn(100000))

	// 根据内容类型确定扩展名
	var ext string
	switch {
	case strings.HasPrefix(contentType, "image/jpeg"):
		ext = ".jpg"
	case strings.HasPrefix(contentType, "image/png"):
		ext = ".png"
	case strings.HasPrefix(contentType, "image/gif"):
		ext = ".gif"
	case strings.HasPrefix(contentType, "image/webp"):
		ext = ".webp"
	default:
		ext = filepath.Ext(contentType)
		if ext == "" {
			ext = ".dat" // 默认扩展名
		}
	}

	// 格式化对象键
	if directory != "" {
		// 确保目录以'/'结尾但不以'/'开头
		directory = strings.TrimPrefix(directory, "/")
		if !strings.HasSuffix(directory, "/") {
			directory += "/"
		}
		return directory + name + ext
	}
	return name + ext
}

// IsEnabled 检查OSS是否启用
func (o *OSSClient) IsEnabled() bool {
	if o == nil {
		return false
	}
	return o.enabled
}

// IsOSSEnabled 全局方法检查OSS是否启用
func IsOSSEnabled() bool {
	if OSSInstance == nil {
		return false
	}
	return OSSInstance.IsEnabled()
}

// 进度监听器
type uploadProgressListener struct{}

func (listener *uploadProgressListener) ProgressChanged(event *oss.ProgressEvent) {
	switch event.EventType {
	// case oss.TransferStartedEvent:
	// 	g.Log().Debug(context.Background(), "开始上传...")
	// case oss.TransferDataEvent:
	// 	g.Log().Debug(context.Background(), "已上传", event.ConsumedBytes, "字节")
	// case oss.TransferCompletedEvent:
	// 	g.Log().Debug(context.Background(), "上传完成")
	// case oss.TransferFailedEvent:
	// 	g.Log().Error(context.Background(), "上传失败")
	}
}
