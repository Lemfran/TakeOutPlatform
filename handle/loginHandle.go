package handle

import (
	"TakeOutPlatform/database"
	"TakeOutPlatform/entity"
	"database/sql"
	"errors"
	/*"crypto/md5"
	"encoding/hex"*/)

var (
    ErrUserNotFound       = errors.New("user not found")
    ErrInvalidCredentials = errors.New("invalid credentials")
)

/*func HashPassword(password string) string {
    hash := md5.New()
    hash.Write([]byte(password))
    return hex.EncodeToString(hash.Sum(nil))
}*/

func RegisterUser(user entity.User) error {
    /*hashedPassword := HashPassword(user.Password)*/

    query := `INSERT INTO users (Username,Password,Email) VALUES (?, ?, ?)`
    _, err := database.DB.Exec(query, user.Username, user.Password,user.Email)
    return err
}

// LoginUser 根据用户名和密码登录用户
func LoginUser(username, password string) (entity.User, error) {
	var user entity.User

	// 构建查询语句
	query := `SELECT userid, username, password, email, manager FROM users WHERE username = ? AND password = ?`
	// 执行查询
	row := database.DB.QueryRow(query, username, password)
	// 扫描查询结果到 user 结构体
	err := row.Scan(&user.UserID, &user.Username, &user.Password, &user.Email, &user.Manager)
	if err != nil {
		if err == sql.ErrNoRows {
			// 如果没有找到用户，返回自定义错误
			return user, errors.New("invalid username or password")
		}
		// 如果发生其他错误，返回该错误
		return user, err
	}
	// 如果查询成功，返回 user 结构体
	return user, nil
}