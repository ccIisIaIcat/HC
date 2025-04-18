package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"strconv"
	"strings"
)

// CalculateMD5 计算文件的MD5值
func CalculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// CompareVersions 比较两个版本号，返回：
// -1 如果 v1 < v2
// 0 如果 v1 = v2
// 1 如果 v1 > v2
func CompareVersions(v1, v2 string) int {
	v1Parts := strings.Split(v1, ".")
	v2Parts := strings.Split(v2, ".")

	// 确保两个版本号都有相同数量的部分
	maxLen := len(v1Parts)
	if len(v2Parts) > maxLen {
		maxLen = len(v2Parts)
	}

	for i := 0; i < maxLen; i++ {
		var v1Num, v2Num int
		var err error

		if i < len(v1Parts) {
			v1Num, err = strconv.Atoi(v1Parts[i])
			if err != nil {
				return -1
			}
		}

		if i < len(v2Parts) {
			v2Num, err = strconv.Atoi(v2Parts[i])
			if err != nil {
				return 1
			}
		}

		if v1Num < v2Num {
			return -1
		}
		if v1Num > v2Num {
			return 1
		}
	}

	return 0
}
