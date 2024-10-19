package models

type User struct {
	Id       string `db:"id"`
	Username string `db:"username"`
}
