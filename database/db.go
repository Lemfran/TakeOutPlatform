package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // 导入MySQL驱动
)

var DB *sql.DB

// InitDB 初始化数据库连接
func InitDB() error {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/takeoutplatform")
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	DB = db

	fmt.Println("Database connection successful")
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}