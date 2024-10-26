package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vbulash/auth/internal/repository"
	"github.com/vbulash/auth/internal/repository/user/model"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/vbulash/auth/pkg/user_v1"
)

type repo struct {
	db *pgxpool.Pool
}

// NewUserRepository Создание репо
func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, info *desc.UserInfo) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx,
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		info.GetName(), info.GetEmail(), info.GetPassword()).Scan(&id)
	if err != nil {
		return 0, nil
	}
	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*desc.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx,
		"SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = $1", id).Scan(
		&user.ID, &user.Info.Name, &user.Info.Email, &user.Info.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Inline converter
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.User{
		Id: user.ID,
		Info: &desc.UserInfo{
			Name:     user.Info.Name,
			Email:    user.Info.Email,
			Password: user.Info.Password,
		},
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}, nil
}

func (r *repo) Update(ctx context.Context, id int64, info *desc.UserInfo) error {
	_, err := r.db.Exec(ctx,
		"UPDATE users SET name = $1, email = $2 WHERE id = $3",
		info.GetName(), info.GetEmail(), id)
	return err
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx,
		"DELETE FROM users WHERE id = $1", id)
	return err
}
