package model

import "time"

type Workspace struct {
	ID          uint64    `gorm:"primarykey;autoIncrement" json:"id"` // 工作区ID
	Name        string    `gorm:"type:varchar(128);not null" json:"name"` // 工作区名称
	Description string    `gorm:"type:text" json:"description"`       // 工作区描述
	OwnerID     uint64    `gorm:"not null;index" json:"owner_id"`     // 工作区所有者ID
	CreatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
	UpdatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"` // 更新时间
}

func (Workspace) TableName() string {
	return "workspaces"
}