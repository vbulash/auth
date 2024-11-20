package model

// User Полная запись пользователя
type User struct {
	ID        int64  `redis:"id"`
	Name      string `redis:"name"`
	Email     string `redis:"email"`
	Password  string `redis:"password"`
	Role      int32  `redis:"role"`
	CreatedAt int64  `redis:"created_at"`
	UpdatedAt *int64 `redis:"updated_at"`
}
