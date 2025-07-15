package entity

type Dish struct {
	DishID      int64   `db:"dishid"`       // 菜品唯一标识
	DishName    string  `db:"dishname"`     // 菜品名称
	Price       float64 `db:"price"`        // 菜品价格
	Description string  `db:"description"`  // 菜品描述
	ImageURL    string  `db:"imageurl"`     // 菜品图片URL
}