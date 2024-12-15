package model

import "time"

type Folder struct {
	ID           uint64    `gorm:"primarykey" json:"id"`
	CollectionID uint64    `gorm:"not null;index" json:"collection_id"` // 关联到集合
	Name         string    `gorm:"type:varchar(64);not null" json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Folder) TableName() string {
	return "folders"
}

type FolderClosure struct {
	Ancestor   uint64 `gorm:"not null;index"` // 祖先文件夹的 ID
	Descendant uint64 `gorm:"not null;index"` // 后代文件夹的 ID
	Depth      int    `gorm:"not null"`       // 祖先与后代之间的距离
}

func (FolderClosure) TableName() string {
	return "folder_closures"
}
