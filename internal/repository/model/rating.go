package model

type Rating struct {
	UserID int32 `db:"user_id"`
	ChatID int32 `db:"chat_id"`
	Rating int32 `db:"rating"`
	Level  int32 `db:"level"`
}
