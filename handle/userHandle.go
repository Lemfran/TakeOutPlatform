package handle

import (
	"TakeOutPlatform/database"
	"TakeOutPlatform/entity"
)

// CreateUser 创建一个新的用户
func CreateUser(user entity.User) error {	
	stmt, err := database.DB.Prepare("INSERT INTO users (username, password, email, manager) VALUES (?, ?, ?, ?)")
    if err != nil {
        return err
    }
    _, err = stmt.Exec(user.Username, user.Password, user.Email, user.Manager)
    if err != nil {
        return err
    }
    return nil
}

func GetUser(userID int64) (entity.User, error) {
    var user entity.User
    err := database.DB.QueryRow("SELECT * FROM users WHERE userid = ?", userID).Scan(&user.UserID, &user.Username, &user.Password, &user.Email, &user.Manager)
    if err != nil {
        return user, err
    }
    return user, nil
}

func DeleteUser(userID int64) error {
    stmt, err := database.DB.Prepare("DELETE FROM users WHERE userid = ?")
    if err != nil {
        return err
    }
    _, err = stmt.Exec(userID)
    if err != nil {
        return err
    }
    
    return nil
}

func GetAllUsers() ([]entity.User, error) {
	rows, err := database.DB.Query("SELECT userid, username, email, manager FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.Email, &user.Manager); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func UpdateUserManager(user entity.User) error {
	// 准备 SQL 更新语句
	stmt, err := database.DB.Prepare("UPDATE users SET manager = ? WHERE userid = ?")
	if err != nil {
		return err
	}
	defer stmt.Close() // 确保语句在函数退出时关闭

	// 执行更新操作
	_, err = stmt.Exec(user.Manager, user.UserID)
	if err != nil {
		return err
	}
	return nil
}

func GetUserByUsername(username string) (entity.User, error) {
    // 定义一个 User 结构体变量用于存储查询结果
    var user entity.User

    // 查询数据库，获取匹配用户名的第一个用户
    err := database.DB.QueryRow("SELECT userid, username, email, manager FROM users WHERE username = ? LIMIT 1", username).Scan(&user.UserID, &user.Username, &user.Email, &user.Manager)
    if err != nil {
        return user, err
    }

    return user, nil
}