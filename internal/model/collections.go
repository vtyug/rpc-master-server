package model

import "time"

type Collections struct {
	ID          uint64    `gorm:"primarykey" json:"id"`               // 收藏夹ID
	WorkspaceID uint64    `gorm:"not null;index" json:"workspace_id"` // 团队ID
	Name        string    `gorm:"type:varchar(64);not null" json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Collections) TableName() string {
	return "collections"
}
