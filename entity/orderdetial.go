package entity

type OrderDetail struct {
	DetailID int64 `db:"detailid"`     // 订单明细唯一标识
	OrderID  int64 `db:"orderid"`      // 所属订单ID
	DishID   int64 `db:"dishid"`        // 菜品ID
	Quantity int64 `db:"quantity"`      // 菜品数量
}