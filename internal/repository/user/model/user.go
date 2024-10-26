package model

import (
	"database/sql"
	"time"
)

// User Полная запись пользователя
type User struct {
	ID        int64        `db:"id"`
	Info      UserInfo     ``
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

// UserInfo Краткая информация по пользователю
type UserInfo struct {
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}
