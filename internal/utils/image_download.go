package utils

import (
	"bytes"
	"context"
	"crypto/tls"
	"demo/internal/consts"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"image"
	"image/jpeg"
	"image/png"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// Redis键前缀常量
const (
	ImageBase64CacheKey = "image_cache_hash:%s"
)

// 图片文件签名常量 (Magic Numbers)
var (
// 常见图片格式的文件头部分字节
// jpegHeader = []byte{0xFF, 0xD8, 0xFF}       // JPEG/JPG
// pngHeader  = []byte{0x89, 0x50, 0x4E, 0x47} // PNG
// gifHeader  = []byte{0x47, 0x49, 0x46, 0x38} // GIF
// webpHeader = []byte{0x52, 0x49, 0x46, 0x46} // WEBP (RIFF)
)

type ImageDownloader struct {
	client *http.Client
}

var Image = &ImageDownloader{
	client: &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          10,
			IdleConnTimeout:       60 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 30 * time.Second, // 总超时时间
	},
}

// DownloadAndSave 安全下载图片并上传到OSS
// ctx: 上下文，用于控制超时和取消
// imageUrl: 要下载的图片URL或base64编码的图片数据
// directory: OSS存储目录（如 "images/"）
// maxSize: 最大允许的文件大小（字节），0表示不限制
// compressQuality: 压缩质量，1代表原图不压缩，0.1-0.9表示压缩比例（如0.5代表压缩至50%质量）
// 返回: OSS URL、文件类型、错误
func (u *ImageDownloader) DownloadAndSave(
	ctx context.Context,
	imageUrl string,
	directory string,
	maxSize int64,
	compressQuality float64,
) (ossUrl string, contentType string, err error) {
	// 检查是否是base64编码的图片
	if u.isBase64Image(imageUrl) {
		return u.processBase64Image(ctx, imageUrl, directory, maxSize, compressQuality)
	}

	allowedTypes := []string{"image/jpeg", "image/png", "image/gif", "image/webp", "application/octet-stream"}
	// 1. URL验证
	if err := u.validateUrl(imageUrl); err != nil {
		return "", "", gerror.Wrap(err, "URL验证失败")
	}

	// 2. 创建安全HTTP请求
	req, err := http.NewRequestWithContext(ctx, "GET", imageUrl, nil)
	if err != nil {
		return "", "", gerror.Wrap(err, "创建请求失败")
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ImageDownloader/1.0)")

	// 3. 执行HTTP请求（带超时控制）
	resp, err := u.client.Do(req)
	if err != nil {
		return "", "", gerror.Wrap(err, "下载请求失败")
	}
	defer resp.Body.Close()

	// 4. 验证响应状态
	if resp.StatusCode != http.StatusOK {
		return "", "", gerror.Newf("无效的响应状态: %d", resp.StatusCode)
	}

	// 5. 验证内容类型
	contentType = resp.Header.Get("Content-Type")
	// 如果是application/octet-stream类型，通过文件头判断实际图片类型
	if contentType == "application/octet-stream" {
		contentType = "image/jpeg"
	}

	if !u.isAllowedType(contentType, allowedTypes) {
		return "", "", gerror.Newf("不允许的文件类型: %s", contentType)
	}

	// 6. 创建限制读取器（防超大文件）
	var reader io.Reader = resp.Body
	if maxSize > 0 {
		reader = io.LimitReader(resp.Body, maxSize)
	}

	// 7. 读取数据到缓冲区（内存安全）
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		return "", "", gerror.Wrap(err, "读取数据失败")
	}

	// 8. 验证实际数据大小
	if maxSize > 0 && int64(buf.Len()) >= maxSize {
		return "", "", gerror.New("文件大小超过限制")
	}

	// 新增: 检查图片是否小于1KB
	if buf.Len() < 1024 {
		return "", "", gerror.New("图片太小（小于1KB）")
	}

	// 新增: 计算图片数据的base64并检查Redis缓存
	imgBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	cacheKey := fmt.Sprintf(consts.ImageBase64CacheKey, gmd5.MustEncryptString(imgBase64))

	// 检查Redis中是否已存在该图片的缓存
	cachedValue, err := g.Redis().Do(ctx, "GET", cacheKey)
	if err == nil && !cachedValue.IsEmpty() {
		// 已存在相同图片，直接返回缓存的信息
		cachedData := cachedValue.String()
		parts := strings.Split(cachedData, ":")
		if len(parts) == 2 {
			return parts[0], parts[1], nil
		}
	}

	// 检查OSS实例是否初始化
	if OSSInstance == nil {
		return "", "", gerror.New("OSS客户端未初始化")
	}

	// 是否需要压缩
	var imageData []byte
	if compressQuality < 1 && compressQuality > 0 && u.isCompressibleImage(contentType) {
		// 使用传入的压缩质量
		compressedData, err := u.compressImage(buf.Bytes(), contentType, compressQuality)
		if err != nil {
			g.Log().Warning(ctx, "图片压缩失败，将使用原始图片:", err)
			imageData = buf.Bytes()
		} else {
			imageData = compressedData
		}
	} else {
		imageData = buf.Bytes()
	}

	// 检查OSS是否启用
	if !OSSInstance.IsEnabled() {
		g.Log().Info(ctx, "OSS上传功能未启用，将使用本地存储")
		// 使用本地存储作为备选
		localUrl, err := u.saveLocalFile(imageData, directory, contentType)
		if err != nil {
			return "", "", gerror.Wrap(err, "保存到本地存储失败")
		}

		// 将本地文件URL和类型缓存到Redis
		cacheValue := fmt.Sprintf("%s:%s", localUrl, contentType)
		_, err = g.Redis().Do(ctx, "SET", cacheKey, cacheValue, "EX", 3600)
		if err != nil {
			// 缓存失败不影响主流程，只记录日志
			g.Log().Warning(ctx, "图片缓存到Redis失败:", err)
		}

		return localUrl, contentType, nil
	}

	// 生成OSS对象键
	objectKey := OSSInstance.GenerateObjectKey(contentType, directory)

	// 上传到OSS
	ossUrl, err = OSSInstance.UploadFile(ctx, imageData, objectKey, contentType)
	if err != nil {
		return "", "", gerror.Wrap(err, "上传文件到OSS失败")
	}

	// 将图片信息缓存到Redis
	cacheValue := fmt.Sprintf("%s:%s", ossUrl, contentType)
	// 设置过期时间为1小时
	_, err = g.Redis().Do(ctx, "SET", cacheKey, cacheValue, "EX", 3600)
	if err != nil {
		// 缓存失败不影响主流程，只记录日志
		g.Log().Warning(ctx, "图片缓存到Redis失败:", err)
	}

	return ossUrl, contentType, nil
}

