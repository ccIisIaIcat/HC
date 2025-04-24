package achievement

import (
	"backend/models"
	"fmt"
	"time"
)

// GiveItemsToUser 给用户添加指定的成就物品
func (c *Check) GiveItemsToUser(achievedItemIDs []uint) ([]models.UserItem, error) {
	if len(achievedItemIDs) == 0 {
		fmt.Println("没有新成就物品需要添加")
		return nil, nil
	}

	// 确保 Check 结构体中有用户ID
	if c.UserID == 0 {
		return nil, fmt.Errorf("未设置用户ID")
	}

	// 为用户添加新达成的成就物品
	var newUserItems []models.UserItem

	for _, itemID := range achievedItemIDs {
		// 创建新的用户物品记录
		userItem := models.UserItem{
			UserID:       c.UserID,
			ItemID:       itemID,
			Quantity:     1,
			ObtainedAt:   time.Now(),
			ObtainedFrom: "achievement", // 标记来源为成就系统
		}

		// 保存到数据库
		err := models.CreateUserItem(&userItem)
		if err != nil {
			return newUserItems, fmt.Errorf("添加物品(ID: %d)失败: %v", itemID, err)
		}

		// 找到对应的物品信息
		var itemName string
		for _, item := range c.Items {
			if item.ID == itemID {
				itemName = item.Name
				break
			}
		}

		fmt.Printf("成功给用户 %d 添加成就物品: %s (ID: %d)\n", c.UserID, itemName, itemID)
		newUserItems = append(newUserItems, userItem)
	}

	return newUserItems, nil
}
