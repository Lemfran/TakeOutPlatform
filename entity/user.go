package entity

type User struct {
	UserID   int64  `db:"userid"` // 注意字段名需要与数据库中的列名对应
	Username string `db:"username"`
	Password string `db:"password"`
	Email    string `db:"email"`
	Manager  bool   `db:"manager"` 
}