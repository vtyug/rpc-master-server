package model

import "time"

type Collections struct {
	ID           uint64    `gorm:"primarykey;autoIncrement" json:"id"`                                                     //
	WorkspaceID  uint64    `gorm:"not null;index" json:"workspace_id"`                                                     // 工作空间ID
	Name         string    `gorm:"type:varchar(128);not null" json:"name"`                                                 // 收藏夹名称
	CollectionID string    `gorm:"type:varchar(128);not null;index" json:"collection_id"`                                  // 收藏夹ID
	CreatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`                             // 创建时间
	UpdatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"` //
}

func (Collections) TableName() string {
	return "collections"
}
