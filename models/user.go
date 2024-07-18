package models

type User struct {
	ID       int64 `xorm:"pk"`
	Nickname string
	Card     string
	Sex      string
	Age      int
	Area     string
	Level    int
	Role     string
	Title    string
}