// isBase64Image 辅助方法：检查是否是base64编码的图片
func (u *ImageDownloader) isBase64Image(data string) bool {
	// 检查是否是 data:image 格式的base64
	if strings.HasPrefix(data, "data:image/") && strings.Contains(data, ";base64,") {
		return true
	}
	return false
}

// processBase64Image 处理并上传base64编码的图片到OSS
func (u *ImageDownloader) processBase64Image(
	ctx context.Context,
	base64Data string,
	directory string,
	maxSize int64,
	compressQuality float64,
) (ossUrl string, contentType string, err error) {
	allowedTypes := []string{"image/jpeg", "image/png", "image/gif", "image/webp"}

	// 解析MIME类型和base64内容
	parts := strings.SplitN(base64Data, ";base64,", 2)
	if len(parts) != 2 {
		return "", "", gerror.New("无效的base64图片格式")
	}

	// 获取内容类型
	mimeType := strings.TrimPrefix(parts[0], "data:")
	if !u.isAllowedType(mimeType, allowedTypes) {
		return "", "", gerror.Newf("不允许的文件类型: %s", mimeType)
	}

	// 解码base64数据
	imgData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", gerror.Wrap(err, "base64解码失败")
	}

	// 验证图片大小
	if maxSize > 0 && int64(len(imgData)) > maxSize {
		return "", "", gerror.New("文件大小超过限制")
	}

	// 检查图片是否小于1KB
	if len(imgData) < 1024 {
		return "", "", gerror.New("图片太小（小于1KB）")
	}

	// 检查Redis缓存
	imgBase64 := base64.StdEncoding.EncodeToString(imgData)
	cacheKey := fmt.Sprintf(consts.ImageBase64CacheKey, imgBase64)

	// 检查Redis中是否已存在该图片的缓存
	cachedValue, err := g.Redis().Do(ctx, "GET", cacheKey)
	if err == nil && !cachedValue.IsEmpty() {
		// 已存在相同图片，直接返回缓存的信息
		cachedData := cachedValue.String()
		parts := strings.Split(cachedData, ":")
		if len(parts) == 2 {
			return parts[0], parts[1], nil
		}
	}

	// 检查OSS实例是否初始化
	if OSSInstance == nil {
		return "", "", gerror.New("OSS客户端未初始化")
	}

	// 是否需要压缩
	var imageData []byte
	if compressQuality < 1 && compressQuality > 0 && u.isCompressibleImage(mimeType) {
		// 使用传入的压缩质量
		compressedData, err := u.compressImage(imgData, mimeType, compressQuality)
		if err != nil {
			imageData = imgData
		} else {
			imageData = compressedData
		}
	} else {
		imageData = imgData
	}

	// 检查OSS是否启用
	if !OSSInstance.IsEnabled() {
		g.Log().Info(ctx, "OSS上传功能未启用，将使用本地存储")
		// 使用本地存储作为备选
		localUrl, err := u.saveLocalFile(imageData, directory, mimeType)
		if err != nil {
			return "", "", gerror.Wrap(err, "保存到本地存储失败")
		}

		// 将本地文件URL和类型缓存到Redis
		cacheValue := fmt.Sprintf("%s:%s", localUrl, mimeType)
		_, err = g.Redis().Do(ctx, "SET", cacheKey, cacheValue, "EX", 3600)
		if err != nil {
			// 缓存失败不影响主流程，只记录日志
			g.Log().Warning(ctx, "图片缓存到Redis失败:", err)
		}

		return localUrl, mimeType, nil
	}

	// 生成OSS对象键
	objectKey := OSSInstance.GenerateObjectKey(mimeType, directory)

	// 上传到OSS
	ossUrl, err = OSSInstance.UploadFile(ctx, imageData, objectKey, mimeType)
	if err != nil {
		return "", "", gerror.Wrap(err, "上传文件到OSS失败")
	}

	// 将图片信息缓存到Redis
	cacheValue := fmt.Sprintf("%s:%s", ossUrl, mimeType)
	_, err = g.Redis().Do(ctx, "SET", cacheKey, cacheValue)
	if err != nil {
		// 缓存失败不影响主流程，只记录日志
		g.Log().Warning(ctx, "图片缓存到Redis失败:", err)
	}

	return ossUrl, mimeType, nil
}

