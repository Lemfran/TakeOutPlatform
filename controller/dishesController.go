package controller

import (
	"TakeOutPlatform/entity"
	"TakeOutPlatform/handle"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
)

var dishListTemplate = template.Must(template.ParseFiles(filepath.Join("templates", "dishlistmanager.html")))
var dishListUserTemplate = template.Must(template.ParseFiles(filepath.Join("templates", "dishlistuser.html")))
var dishAddTemplate = template.Must(template.ParseFiles(filepath.Join("templates", "dishadd.html")))
var dishEditTemplate = template.Must(template.ParseFiles(filepath.Join("templates", "dishedit.html")))

// DishListHandler 显示所有菜品信息
func DishListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		Dishes, err := handle.GetAllDishes()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dishListTemplate.Execute(w, Dishes)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// DishListHandler 显示所有菜品信息
func DishListUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		Dishes, err := handle.GetAllDishes()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dishListUserTemplate.Execute(w, Dishes)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// DishAddHandler 添加新菜品
func DishAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		dishAddTemplate.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		var dish entity.Dish
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dish.DishName = r.FormValue("dishname")
		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			http.Error(w, "Invalid price format", http.StatusBadRequest)
			return
		}
		dish.Price = price
		dish.Description = r.FormValue("description")
		dish.ImageURL = r.FormValue("imageurl")

		err = handle.CreateDish(dish)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/dishes", http.StatusFound)
	}
}

// DishEditHandler 修改菜品信息
func DishEditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		dishID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid dish ID", http.StatusBadRequest)
			return
		}
		dish, err := handle.GetDish(dishID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dishEditTemplate.Execute(w, dish)
	} else if r.Method == http.MethodPost {
		var dish entity.Dish
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dishID, err := strconv.ParseInt(r.FormValue("dishid"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid dish ID", http.StatusBadRequest)
			return
		}
		dish.DishID = dishID
		dish.DishName = r.FormValue("dishname")
		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			http.Error(w, "Invalid price format", http.StatusBadRequest)
			return
		}
		dish.Price = price
		dish.Description = r.FormValue("description")
		dish.ImageURL = r.FormValue("imageurl")

		err = handle.UpdateDish(dish)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/manager/dishes", http.StatusFound)
	}
}

// DishDeleteHandler 删除菜品
func DishDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		dishID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid dish ID", http.StatusBadRequest)
			return
		}
		err = handle.DeleteDish(dishID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/manager/dishes", http.StatusFound)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func DishSearchController(w http.ResponseWriter, r *http.Request) {
	// 检查请求方法是否为GET
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method. Only GET is allowed.", http.StatusMethodNotAllowed)
		return
	}

	// 从请求中获取搜索关键词
	keyword := r.URL.Query().Get("query")
	if keyword == "" {
		http.Redirect(w, r, "/manager/dishes", http.StatusFound)

	}

	// 调用搜索函数
	Dishes, err := handle.SearchDishesByName(keyword)
	if err != nil {
		http.Error(w, "Failed to search dishes", http.StatusInternalServerError)
		return
	}

	// 渲染HTML模板
	tmpl := template.Must(template.ParseFiles("templates/dishlistmanager.html"))
	tmpl.Execute(w, Dishes)
}

func DishSearchUserController(w http.ResponseWriter, r *http.Request) {
	// 检查请求方法是否为GET
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method. Only GET is allowed.", http.StatusMethodNotAllowed)
		return
	}

	// 从请求中获取搜索关键词
	keyword := r.URL.Query().Get("query")
	if keyword == "" {
		http.Redirect(w, r, "/user/dishes", http.StatusFound)

	}

	// 调用搜索函数
	Dishes, err := handle.SearchDishesByName(keyword)
	if err != nil {
		http.Error(w, "Failed to search dishes", http.StatusInternalServerError)
		return
	}

	// 渲染HTML模板
	tmpl := template.Must(template.ParseFiles("templates/dishlistuser.html"))
	tmpl.Execute(w, Dishes)
}

// TopDishesHandler 处理获取点菜次数排行前十的请求
func TopDishesHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
        return
    }

    topDishes, err := handle.GetTopDishes()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 将结果传递给前端模板
    tmpl := template.Must(template.ParseFiles("templates/topdishes.html"))
    tmpl.Execute(w, topDishes)
}