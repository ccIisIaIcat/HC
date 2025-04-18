package models

import (
	"time"

	"gorm.io/gorm"
)

// AppUpdate 应用程序更新信息
type AppUpdate struct {
	gorm.Model
	Version     string    `json:"version" gorm:"type:varchar(50);not null"`
	Platform    string    `json:"platform" gorm:"type:varchar(20);not null"`
	UpdateURL   string    `json:"update_url" gorm:"type:varchar(255);not null"`
	ForceUpdate bool      `json:"force_update" gorm:"default:false"`
	UpdateNotes string    `json:"update_notes" gorm:"type:text"`
	ReleaseDate time.Time `json:"release_date" gorm:"not null"`
	FileSize    int64     `json:"file_size" gorm:"not null"`
	MD5         string    `json:"md5" gorm:"type:varchar(32);not null"`
}
