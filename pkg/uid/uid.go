package uid

import (
	"errors"

	"github.com/google/uuid"
)

// NewUUID 生成一个新的 UUID
func NewUUID() string {
	return uuid.New().String()
}

// ParseUUID 解析 UUID 字符串
func ParseUUID(u string) (uuid.UUID, error) {
	parsedUUID, err := uuid.Parse(u)
	if err != nil {
		return uuid.Nil, errors.New("无效的 UUID 格式")
	}
	return parsedUUID, nil
}

// IsValidUUID 检查字符串是否是有效的 UUID
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// MustParseUUID 解析 UUID 字符串，如果解析失败则 panic
func MustParseUUID(u string) uuid.UUID {
	parsedUUID, err := ParseUUID(u)
	if err != nil {
		panic(err)
	}
	return parsedUUID
}

// UUIDVersion 返回 UUID 的版本号
func UUIDVersion(u string) (int, error) {
	parsedUUID, err := ParseUUID(u)
	if err != nil {
		return 0, err
	}
	return int(parsedUUID.Version()), nil
}
