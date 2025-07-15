package controller

import (
	"TakeOutPlatform/entity"
	"TakeOutPlatform/handle"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

var registerTemplate = template.Must(template.ParseFiles(filepath.Join("templates", "register.html")))
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        registerTemplate.Execute(w, nil)
    } else if r.Method == http.MethodPost {
        var user entity.User
        err := r.ParseForm()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        user.Username = r.FormValue("username")
        user.Password = r.FormValue("password")
        user.Email = r.FormValue("email")

        err = handle.RegisterUser(user)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        http.Redirect(w, r, "/login", http.StatusFound)
    }
}

var loginTemplate = template.Must(template.ParseFiles(filepath.Join("templates", "login.html")))

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        // 发送登录页面
        loginTemplate.Execute(w, nil)
    } else if r.Method == http.MethodPost {
        // 处理登录逻辑
        username := r.FormValue("username")
        password := r.FormValue("password")
        user, err := handle.LoginUser(username, password)
		cookie := &http.Cookie{
        Name:  "userID",  // Cookie的名称
        Value: fmt.Sprintf("%d", user.UserID), // Cookie的值，这里是用户ID
        Path:  "/",       // Cookie的有效路径，"/"表示整个网站
        // 你可以设置Cookie的过期时间
        Expires: time.Now().Add(24 * time.Hour),
        // 设置Cookie为HttpOnly，这样JavaScript就不能访问Cookie了，增加安全性
        HttpOnly: true,
        Secure: true,
  	 	}
  		// 将Cookie发送给客户端
    	http.SetCookie(w, cookie)
        if err != nil {
            if err == handle.ErrUserNotFound || err == handle.ErrInvalidCredentials {
                // 登录失败，可以重定向回登录页面并显示错误信息
                // 或者直接在响应中显示错误信息
                http.Error(w, "Invalid username or password", http.StatusUnauthorized)
            } else {
                // 其他错误
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
            return
        } else if user.Manager{
            // 登录成功，重定向到主页或其他页面
     		http.Redirect(w, r, "/manager/dishes", http.StatusFound)
        } else if !user.Manager{
			http.Redirect(w, r, "/user/dishes", http.StatusFound)
		}
    }
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    // 创建一个新的Cookie，过期时间设置为过去的时间，这样Cookie就会被删除
    cookie := &http.Cookie{
        Name:   "userID",
        Value:  "",
        Path:   "/",
        Expires: time.Unix(0, 0), // 过期时间设置为Unix时间的起点，即1970年1月1日
    }

    // 将Cookie发送给客户端
    http.SetCookie(w, cookie)

    // 重定向到登录页面或主页
    http.Redirect(w, r, "/login", http.StatusFound)
}

func UserListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// 调用服务层获取所有用户信息（不包括密码）
		users, err := handle.GetAllUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 使用模板渲染用户列表页面
		userListTemplate := template.Must(template.ParseFiles("templates/userlist.html"))
		userListTemplate.Execute(w, users)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// UpdateManagerHandler 处理用户管理员状态的更新
func UpdateManagerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// 解析表单数据
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		// 获取用户ID
		userIDStr := r.FormValue("userId")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// 获取当前管理员状态
		managerStr := r.FormValue("manager")
		manager := managerStr == "true"

		// 创建用户对象
		user := entity.User{
			UserID:  userID,
			Manager: !manager, // 切换状态
		}

		// 调用服务层更新管理员状态
		err = handle.UpdateUserManager(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 重定向回用户列表页面
		http.Redirect(w, r, "/userlist", http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
