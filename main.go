package main

import (
	"TakeOutPlatform/controller"
	"TakeOutPlatform/database"
	"log"
	"net/http"
)

func main(){
	 // 初始化数据库连接
    err := database.InitDB()
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer database.CloseDB()
    log.Println("Database connected successfully")

    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
   // 开启server服务


    // 设置 HTTP 路由
    mux := http.NewServeMux()
    mux.HandleFunc("/register", controller.RegisterHandler)
    mux.HandleFunc("/login", controller.LoginHandler)
    mux.HandleFunc("/logout", controller.LogoutHandler)
    mux.HandleFunc("/manager/dishes",controller.DishListHandler)
    mux.HandleFunc("/manager/dishes/add",controller.DishAddHandler)
    mux.HandleFunc("/manager/dishes/edit",controller.DishEditHandler)
    mux.HandleFunc("/manager/dishes/delete",controller.DishDeleteHandler)
    mux.HandleFunc("/manager/userlist",controller.UserListHandler)
    mux.HandleFunc("/manager/userlist/updateManager",controller.UpdateManagerHandler)
    mux.HandleFunc("/manager/orders",controller.OrdersHandler)
    mux.HandleFunc("/manager/orders/details",controller.GetOrderDetailsHandler)
    mux.HandleFunc("/user/dishes",controller.DishListUserHandler)
    mux.HandleFunc("/user/dishes/add",controller.CreateOrderHandler)
    mux.HandleFunc("/user/orders",controller.OrdersUserHandler)
    mux.HandleFunc("/user/orders/details",controller.GetOrderDetailsUserHandler)
    mux.HandleFunc("/user/orders/pay",controller.PayOrderHandler)
    mux.HandleFunc("/user/orders/cancel",controller.CancelOrderHandler)
    mux.HandleFunc("/manager/dishes/search",controller.DishSearchController)
    mux.HandleFunc("/user/dishes/search",controller.DishSearchUserController)
    mux.HandleFunc("/user/topdishes",controller.TopDishesHandler)

    // 启动 HTTP 服务器
    log.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}