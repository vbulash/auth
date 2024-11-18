package converter

import (
	"database/sql"
	"time"

	"github.com/vbulash/auth/internal/model"
	modelRepo "github.com/vbulash/auth/internal/repository/user/redis/model"
)

// ToUserFromRepo ...
func ToUserFromRepo(user *modelRepo.User) *model.User {
	var updatedAt sql.NullTime
	if user.UpdatedAt != nil {
		updatedAt = sql.NullTime{
			Time:  time.Unix(0, *user.UpdatedAt),
			Valid: true,
		}
	}

	return &model.User{
		ID: user.ID,
		Info: model.UserInfo{
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
			Role:     user.Role,
		},
		CreatedAt: time.Unix(0, user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}
