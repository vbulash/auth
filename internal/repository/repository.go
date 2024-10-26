package repository

import (
	"context"

	desc "github.com/vbulash/auth/pkg/user_v1"
)

// UserRepository Репо пользователя
type UserRepository interface {
	Create(ctx context.Context, info *desc.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*desc.User, error)
	Update(ctx context.Context, id int64, info *desc.UserInfo) error
	Delete(ctx context.Context, id int64) error
}
