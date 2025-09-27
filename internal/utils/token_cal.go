package utils

import (
	"math"
	"os"
)

// tokens计算
type TokenCalder struct {
}

var TokenCal = &TokenCalder{}

func (u *TokenCalder) CalImage(filePath string) (tokens int, err error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}

	sizeMB := float64(fileInfo.Size()) / (1024 * 1024)
	tokens = int(math.Ceil(sizeMB * 512)) // 假设每MB约512个Token
	return tokens, nil
}

func (u *TokenCalder) CalText(text string) (tokens int, err error) {
	// 英文：1 token ≈ 4字符
	// 中文：1 token ≈ 1.5字符
	charCount := len([]rune(text))
	if isChinese(text) {
		return int(float64(charCount) * 1.5), nil
	}
	return (charCount + 3) / 4, nil
}

func isChinese(text string) bool {
	for _, r := range text {
		if r >= '\u4e00' && r <= '\u9fff' {
			return true
		}
	}
	return false
}