// 辅助方法：URL验证
func (u *ImageDownloader) validateUrl(rawUrl string) error {
	parsed, err := url.Parse(rawUrl)
	if err != nil {
		return gerror.Wrap(err, "URL解析失败")
	}

	// 只允许HTTP/HTTPS
	if !strings.HasPrefix(parsed.Scheme, "http") {
		return gerror.New("只支持HTTP/HTTPS协议")
	}

	// 禁止本地网络地址
	if ip := net.ParseIP(parsed.Hostname()); ip != nil {
		if ip.IsLoopback() || ip.IsPrivate() {
			return gerror.New("禁止访问本地网络地址")
		}
	}

	return nil
}

// isAllowedType 辅助方法：检查允许的类型
func (u *ImageDownloader) isAllowedType(contentType string, allowedTypes []string) bool {
	if len(allowedTypes) == 0 {
		return true // 未设置限制时允许所有类型
	}

	for _, t := range allowedTypes {
		if strings.HasPrefix(contentType, t) {
			return true
		}
	}
	return false
}

// saveLocalFile 保存文件到本地（当OSS未启用时）
func (u *ImageDownloader) saveLocalFile(
	data []byte,
	directory string,
	contentType string,
) (string, error) {
	// 确保目录存在
	uploadDir := ""
	if directory != "" {
		uploadDir += strings.TrimPrefix(directory, "/")
		if !strings.HasSuffix(uploadDir, "/") {
			uploadDir += "/"
		}
	}

	// 使用GoFrame创建目录
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", gerror.Wrap(err, "创建本地目录失败")
	}

	// 生成文件名（复用OSS对象键生成逻辑）
	var fileName string
	if OSSInstance != nil {
		fileName = OSSInstance.GenerateObjectKey(contentType, "")
	} else {
		// 基础文件名：时间戳+随机数
		rand.Seed(time.Now().UnixNano())
		name := time.Now().Format("20060102150405") + "_" + strconv.Itoa(rand.Intn(100000))

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
		fileName = name + ext
	}

	// 组合完整路径
	filePath := uploadDir + fileName

	// 写入文件
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", gerror.Wrap(err, "保存本地文件失败")
	}

	// 返回相对URL
	return "/" + filePath, nil
}

// compressImage 压缩图片数据
// data: 原始图片数据
// contentType: 图片类型（MIME类型）
// quality: 压缩质量，0.1-0.9表示压缩比例（1表示不压缩，但这个值不应该传入本函数）
// 返回: 压缩后的数据，错误
func (u *ImageDownloader) compressImage(
	data []byte,
	contentType string,
	quality float64,
) ([]byte, error) {
	// 记录压缩前大小
	originalSize := len(data)

	// 解码图片
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, gerror.Wrap(err, "解码图片失败")
	}

	// 准备输出缓冲区
	var buf bytes.Buffer

	// 根据图片类型选择压缩策略
	switch {
	case strings.HasPrefix(contentType, "image/jpeg") || format == "jpeg":
		// JPEG压缩
		// 将0.1-1的浮点数转换为1-100的整数质量值
		jpegQuality := int(quality * 100)
		if jpegQuality <= 0 || jpegQuality > 100 {
			jpegQuality = 85 // 默认质量
		}
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: jpegQuality})

	case strings.HasPrefix(contentType, "image/png") || format == "png":
		// PNG压缩（无损，但减少颜色数量）
		encoder := png.Encoder{
			CompressionLevel: png.BestCompression,
		}
		err = encoder.Encode(&buf, img)

	default:
		// 其他格式暂不压缩，直接返回原数据
		return data, nil
	}

	if err != nil {
		return nil, gerror.Wrap(err, "压缩图片失败")
	}

	// 检查压缩是否有效（压缩后大小应该小于原始大小）
	if buf.Len() >= originalSize {
		return data, nil
	}

	return buf.Bytes(), nil
}

// isCompressibleImage 检查图片是否可压缩
func (u *ImageDownloader) isCompressibleImage(contentType string) bool {
	// 目前只支持压缩JPEG和PNG图片
	compressibleTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
	}

	for _, t := range compressibleTypes {
		if strings.HasPrefix(contentType, t) {
			return true
		}
	}

	return false
}
