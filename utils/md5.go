package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// MD5String 计算字符串的 MD5 哈希值
func MD5String(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// MD5Bytes 计算字节数组的 MD5 哈希值
func MD5Bytes(data []byte) string {
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

// MD5File 计算文件的 MD5 哈希值
func MD5File(filePath string) (string, error) {
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
