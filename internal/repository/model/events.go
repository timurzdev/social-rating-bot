package model

import "time"

type Event struct {
	MessageID int32     `db:"message_id"`
	UserID    int32     `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	IsDeleted bool      `db:"is_deleted"`
}
