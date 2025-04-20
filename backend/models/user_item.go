package models

import (
	"time"

	"gorm.io/gorm"
)

// UserItem 用户物品关联模型
type UserItem struct {
	gorm.Model             // 包含 ID、CreatedAt、UpdatedAt、DeletedAt
	UserID       uint      `json:"user_id" gorm:"index;not null"`      // 用户ID
	ItemID       uint      `json:"item_id" gorm:"index;not null"`      // 物品ID
	Quantity     int       `json:"quantity" gorm:"not null;default:1"` // 数量
	ObtainedAt   time.Time `json:"obtained_at" gorm:"not null"`        // 获得时间
	ObtainedFrom string    `json:"obtained_from" gorm:"size:100"`      // 具体获得方式
	Item         Item      `json:"item" gorm:"foreignKey:ItemID"`      // 关联物品信息
}

// CreateUserItem 创建新的用户物品记录
func CreateUserItem(userItem *UserItem) error {
	return DB.Create(userItem).Error
}

// GetUserItems 获取用户的所有物品（根据物品和获得来源聚合数量）
func GetUserItems(userID uint) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := DB.Table("user_items").
		Select(`
			items.id, 
			items.name, 
			items.description, 
			items.source, 
			items.icon_url, 
			items.image_url,
			items.created_at, 
			items.updated_at,
			user_items.obtained_from,
			SUM(user_items.quantity) as total_quantity
		`).
		Joins("LEFT JOIN items ON items.id = user_items.item_id").
		Where("user_items.user_id = ? AND user_items.deleted_at IS NULL", userID).
		Group("items.id, user_items.obtained_from").
		Order("items.name, user_items.obtained_from").
		Scan(&results).Error

	return results, err
}

// GetUserItemDetails 获取用户特定物品的详细信息
func GetUserItemDetails(userID, itemID uint) (*UserItem, error) {
	var userItem UserItem
	err := DB.Preload("Item").
		Where("user_id = ? AND item_id = ?", userID, itemID).
		First(&userItem).Error
	if err != nil {
		return nil, err
	}
	return &userItem, nil
}

// UpdateUserItemQuantity 更新用户物品数量
func UpdateUserItemQuantity(userID, itemID uint, quantity int) error {
	return DB.Model(&UserItem{}).
		Where("user_id = ? AND item_id = ?", userID, itemID).
		Update("quantity", quantity).Error
}

// DeleteUserItem 删除用户物品
func DeleteUserItem(userID, itemID uint) error {
	return DB.Where("user_id = ? AND item_id = ?", userID, itemID).
		Delete(&UserItem{}).Error
}

// HasUserItem 检查用户是否拥有特定物品
func HasUserItem(userID, itemID uint) (bool, error) {
	var count int64
	err := DB.Model(&UserItem{}).
		Where("user_id = ? AND item_id = ?", userID, itemID).
		Count(&count).Error
	return count > 0, err
}

// GetUserItemsBySource 获取用户特定来源的物品
func GetUserItemsBySource(userID uint, source string) ([]UserItem, error) {
	var userItems []UserItem
	err := DB.Preload("Item").
		Joins("LEFT JOIN items ON items.id = user_items.item_id").
		Where("user_items.user_id = ? AND items.source = ?", userID, source).
		Find(&userItems).Error
	return userItems, err
}

// GetUserItemCount 获取用户物品总数
func GetUserItemCount(userID uint) (int64, error) {
	var count int64
	err := DB.Model(&UserItem{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}
