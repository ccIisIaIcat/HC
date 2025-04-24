package achievement

import (
	"fmt"
	"strconv"
)

// AchievementCheckFunc 定义成就检查函数的类型
type AchievementCheckFunc func() bool

// GetAchievementChecks 返回物品ID到成就检查函数的映射
func (c *Check) GetAchievementChecks() map[uint]AchievementCheckFunc {
	// 创建物品名称到ID的映射，方便查找
	nameToID := make(map[string]uint)
	for _, item := range c.Items {
		nameToID[item.Name] = item.ID
	}

	// 创建ID到检查函数的映射
	checks := make(map[uint]AchievementCheckFunc)

	// 为每个成就物品添加对应的检查函数

	// 饭团猫的礼物
	if id, ok := nameToID["饭团猫的礼物"]; ok {
		checks[id] = func() bool {
			return c.CheckFirstFoodRecord()
		}
	}

	// 第七天的饭勺
	if id, ok := nameToID["第七天的饭勺"]; ok {
		checks[id] = func() bool {
			return c.CheckSevenDaysCheckin()
		}
	}

	// 超人披风
	if id, ok := nameToID["超人披风"]; ok {
		checks[id] = func() bool {
			return c.CheckTwentyOneDaysCheckin()
		}
	}

	// 不倒翁餐盘
	if id, ok := nameToID["不倒翁餐盘"]; ok {
		checks[id] = func() bool {
			return c.CheckReboundCheckin()
		}
	}

	// 太阳花
	if id, ok := nameToID["太阳花"]; ok {
		checks[id] = func() bool {
			return c.CheckEarlyBreakfast()
		}
	}

	// 奇怪的鸡蛋
	if id, ok := nameToID["奇怪的鸡蛋"]; ok {
		checks[id] = func() bool {
			return c.CheckHighProtein()
		}
	}

	// 大小适中的碗
	if id, ok := nameToID["大小适中的碗"]; ok {
		checks[id] = func() bool {
			return c.CheckBalancedMeal()
		}
	}

	// 几个糖分子
	if id, ok := nameToID["几个糖分子"]; ok {
		checks[id] = func() bool {
			return c.CheckLowCarb()
		}
	}

	// 美食记录本
	if id, ok := nameToID["美食记录本"]; ok {
		checks[id] = func() bool {
			return c.CheckConsistentRecording()
		}
	}

	// 正常的天秤
	if id, ok := nameToID["正常的天秤"]; ok {
		checks[id] = func() bool {
			return c.CheckBalancedWeeklyMeals()
		}
	}

	// 啄木鸟闹钟
	if id, ok := nameToID["啄木鸟闹钟"]; ok {
		checks[id] = func() bool {
			return c.CheckWeeklyBreakfast()
		}
	}

	// 安眠小夜灯
	if id, ok := nameToID["安眠小夜灯"]; ok {
		checks[id] = func() bool {
			return c.CheckNoLateNightEating()
		}
	}

	// 美食探险背包
	if id, ok := nameToID["美食探险背包"]; ok {
		checks[id] = func() bool {
			return c.CheckDiverseFoods()
		}
	}

	return checks
}

// CheckAchievements 检查用户的所有成就并返回已达成但未获得的成就ID列表
func (c *Check) CheckAchievements() []uint {
	achievementChecks := c.GetAchievementChecks()

	// 将用户已有的物品ID存入map，方便查找
	userItemIDs := make(map[uint]bool)
	for _, userItem := range c.UserItems {
		// 根据SQL查询，物品ID应该是"id"字段
		if itemID, ok := userItem["id"].(uint); ok {
			userItemIDs[itemID] = true
		} else if itemIDFloat, ok := userItem["id"].(float64); ok {
			// 处理JSON解析可能将uint转为float64的情况
			userItemIDs[uint(itemIDFloat)] = true
		} else if itemIDInt, ok := userItem["id"].(int); ok {
			// 处理可能是int的情况
			userItemIDs[uint(itemIDInt)] = true
		} else if itemIDInt64, ok := userItem["id"].(int64); ok {
			// 处理int64类型
			userItemIDs[uint(itemIDInt64)] = true
		} else if itemIDStr, ok := userItem["id"].(string); ok {
			// 处理可能是字符串的情况
			if id, err := strconv.ParseUint(itemIDStr, 10, 32); err == nil {
				userItemIDs[uint(id)] = true
			}
		}
	}

	// 存储符合条件但用户尚未获得的成就ID，使用map确保唯一性
	uniqueAchievedItemIDs := make(map[uint]bool)

	// 遍历所有成就检查函数
	for itemID, checkFunc := range achievementChecks {
		// 如果用户还没有这个成就物品，且满足获取条件
		if !userItemIDs[itemID] && checkFunc() {
			// 使用map确保唯一性
			uniqueAchievedItemIDs[itemID] = true
			fmt.Printf("用户达成了成就，物品ID: %d\n", itemID)
		}
	}

	// 转换回切片
	var achievedItemIDs []uint
	for itemID := range uniqueAchievedItemIDs {
		achievedItemIDs = append(achievedItemIDs, itemID)
	}

	return achievedItemIDs
}
