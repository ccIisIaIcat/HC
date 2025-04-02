package models

import (
	"time"

	"gorm.io/gorm"
)

// FoodRecord 表示用户的饮食记录
type FoodRecord struct {
	gorm.Model           // 包含ID、CreatedAt、UpdatedAt、DeletedAt
	UserID     uint      `json:"user_id" gorm:"index"`     // 用户ID，关联User表
	RecordTime time.Time `json:"record_time" gorm:"index"` // 记录时间
	FoodName   string    `json:"food_name"`                // 食物名称
	Weight     float64   `json:"weight"`                   // 食物重量(克)

	// 基本营养素
	Calories       float64 `json:"calories"`        // 热量（卡路里）
	Protein        float64 `json:"protein"`         // 蛋白质含量（克）
	TotalFat       float64 `json:"total_fat"`       // 总脂肪（克）
	SaturatedFat   float64 `json:"saturated_fat"`   // 饱和脂肪（克）
	TransFat       float64 `json:"trans_fat"`       // 反式脂肪（克）
	UnsaturatedFat float64 `json:"unsaturated_fat"` // 不饱和脂肪（克）
	Carbohydrates  float64 `json:"carbohydrates"`   // 碳水化合物（克）
	Sugar          float64 `json:"sugar"`           // 糖分（克）
	Fiber          float64 `json:"fiber"`           // 膳食纤维（克）

	// 维生素 (选择性存储主要维生素)
	VitaminA  float64 `json:"vitamin_a"`  // 维生素A（μg）
	VitaminC  float64 `json:"vitamin_c"`  // 维生素C（mg）
	VitaminD  float64 `json:"vitamin_d"`  // 维生素D（μg）
	VitaminB1 float64 `json:"vitamin_b1"` // 维生素B1（mg）
	VitaminB2 float64 `json:"vitamin_b2"` // 维生素B2（mg）

	// 矿物质 (选择性存储主要矿物质)
	Calcium   float64 `json:"calcium"`   // 钙（mg）
	Iron      float64 `json:"iron"`      // 铁（mg）
	Sodium    float64 `json:"sodium"`    // 钠（mg）
	Potassium float64 `json:"potassium"` // 钾（mg）

	// 记录类型和备注
	MealType  string `json:"meal_type"`  // 餐食类型：早餐/午餐/晚餐/加餐
	Notes     string `json:"notes"`      // 备注
	ImagePath string `json:"image_path"` // 图片路径（可选）
}

// 将食物分析结果转换为食物记录
func CreateFoodRecordFromAnalysis(userID uint, foodAnalysis *FoodAnalysis, mealType string, notes string) *FoodRecord {
	record := &FoodRecord{
		UserID:     userID,
		RecordTime: time.Now(),
		FoodName:   foodAnalysis.FoodType,
		Weight:     foodAnalysis.Weight,

		// 基本营养素
		Calories:       foodAnalysis.Nutrition.Calories,
		Protein:        foodAnalysis.Nutrition.Protein,
		TotalFat:       foodAnalysis.Nutrition.TotalFat,
		SaturatedFat:   foodAnalysis.Nutrition.SaturatedFat,
		TransFat:       foodAnalysis.Nutrition.TransFat,
		UnsaturatedFat: foodAnalysis.Nutrition.UnsaturatedFat,
		Carbohydrates:  foodAnalysis.Nutrition.Carbohydrates,
		Sugar:          foodAnalysis.Nutrition.Sugar,
		Fiber:          foodAnalysis.Nutrition.Fiber,

		// 维生素
		VitaminA:  foodAnalysis.Nutrition.Vitamins.VitaminA,
		VitaminC:  foodAnalysis.Nutrition.Vitamins.VitaminC,
		VitaminD:  foodAnalysis.Nutrition.Vitamins.VitaminD,
		VitaminB1: foodAnalysis.Nutrition.Vitamins.VitaminB.B1,
		VitaminB2: foodAnalysis.Nutrition.Vitamins.VitaminB.B2,

		// 矿物质
		Calcium:   foodAnalysis.Nutrition.Minerals.Calcium,
		Iron:      foodAnalysis.Nutrition.Minerals.Iron,
		Sodium:    foodAnalysis.Nutrition.Minerals.Sodium,
		Potassium: foodAnalysis.Nutrition.Minerals.Potassium,

		// 记录信息
		MealType: mealType,
		Notes:    notes,
	}

	return record
}

// 获取用户一段时间内的饮食记录
func GetUserFoodRecords(userID uint, startTime, endTime time.Time) ([]FoodRecord, error) {
	var records []FoodRecord

	result := DB.Where("user_id = ? AND record_time BETWEEN ? AND ?",
		userID, startTime, endTime).
		Order("record_time DESC").
		Find(&records)

	if result.Error != nil {
		return nil, result.Error
	}

	return records, nil
}

// 获取用户的所有饮食记录
func GetAllUserFoodRecords(userID uint) ([]FoodRecord, error) {
	var records []FoodRecord

	result := DB.Where("user_id = ?", userID).
		Order("record_time DESC").
		Find(&records)

	if result.Error != nil {
		return nil, result.Error
	}

	return records, nil
}

// 获取指定ID的食物记录
func GetFoodRecordByID(recordID uint) (*FoodRecord, error) {
	var record FoodRecord

	result := DB.First(&record, recordID)
	if result.Error != nil {
		return nil, result.Error
	}

	return &record, nil
}

// 创建食物记录
func CreateFoodRecord(record *FoodRecord) error {
	result := DB.Create(record)
	return result.Error
}

// 更新食物记录
func UpdateFoodRecord(record *FoodRecord) error {
	result := DB.Save(record)
	return result.Error
}

// 删除食物记录
func DeleteFoodRecord(recordID uint) error {
	result := DB.Delete(&FoodRecord{}, recordID)
	return result.Error
}
