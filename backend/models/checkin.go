package models

import (
	"time"

	"gorm.io/gorm"
)

// CheckIn 用户打卡记录模型
type CheckIn struct {
	gorm.Model
	UserID    uint      `json:"user_id" gorm:"index;not null"` // 用户ID，添加索引
	CheckInAt time.Time `json:"check_in_at"`                   // 打卡时间
	Content   string    `json:"content" gorm:"size:500"`       // 打卡内容
	ImageURL  string    `json:"image_url" gorm:"size:255"`     // 打卡图片URL（可选）
}

// 获取用户最近的打卡记录
func GetUserRecentCheckIns(userID uint, limit int) ([]CheckIn, error) {
	var checkIns []CheckIn
	result := DB.Where("user_id = ?", userID).
		Order("check_in_at DESC").
		Limit(limit).
		Find(&checkIns)
	return checkIns, result.Error
}

// 获取用户今日是否已打卡
func HasUserCheckedInToday(userID uint) (bool, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	result := DB.Model(&CheckIn{}).
		Where("user_id = ? AND DATE(check_in_at) = ?", userID, today).
		Count(&count)
	return count > 0, result.Error
}

// 创建打卡记录
func CreateCheckIn(checkIn *CheckIn) error {
	return DB.Create(checkIn).Error
}

// GetUserCheckIns 获取用户所有打卡记录
func GetUserCheckIns(userID uint) ([]CheckIn, error) {
	var checkIns []CheckIn
	result := DB.Where("user_id = ?", userID).
		Order("check_in_at DESC").
		Find(&checkIns)
	return checkIns, result.Error
}
