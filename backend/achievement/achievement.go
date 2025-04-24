package achievement

import (
	"sort"
	"time"
)

// CheckFirstFoodRecord 检查是否满足"饭团猫的礼物"成就
// 条件：记录第一次食物
func (c *Check) CheckFirstFoodRecord() bool {
	return len(c.FoodRecords) > 0
}

// CheckSevenDaysCheckin 检查是否满足"第七天的饭勺"成就
// 条件：累计打卡7天
func (c *Check) CheckSevenDaysCheckin() bool {
	return len(c.CheckIns) >= 7
}

// CheckTwentyOneDaysCheckin 检查是否满足"超人披风"成就
// 条件：累计打卡21天
func (c *Check) CheckTwentyOneDaysCheckin() bool {
	return len(c.CheckIns) >= 21
}

// CheckReboundCheckin 检查是否满足"不倒翁餐盘"成就
// 条件：在打卡中断后重新打卡3天
func (c *Check) CheckReboundCheckin() bool {
	if len(c.CheckIns) < 3 {
		return false
	}

	// 创建一个日期到打卡记录的映射
	dateToCheckin := make(map[string]bool)
	for _, checkin := range c.CheckIns {
		// 使用日期作为键（忽略时间部分）
		dateStr := checkin.CheckInAt.Format("2006-01-02")
		dateToCheckin[dateStr] = true
	}

	// 获取日期列表并排序
	var dates []string
	for dateStr := range dateToCheckin {
		dates = append(dates, dateStr)
	}
	sort.Strings(dates) // 按日期从早到晚排序

	// 查找中断后连续3天的模式
	consecutiveDays := 0
	hasBreak := false

	for i := 1; i < len(dates); i++ {
		// 解析当前日期和前一天日期
		currentDate, _ := time.Parse("2006-01-02", dates[i])
		prevDate, _ := time.Parse("2006-01-02", dates[i-1])

		// 计算日期差
		daysDiff := int(currentDate.Sub(prevDate).Hours() / 24)

		if daysDiff == 1 {
			// 连续的日期
			consecutiveDays++
			// 如果之前有中断且现在已经连续3天，满足条件
			if hasBreak && consecutiveDays >= 3 {
				return true
			}
		} else if daysDiff > 1 {
			// 发现中断
			hasBreak = true
			// 重置连续天数计数（当前天算第一天）
			consecutiveDays = 1
		}
	}

	return false
}

// CheckEarlyBreakfast 检查是否满足"太阳花"成就
// 条件：记录一次早餐，且时间在上午9点前
func (c *Check) CheckEarlyBreakfast() bool {
	for _, record := range c.FoodRecords {
		if record.MealType == "早餐" {
			hour := record.RecordTime.Hour()
			if hour < 9 {
				return true
			}
		}
	}
	return false
}

// CheckHighProtein 检查是否满足"奇怪的鸡蛋"成就
// 条件：记录一次高蛋白（超过20克）的餐食
func (c *Check) CheckHighProtein() bool {
	for _, record := range c.FoodRecords {
		if record.Protein >= 20 {
			return true
		}
	}
	return false
}

// CheckBalancedMeal 检查是否满足"大小适中的碗"成就
// 条件：记录一次热量适中（400-600卡）且三大营养素比例均衡的正餐
func (c *Check) CheckBalancedMeal() bool {
	for _, record := range c.FoodRecords {
		// 检查热量是否在400-600卡范围内
		if record.Calories >= 400 && record.Calories <= 600 {
			// 计算三大营养素的总量
			total := record.Protein + record.TotalFat + record.Carbohydrates
			if total == 0 {
				continue
			}

			// 检查各营养素比例是否均衡
			// 简化判断：蛋白质10-35%，脂肪20-35%，碳水45-65%
			proteinRatio := record.Protein / total
			fatRatio := record.TotalFat / total
			carbRatio := record.Carbohydrates / total

			if proteinRatio >= 0.1 && proteinRatio <= 0.35 &&
				fatRatio >= 0.2 && fatRatio <= 0.35 &&
				carbRatio >= 0.45 && carbRatio <= 0.65 {
				return true
			}
		}
	}
	return false
}

// CheckLowCarb 检查是否满足"几个糖分子"成就
// 条件：记录一次低糖（碳水化合物少于30克）的餐食
func (c *Check) CheckLowCarb() bool {
	for _, record := range c.FoodRecords {
		if record.Carbohydrates < 30 {
			return true
		}
	}
	return false
}

