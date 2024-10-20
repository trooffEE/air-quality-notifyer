package models

type User struct {
	Id         int64
	Username   string `db:"username"`
	TelegramId string `db:"telegram_id"`
}
