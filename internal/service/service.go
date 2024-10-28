package service

import (
	"context"

	"github.com/vbulash/auth/internal/model"
)

// UserService Сервис пользователей
type UserService interface {
	Create(ctx context.Context, info *model.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, id int64, info *model.UserInfo) error
	Delete(ctx context.Context, id int64) error
}
