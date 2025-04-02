package models

// FoodAnalysis 表示食物分析的结果
type FoodAnalysis struct {
	HasFood   bool               `json:"hasFood"`  // 是否检测到食物
	FoodType  string             `json:"foodType"` // 食物类型
	Weight    float64            `json:"weight"`   // 估计重量（克）
	Nutrition `json:"nutrition"` // 营养成分
}

// Nutrition 表示食物的营养成分
type Nutrition struct {
	Calories       float64  `json:"calories"`       // 热量（卡路里）
	Protein        float64  `json:"protein"`        // 蛋白质含量（克）
	TotalFat       float64  `json:"totalFat"`       // 总脂肪（克）
	SaturatedFat   float64  `json:"saturatedFat"`   // 饱和脂肪（克）
	TransFat       float64  `json:"transFat"`       // 反式脂肪（克）
	UnsaturatedFat float64  `json:"unsaturatedFat"` // 不饱和脂肪（克）
	Carbohydrates  float64  `json:"carbohydrates"`  // 碳水化合物（克）
	Sugar          float64  `json:"sugar"`          // 糖分（克）
	Fiber          float64  `json:"fiber"`          // 膳食纤维（克）
	Vitamins       Vitamins `json:"vitamins"`       // 维生素
	Minerals       Minerals `json:"minerals"`       // 矿物质
}

// Vitamins 表示维生素含量
type Vitamins struct {
	VitaminA float64 `json:"vitaminA"` // 维生素A（μg）
	VitaminC float64 `json:"vitaminC"` // 维生素C（mg）
	VitaminD float64 `json:"vitaminD"` // 维生素D（μg）
	VitaminE float64 `json:"vitaminE"` // 维生素E（mg）
	VitaminK float64 `json:"vitaminK"` // 维生素K（μg）
	VitaminB Complex `json:"vitaminB"` // B族维生素
}

// Complex 表示B族维生素
type Complex struct {
	B1  float64 `json:"b1"`  // 维生素B1（mg）
	B2  float64 `json:"b2"`  // 维生素B2（mg）
	B6  float64 `json:"b6"`  // 维生素B6（mg）
	B12 float64 `json:"b12"` // 维生素B12（μg）
}

// Minerals 表示矿物质含量
type Minerals struct {
	Calcium   float64 `json:"calcium"`   // 钙（mg）
	Iron      float64 `json:"iron"`      // 铁（mg）
	Sodium    float64 `json:"sodium"`    // 钠（mg）
	Potassium float64 `json:"potassium"` // 钾（mg）
	Zinc      float64 `json:"zinc"`      // 锌（mg）
	Magnesium float64 `json:"magnesium"` // 镁（mg）
}
