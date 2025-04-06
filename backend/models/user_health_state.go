package models

import (
	"time"

	"gorm.io/gorm"
)

// UserHealthState 表示用户的健康状态
type UserHealthState struct {
	gorm.Model           // 包含ID、CreatedAt、UpdatedAt、DeletedAt
	UserID     uint      `json:"user_id" gorm:"index"`     // 用户ID，关联User表
	RecordTime time.Time `json:"record_time" gorm:"index"` // 记录时间

	// 基本生命体征

	// 身体测量指标
	Height              float64 `json:"height"`               // 身高，单位：cm
	Weight              float64 `json:"weight"`               // 体重，单位：kg
	BMI                 float64 `json:"bmi"`                  // 体重指数，单位：kg/m²
	BodyFatPercentage   float64 `json:"body_fat_percentage"`  // 体脂率，单位：%
	Temperature         float64 `json:"temperature"`          // 体温，单位：℃
	HeartRate           int     `json:"heart_rate"`           // 心率，单位：次/分钟
	RespiratoryRate     int     `json:"respiratory_rate"`     // 呼吸频率，单位：次/分钟
	FastingGlucose      float64 `json:"fasting_glucose"`      // 空腹血糖，单位：mmol/L
	PostprandialGlucose float64 `json:"postprandial_glucose"` // 餐后血糖，单位：mmol/L
	TotalCholesterol    float64 `json:"total_cholesterol"`    // 总胆固醇，单位：mmol/L

	// 备注
	Notes string `json:"notes"` // 备注信息
}

// 创建用户健康状态记录
func CreateUserHealthState(record *UserHealthState) error {
	result := DB.Create(record)
	return result.Error
}

// 获取用户最新的健康状态记录
func GetLatestUserHealthState(userID uint) (*UserHealthState, error) {
	var record UserHealthState

	result := DB.Where("user_id = ?", userID).
		Order("record_time DESC").
		First(&record)

	if result.Error != nil {
		return nil, result.Error
	}

	return &record, nil
}

// 获取用户一段时间内的健康状态记录
func GetUserHealthStates(userID uint, startTime, endTime time.Time) ([]UserHealthState, error) {
	var records []UserHealthState

	result := DB.Where("user_id = ? AND record_time BETWEEN ? AND ?",
		userID, startTime, endTime).
		Order("record_time DESC").
		Find(&records)

	if result.Error != nil {
		return nil, result.Error
	}

	return records, nil
}

// 根据ID获取用户健康状态记录
func GetUserHealthStateByID(recordID uint) (*UserHealthState, error) {
	var record UserHealthState

	result := DB.First(&record, recordID)

	if result.Error != nil {
		return nil, result.Error
	}

	return &record, nil
}

// 更新用户健康状态记录
func UpdateUserHealthState(record *UserHealthState) error {
	result := DB.Save(record)
	return result.Error
}

// 删除用户健康状态记录
func DeleteUserHealthState(recordID uint) error {
	result := DB.Delete(&UserHealthState{}, recordID)
	return result.Error
}
