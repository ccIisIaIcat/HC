package models

import (
	"gorm.io/gorm"
)

// Item 物品模型
type Item struct {
	gorm.Model             // 包含 ID、CreatedAt、UpdatedAt、DeletedAt
	Name        string     `json:"name" gorm:"size:100;not null"`                 // 物品名称
	Description string     `json:"description" gorm:"type:text"`                  // 物品介绍
	Source      string     `json:"source" gorm:"size:100;not null"`               // 获取来源
	IconURL     string     `json:"icon_url" gorm:"size:255;not null"`             // 图标路径
	ImageURL    string     `json:"image_url" gorm:"size:255;not null"`            // 图片路径
	UserItems   []UserItem `json:"user_items,omitempty" gorm:"foreignKey:ItemID"` // 关联用户物品
}

// CreateItem 创建新物品
func CreateItem(item *Item) error {
	return DB.Create(item).Error
}

// GetItemByID 根据ID获取物品
func GetItemByID(itemID uint) (*Item, error) {
	var item Item
	err := DB.First(&item, itemID).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetAllItems 获取所有物品（支持分页）
func GetAllItems(page, pageSize int) ([]Item, int64, error) {
	var items []Item
	var total int64

	// 获取总数
	err := DB.Model(&Item{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = DB.Offset(offset).Limit(pageSize).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// UpdateItem 更新物品信息
func UpdateItem(item *Item) error {
	return DB.Save(item).Error
}

// DeleteItem 删除物品
func DeleteItem(itemID uint) error {
	return DB.Delete(&Item{}, itemID).Error
}

// GetItemsByName 根据名称搜索物品
func GetItemsByName(name string) ([]Item, error) {
	var items []Item
	err := DB.Where("name LIKE ?", "%"+name+"%").Find(&items).Error
	return items, err
}

// GetItemsBySource 根据来源搜索物品
func GetItemsBySource(source string) ([]Item, error) {
	var items []Item
	err := DB.Where("source = ?", source).Find(&items).Error
	return items, err
}

// GetItemSources 获取所有物品来源（用于筛选）
func GetItemSources() ([]string, error) {
	var sources []string
	err := DB.Model(&Item{}).Distinct().Pluck("source", &sources).Error
	return sources, err
}

// BatchCreateItems 批量创建物品
func BatchCreateItems(items []Item) error {
	return DB.Create(&items).Error
}

// ExistsItemByName 检查指定名称的物品是否存在
func ExistsItemByName(name string) (bool, error) {
	var count int64
	err := DB.Model(&Item{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}