func (c *Check) CheckConsistentRecording() bool {
	if len(c.FoodRecords) < 6 { // 至少需要6条记录才可能有3天每天2餐
		return false
	}

	// 按日期分组记录
	dailyRecords := make(map[string]int)
	for _, record := range c.FoodRecords {
		date := record.RecordTime.Format("2006-01-02")
		dailyRecords[date]++
	}

	// 将日期提取出来并排序
	var dates []string
	for date := range dailyRecords {
		dates = append(dates, date)
	}
	sort.Strings(dates) // 按日期从早到晚排序

	// 检查是否有连续3天每天至少2条记录
	if len(dates) < 3 {
		return false
	}

	consecutiveDays := 0
	for i := 0; i < len(dates); i++ {
		// 当前日期有至少2条记录
		if dailyRecords[dates[i]] >= 2 {
			// 如果是第一天或者与前一天相邻
			if i == 0 || isConsecutiveDate(dates[i-1], dates[i]) {
				consecutiveDays++
				if consecutiveDays >= 3 {
					return true
				}
			} else {
				// 日期不连续，重置计数
				consecutiveDays = 1
			}
		} else {
			// 记录不足2条，重置计数
			consecutiveDays = 0
		}
	}
	return false
}

// isConsecutiveDate 检查两个日期是否相邻
func isConsecutiveDate(date1, date2 string) bool {
	t1, err := time.Parse("2006-01-02", date1)
	if err != nil {
		return false
	}
	t2, err := time.Parse("2006-01-02", date2)
	if err != nil {
		return false
	}

	// 计算两个日期之间的差距（天数）
	diff := int(t2.Sub(t1).Hours() / 24)

	// 如果差距为1天，则是连续的
	return diff == 1
}

// CheckBalancedWeeklyMeals 检查是否满足"正常的天秤"成就
// 条件：一周内有3次餐食的三大营养素比例符合推荐标准
func (c *Check) CheckBalancedWeeklyMeals() bool {
	// 获取最近一周的记录
	oneWeekAgo := time.Now().AddDate(0, 0, -7)

	balancedMeals := 0
	for _, record := range c.FoodRecords {
		if record.RecordTime.After(oneWeekAgo) {
			// 计算三大营养素的总量
			total := record.Protein + record.TotalFat + record.Carbohydrates
			if total == 0 {
				continue
			}

			// 检查各营养素比例是否均衡
			proteinRatio := record.Protein / total
			fatRatio := record.TotalFat / total
			carbRatio := record.Carbohydrates / total

			if proteinRatio >= 0.1 && proteinRatio <= 0.35 &&
				fatRatio >= 0.2 && fatRatio <= 0.35 &&
				carbRatio >= 0.45 && carbRatio <= 0.65 {
				balancedMeals++
				if balancedMeals >= 3 {
					return true
				}
			}
		}
	}
	return false
}

// CheckWeeklyBreakfast 检查是否满足"啄木鸟闹钟"成就
// 条件：一周内记录至少4次早餐
func (c *Check) CheckWeeklyBreakfast() bool {
	// 获取最近一周的记录
	oneWeekAgo := time.Now().AddDate(0, 0, -7)

	breakfastCount := 0
	for _, record := range c.FoodRecords {
		if record.RecordTime.After(oneWeekAgo) && record.MealType == "早餐" {
			breakfastCount++
			if breakfastCount >= 4 {
				return true
			}
		}
	}
	return false
}

// CheckNoLateNightEating 检查是否满足"安眠小夜灯"成就
// 条件：连续7天没有记录晚上9点后的进食
func (c *Check) CheckNoLateNightEating() bool {
	// 获取最近7天的日期
	today := time.Now()
	dates := make([]time.Time, 7)
	for i := 0; i < 7; i++ {
		dates[i] = today.AddDate(0, 0, -i)
	}

	// 检查每一天是否有晚上9点后的进食记录
	for _, date := range dates {
		year, month, day := date.Date()
		nightStart := time.Date(year, month, day, 21, 0, 0, 0, time.Local)
		nextDayStart := time.Date(year, month, day+1, 0, 0, 0, 0, time.Local)

		for _, record := range c.FoodRecords {
			recordTime := record.RecordTime
			if recordTime.After(nightStart) && recordTime.Before(nextDayStart) {
				return false
			}
		}
	}
	return true
}

// CheckDiverseFoods 检查是否满足"美食探险背包"成就
// 条件：记录20种从未记录过的新食物
func (c *Check) CheckDiverseFoods() bool {
	// 统计不同食物的种类
	foodTypes := make(map[string]bool)
	for _, record := range c.FoodRecords {
		foodTypes[record.FoodName] = true
	}

	return len(foodTypes) >= 20
}
