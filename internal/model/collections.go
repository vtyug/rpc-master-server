package model

import (
	"time"
)

// 定义请求类型的枚举
type CollectionType int

const (
	HTTP CollectionType = iota + 1
	GRPC
)

type Collections struct {
	ID           uint64         `gorm:"primarykey;autoIncrement" json:"id"`
	Name         string         `gorm:"type:varchar(128);not null" json:"name"`
	OwnerID      uint64         `gorm:"not null;index" json:"owner_id"`
	Protocol     CollectionType `gorm:"type:int(10);not null" json:"protocol"`
	WorkspaceID  uint64         `gorm:"not null;index" json:"workspace_id"`
	Description  string         `gorm:"type:text;not null" json:"description"`
	MembersCount int            `gorm:"type:int(10);not null;default:1" json:"members_count"`
	CollectionID string         `gorm:"type:varchar(128);not null;index" json:"collection_id"`
	CreatedAt    time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Collections) TableName() string {
	return "collections"
}

func ReturnString(c CollectionType) string {
	switch c {
	case HTTP:
		return "HTTP"
	case GRPC:
		return "gRPC"
	default:
		return "Unknown"
	}
}

func FromString(s string) CollectionType {
	switch s {
	case "HTTP":
		return HTTP
	case "gRPC":
		return GRPC
	default:
		return 0
	}
}
