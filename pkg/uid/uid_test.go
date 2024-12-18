// 文件路径: pkg/uid/uid_test.go
package uid

import (
	"testing"
)

func TestNewUUID(t *testing.T) {
	uuid := NewUUID()
	if !IsValidUUID(uuid) {
		t.Errorf("生成的 UUID 无效: %s", uuid)
	}
	t.Logf("生成的 UUID: %s", uuid)
}

func TestParseUUID(t *testing.T) {
	validUUID := NewUUID()
	parsedUUID, err := ParseUUID(validUUID)
	if err != nil {
		t.Errorf("解析有效 UUID 失败: %v", err)
	}
	if parsedUUID.String() != validUUID {
		t.Errorf("解析结果不匹配: got %s, want %s", parsedUUID.String(), validUUID)
	}

	invalidUUID := "invalid-uuid"
	_, err = ParseUUID(invalidUUID)
	if err == nil {
		t.Error("解析无效 UUID 应该失败")
	}
}

func TestIsValidUUID(t *testing.T) {
	validUUID := NewUUID()
	if !IsValidUUID(validUUID) {
		t.Errorf("有效 UUID 被识别为无效: %s", validUUID)
	}

	invalidUUID := "invalid-uuid"
	if IsValidUUID(invalidUUID) {
		t.Error("无效 UUID 被识别为有效")
	}
}

func TestMustParseUUID(t *testing.T) {
	validUUID := NewUUID()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MustParseUUID 应该不会 panic 对于有效 UUID: %s", validUUID)
		}
	}()
	MustParseUUID(validUUID)

	invalidUUID := "invalid-uuid"
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustParseUUID 应该 panic 对于无效 UUID")
		}
	}()
	MustParseUUID(invalidUUID)
}

func TestUUIDVersion(t *testing.T) {
	validUUID := NewUUID()
	version, err := UUIDVersion(validUUID)
	if err != nil {
		t.Errorf("获取 UUID 版本失败: %v", err)
	}
	if version != 4 {
		t.Errorf("UUID 版本不正确: got %d, want 4", version)
	}

	invalidUUID := "invalid-uuid"
	_, err = UUIDVersion(invalidUUID)
	if err == nil {
		t.Error("获取无效 UUID 版本应该失败")
	}
}
