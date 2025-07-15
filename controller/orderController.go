package controller

import (
	"TakeOutPlatform/entity"
	"TakeOutPlatform/handle"
	"html/template"
	"net/http"
	"strconv"
)

var ordersTemplate = template.Must(template.ParseFiles("templates/orders.html"))
var orderDetailsTemplate = template.Must(template.ParseFiles("templates/orderdetails.html"))
var ordersUserTemplate = template.Must(template.ParseFiles("templates/ordersuser.html"))
var orderDetailsUserTemplate = template.Must(template.ParseFiles("templates/orderdetailsuser.html"))
// OrdersHandler 显示订单列表
func OrdersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	orders, err := handle.GetAllOrders()
	if err != nil {
		ordersTemplate.Execute(w, nil)
		return
	}

	ordersTemplate.Execute(w, orders)
}

// OrdersHandler 显示订单列表
func OrdersUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("userID")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userID, err := strconv.ParseInt(cookie.Value, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	orders, err := handle.GetOrdersByUserID(userID)
	if err != nil {
		ordersUserTemplate.Execute(w, nil)
		return
	}

	ordersUserTemplate.Execute(w, orders)
}

func GetOrderDetailsHandler(w http.ResponseWriter, r *http.Request) {
    orderIDStr := r.URL.Query().Get("orderid")
    orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid order ID", http.StatusBadRequest)
        return
    }

    details, err := handle.GetOrderDetailsByOrderID(orderID)
    if err != nil {
        orderDetailsTemplate.Execute(w, nil)
        return
    }

    // 计算总金额
	var totalAmount float64
	for _, detail := range details {
		totalAmount += detail.Price
	}

	type OrderDetails struct {
		Details       []struct {
			DetailID  int64
			OrderID   int64
			DishName  string
			Quantity int64
			Price    float64
		}
		TotalAmount float64
	}

	// 将总金额添加到模板数据中
	var data OrderDetails
	data.Details = details
	data.TotalAmount = totalAmount
	orderDetailsTemplate.Execute(w, data)
}

func GetOrderDetailsUserHandler(w http.ResponseWriter, r *http.Request) {
    orderIDStr := r.URL.Query().Get("orderid")
    orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid order ID", http.StatusBadRequest)
        return
    }

    details, err := handle.GetOrderDetailsByOrderID(orderID)
    if err != nil {
        orderDetailsUserTemplate.Execute(w, nil)
    }

    // 计算总金额
	var totalAmount float64
	for _, detail := range details {
		totalAmount += detail.Price
	}

	type OrderDetails struct {
		Details       []struct {
			DetailID  int64
			OrderID   int64
			DishName  string
			Quantity int64
			Price    float64
		}
		TotalAmount float64
		OrderStatus int64
		OrderID int64
	}

	// 将总金额添加到模板数据中
	var data OrderDetails
	data.Details = details
	data.TotalAmount = totalAmount
	data.OrderStatus = handle.GetOrderStatusByOrderID(orderID)
	data.OrderID = orderID

	orderDetailsUserTemplate.Execute(w, data)
}

// 创建订单并添加菜品
func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		cookie, err := r.Cookie("userID")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userID, err := strconv.ParseInt(cookie.Value, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		orderID, err := handle.CreateOrderIfNeeded(userID)
		if err != nil {
			http.Redirect(w, r, "/user/dishes", http.StatusSeeOther)
		}
		dishID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid dish ID", http.StatusBadRequest)
			return
		}
		detail := entity.OrderDetail{
			OrderID:  orderID,
			DishID:   dishID,
			Quantity: 1,
		}
		err = handle.AddOrderDetail(detail)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 重定向到订单详情页面或返回订单详情
		http.Redirect(w, r, "/user/dishes", http.StatusSeeOther)
	}
}

// PayOrderHandler 处理支付订单的请求
func PayOrderHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
		// 解析表单数据
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
	
    orderIDstr := r.FormValue("orderID") // 从请求中获取订单ID
	orderID, err := strconv.ParseInt(orderIDstr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid orderID", http.StatusBadRequest)
			return
		}
    newStatus := 1

    // 更新订单状态
    err = handle.UpdateOrderStatus(orderID, newStatus)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 重定向到订单列表页面或其他页面
    http.Redirect(w, r, "/user/orders", http.StatusSeeOther)
	}
}

func CancelOrderHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
		// 解析表单数据
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
	
    orderIDstr := r.FormValue("orderID") // 从请求中获取订单ID
	orderID, err := strconv.ParseInt(orderIDstr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid orderID", http.StatusBadRequest)
			return
		}
    newStatus := 2

    // 更新订单状态
    err = handle.UpdateOrderStatus(orderID, newStatus)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 重定向到订单列表页面或其他页面
    http.Redirect(w, r, "/user/orders", http.StatusSeeOther)
	}
}

