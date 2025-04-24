package achievement

import (
	"backend/models"
	"fmt"
	"math"
)

type Check struct {
	ID          int
	UserID      uint
	UserItems   []map[string]interface{}
	FoodRecords []models.FoodRecord
	CheckIns    []models.CheckIn
	Items       []models.Item
}

func CheckAchievements(UserID uint) {
	c := GenCheck(UserID)
	achievedItemIDs := c.CheckAchievements()
	_, err := c.GiveItemsToUser(achievedItemIDs)
	if err != nil {
		fmt.Printf("给用户 %d 添加成就物品失败: %v\n", UserID, err)
	}
}

func GenCheck(UserID uint) *Check {
	c := &Check{}
	c.UserID = UserID
	c.UserItems = make([]map[string]interface{}, 0)
	c.FoodRecords = make([]models.FoodRecord, 0)
	c.CheckIns = make([]models.CheckIn, 0)
	c.Items = make([]models.Item, 0)
	c.GatherData()
	c.CheckAchievements()
	return c
}

func (c *Check) GatherData() error {
	// 获取用户的打卡记录
	checkIns, err := models.GetUserCheckIns(c.UserID)
	if err != nil {
		return fmt.Errorf("获取用户打卡记录失败: %v", err)
	}
	c.CheckIns = checkIns
	fmt.Printf("用户 %d 的打卡记录已获取，共 %d 条记录\n", c.UserID, len(checkIns))

	// 获取用户的饮食记录
	foodRecords, err := models.GetAllUserFoodRecords(c.UserID)
	if err != nil {
		return fmt.Errorf("获取用户饮食记录失败: %v", err)
	}
	c.FoodRecords = foodRecords
	fmt.Printf("用户 %d 的饮食记录已获取，共 %d 条记录\n", c.UserID, len(foodRecords))

	// 获取用户的物品
	// 注意：这里假设有一个 GetUserItems 函数来获取用户的物品
	userItems, err := models.GetUserItems(c.UserID)
	if err != nil {
		return fmt.Errorf("获取用户物品失败: %v", err)
	}
	c.UserItems = userItems
	fmt.Printf("用户 %d 的物品已获取，共 %d 个物品\n", c.UserID, len(userItems))

	// 获取所有物品
	items, total, err := models.GetAllItems(1, 100) // 假设总物品数不超过100
	if err != nil {
		return fmt.Errorf("获取所有物品失败: %v", err)
	}

	// 如果物品数量超过了一页的限制，继续获取剩余物品
	if total > 100 {
		pageCount := int(math.Ceil(float64(total) / 100))
		for page := 2; page <= pageCount; page++ {
			moreItems, _, err := models.GetAllItems(page, 100)
			if err != nil {
				return fmt.Errorf("获取第 %d 页物品失败: %v", page, err)
			}
			items = append(items, moreItems...)
		}
	}

	c.Items = items
	fmt.Printf("所有物品已获取，共 %d 个物品\n", len(items))

	return nil
}
