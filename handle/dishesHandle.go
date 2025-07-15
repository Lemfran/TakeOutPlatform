package handle

import (
	"TakeOutPlatform/database"
	"TakeOutPlatform/entity"
)

// CreateDish 创建一个新的菜品
func CreateDish(dish entity.Dish) error {
	stmt, err := database.DB.Prepare("INSERT INTO dishes (dishname, price, description, imageurl) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(dish.DishName, dish.Price, dish.Description, dish.ImageURL)
	if err != nil {
		return err
	}
	return nil
}

// GetDish 根据菜品ID获取菜品信息
func GetDish(dishID int64) (entity.Dish, error) {
	var dish entity.Dish
	err := database.DB.QueryRow("SELECT * FROM dishes WHERE dishid = ?", dishID).Scan(&dish.DishID, &dish.DishName, &dish.Price, &dish.Description, &dish.ImageURL)
	if err != nil {
		return dish, err
	}
	return dish, nil
}

// UpdateDish 更新菜品信息
func UpdateDish(dish entity.Dish) error {
	stmt, err := database.DB.Prepare("UPDATE dishes SET dishname = ?, price = ?, description = ?, imageurl = ? WHERE dishid = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(dish.DishName, dish.Price, dish.Description, dish.ImageURL, dish.DishID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteDish 根据菜品ID删除菜品
func DeleteDish(dishID int64) error {
	stmt, err := database.DB.Prepare("DELETE FROM dishes WHERE dishid = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(dishID)
	if err != nil {
		return err
	}
	return nil
}

// SearchDishesByName 根据菜名模糊查询菜品
func SearchDishesByName(keyword string) ([]entity.Dish, error) {
	// 使用 LIKE 进行模糊查询，注意防止 SQL 注入
	query := "SELECT dishid, dishname, price, description, imageurl FROM dishes WHERE dishname LIKE ?"
	rows, err := database.DB.Query(query, "%"+keyword+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dishes []entity.Dish
	for rows.Next() {
		var dish entity.Dish
		if err := rows.Scan(&dish.DishID, &dish.DishName, &dish.Price, &dish.Description, &dish.ImageURL); err != nil {
			return nil, err
		}
		dishes = append(dishes, dish)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dishes, nil
}

// GetAllDishes 获取所有菜品信息
func GetAllDishes() ([]entity.Dish, error) {
	rows, err := database.DB.Query("SELECT dishid, dishname, price, description, imageurl FROM dishes")
	if err != nil {
		return nil, err
	}
	var dishes []entity.Dish
	for rows.Next() {
		var dish entity.Dish
		if err := rows.Scan(&dish.DishID, &dish.DishName, &dish.Price, &dish.Description, &dish.ImageURL); err != nil {
			return nil, err
		}
		dishes = append(dishes, dish)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dishes, nil
}


// GetTopDishes 获取点菜次数排行前十的菜品
func GetTopDishes() ([]entity.Dish, error) {
    var Dishes []entity.Dish
    rows, err := database.DB.Query(`
        SELECT od.dishid AS total_quantity
		FROM orderdetails od
		JOIN orders o ON od.orderid = o.orderid
		WHERE o.orderstatus = 1
		GROUP BY od.dishid
		ORDER BY total_quantity DESC
		LIMIT 10;
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var dish entity.Dish
        if err := rows.Scan(&dish.DishID); err != nil {
            return nil, err
        }
		dish1,_:=GetDish(dish.DishID)
        Dishes = append(Dishes, dish1)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return Dishes, nil
}