package model

import (
	"database/sql"
	"time"
)

// User Полная запись пользователя
type User struct {
	ID        int64
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

// UserInfo Краткая информация по пользователю
type UserInfo struct {
	Name     string
	Email    string
	Password string
	Role     int32
}