package entity

type Order struct {
	OrderID int64 `db:"orderid"`       // 订单唯一标识
	UserID  int64 `db:"userid"`        // 下单用户ID
	OrderStatus int `db:"orderstatus"`  // 订单状态（0：未支付，1：已支付，2：已取消）
}