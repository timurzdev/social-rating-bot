package model

type Permissions struct {
	Level    int32 `db:"level"`
	Media    bool  `db:"media"`
	Messages bool  `db:"messages"`
}
