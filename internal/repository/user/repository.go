package user

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
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
	creates := make(map[string]interface{})
	creates["name"] = info.GetName()
	creates["email"] = info.GetEmail()
	creates["password"] = info.GetPassword()

	query, args, err := squirrel.Insert("chats").
		SetMap(creates).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, nil
	}

	var id int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, nil
	}
	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*desc.User, error) {
	query, args, err := squirrel.
		Select("id, name, email, created_at, updated_at").
		From("users").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	var user model.User
	err = r.db.QueryRow(ctx, query, args...).
		Scan(&user.ID, &user.Info.Name, &user.Info.Email, &user.CreatedAt, &user.UpdatedAt)
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
	bUpdated := false
	updates := make(map[string]interface{})
	if len(info.GetName()) > 0 {
		updates["name"] = info.GetName()
		bUpdated = true
	}
	if len(info.GetEmail()) > 0 {
		updates["email"] = info.GetEmail()
		bUpdated = true
	}
	if bUpdated {
		updates["updated_at"] = time.Now()
	}

	query, args, err := squirrel.Update("users").
		SetMap(updates).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, query, args...)
	return err
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	query, args, err := squirrel.Delete("users").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, query, args...)
	return err
}
