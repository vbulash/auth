package operations

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	// pq Используется sqlx для работы с postgres
	_ "github.com/lib/pq"
	"github.com/vbulash/auth/config"
)

// Get Получение данных из таблицы notes
func Get(db *sqlx.DB) (*[]config.UserType, error) {
	users := []config.UserType{}
	err := db.Select(&users, "SELECT * FROM users")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &users, nil
}
