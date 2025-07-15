package handle

import (
	"TakeOutPlatform/database"
	"TakeOutPlatform/entity"
	"errors"
	"database/sql"
)

// CreateOrderIfNeeded 创建订单（如果不存在则创建）
func CreateOrderIfNeeded(userID int64) (int64, error) {
    var orderID int64
    err := database.DB.QueryRow("SELECT orderid FROM orders WHERE userid = ? AND orderstatus=0", userID).Scan(&orderID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            // 如果没有找到订单，创建一个新的订单
            order := entity.Order{
                UserID:   userID,
                OrderStatus: 0, // 假设0表示未支付状态
            }
            err = CreateOrder(order)
            if err != nil {
                return 0, err
            }
            var orderid int64
            err:=database.DB.QueryRow("SELECT orderid FROM orders WHERE userid = ? AND orderstatus=0", userID).Scan(&orderid)
            if err !=nil{
                return 0, err
            }
            return orderid, nil
        } else {
            return 0, err
        }
    }
    return orderID, nil
}

// CreateOrder 创建一个新的订单
func CreateOrder(order entity.Order) error {
    // 插入订单记录
    stmt, err := database.DB.Prepare("INSERT INTO orders (userid, orderstatus) VALUES (?, ?)")
    if err != nil {
        return err
    }
    _, err = stmt.Exec(order.UserID, order.OrderStatus)
    if err != nil {
        return err
    }
    return nil
}

// AddOrderDetail 添加订单明细
func AddOrderDetail(detail entity.OrderDetail) error {
	stmt, err := database.DB.Prepare("INSERT INTO orderdetails (orderid, dishid, quantity) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(detail.OrderID, detail.DishID, detail.Quantity)
	return err
}


// UpdateOrderStatus 更新订单状态
func UpdateOrderStatus(orderID int64, status int) error {
	stmt, err := database.DB.Prepare("UPDATE orders SET orderstatus = ? WHERE orderid = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(status, orderID)
	return err
}

// DeleteOrder 删除订单
func DeleteOrder(orderID int64) error {
	// 删除订单明细
	if _, err := database.DB.Exec("DELETE FROM order_details WHERE orderid = ?", orderID); err != nil {
		return err
	}

	// 删除订单
	if _, err := database.DB.Exec("DELETE FROM orders WHERE orderid = ?", orderID); err != nil {
		return err
	}

	return nil
}



func GetAllOrders() ([]struct {
	OrderID    int64
	Username  string
	OrderStatus int
}, error) {
	rows, err := database.DB.Query(`
		SELECT 
			o.orderid,
			u.username,
			o.orderstatus
		FROM 
			orders o
		JOIN 
			users u ON o.userid = u.userid
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []struct {
		OrderID    int64
		Username  string
		OrderStatus int
	}
	for rows.Next() {
		var order struct {
			OrderID    int64
			Username  string
			OrderStatus int
		}
		if err := rows.Scan(&order.OrderID, &order.Username, &order.OrderStatus); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func GetOrderStatusByOrderID(orderID int64) int64 {
    var status int64
    // 准备SQL查询语句
    err := database.DB.QueryRow("SELECT orderstatus FROM orders WHERE orderid = ?", orderID).Scan(&status)
    if err != nil {
        if err == sql.ErrNoRows {
            // 如果没有找到记录，返回错误或特定状态值
            return 0 // 或者返回一个特定的错误
        }
        return 0
    }
    return status
}

func GetOrderDetailsByOrderID(orderID int64) ([]struct {
    DetailID  int64
    OrderID   int64
    DishName string
    Quantity int64
    Price   float64
}, error) {
    rows, err := database.DB.Query(`
        SELECT 
            od.detailid,
            od.orderid,
            d.dishname,
            od.quantity,
            d.price
        FROM 
            orderdetails od
        JOIN 
            dishes d ON od.dishid = d.dishid
        WHERE 
            od.orderid = ?
    `, orderID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var details []struct {
        DetailID  int64
        OrderID   int64
        DishName string
        Quantity int64
        Price   float64
    }
    for rows.Next() {
        var detail struct {
            DetailID  int64
            OrderID   int64
            DishName string
            Quantity int64
            Price   float64
        }
        if err := rows.Scan(&detail.DetailID, &detail.OrderID, &detail.DishName, &detail.Quantity, &detail.Price); err != nil {
            return nil, err
        }
        details = append(details, detail)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return details, nil
}

// GetOrdersByUserID 根据用户ID获取订单
func GetOrdersByUserID(userID int64) ([]struct {
    OrderID    int64
    Username  string
    OrderStatus int
}, error) {
    // 准备SQL查询语句
    rows, err := database.DB.Query(`
        SELECT 
            o.orderid,
            u.username,
            o.orderstatus
        FROM 
            orders o
        JOIN 
            users u ON o.userid = u.userid
        WHERE 
            o.userid = ?
    `, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var orders []struct {
        OrderID    int64
        Username  string
        OrderStatus int
    }
    for rows.Next() {
        var order struct {
            OrderID    int64
            Username  string
            OrderStatus int
        }
        if err := rows.Scan(&order.OrderID, &order.Username, &order.OrderStatus); err != nil {
            return nil, err
        }
        orders = append(orders, order)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    if len(orders) == 0 {
        return nil, errors.New("no orders found for the given user ID")
    }

    return orders, nil
}