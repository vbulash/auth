package pg

import (
	"context"
	"fmt"

	innermodel "github.com/vbulash/auth/internal/repository/user/pg/model"

	"strings"
	"time"

	outermodel "github.com/vbulash/auth/internal/model"

	"github.com/vbulash/platform_common/pkg/client/db"

	"github.com/Masterminds/squirrel"
	"github.com/vbulash/auth/internal/repository"

	desc "github.com/vbulash/auth/pkg/user_v1"
)

const (
	tableName string = "users"

	idColumn        string = "id"
	nameColumn      string = "name"
	emailColumn     string = "email"
	passwordColumn  string = "password"
	roleColumn      string = "role"
	createdAtColumn string = "created_at"
	updatedAtColumn string = "updated_at"
)

type repoLayer struct {
	db db.Client
}

// NewUserRepository Создание репо-слоя
func NewUserRepository(db db.Client) repository.UserRepository {
	return &repoLayer{db: db}
}

func (r *repoLayer) Create(ctx context.Context, info *desc.UserInfo) (int64, error) {
	creates := make(map[string]interface{})
	creates[nameColumn] = info.GetName()
	creates[emailColumn] = info.GetEmail()
	creates[passwordColumn] = info.GetPassword()
	creates[roleColumn] = int32(info.GetRole())

	query, args, err := squirrel.Insert(tableName).
		SetMap(creates).
		Suffix(fmt.Sprintf("RETURNING \"%s\"", idColumn)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, nil
	}

	var id int64
	q := db.Query{
		Name:     tableName + ".Create",
		QueryRaw: query,
	}
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, nil
	}
	return id, nil
}

func (r *repoLayer) Get(ctx context.Context, id int64) (*outermodel.User, error) {
	columns := []string{
		idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn,
	}
	query, args, err := squirrel.
		Select(strings.Join(columns, ", ")).
		From(tableName).
		Where(squirrel.Eq{idColumn: id}).
		PlaceholderFormat(squirrel.Dollar).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}
	var user innermodel.User
	q := db.Query{
		Name:     tableName + ".Get",
		QueryRaw: query,
	}
	err = r.db.DB().QueryRowContext(ctx, q, args...).
		Scan(&user.ID, &user.Info.Name, &user.Info.Email, &user.Info.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Преобразование внутренней модели во внешнюю - не стал выносить в конвертер
	// innermodel -> outermodel
	return &outermodel.User{
		ID: user.ID,
		Info: outermodel.UserInfo{
			Name:     user.Info.Name,
			Email:    user.Info.Email,
			Role:     user.Info.Role,
			Password: user.Info.Password,
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *repoLayer) Update(ctx context.Context, id int64, info *desc.UserInfo) error {
	updates := make(map[string]interface{})
	if len(info.GetName()) > 0 {
		updates[nameColumn] = info.GetName()
	}
	if len(info.GetEmail()) > 0 {
		updates[emailColumn] = info.GetEmail()
	}
	if info.GetRole() != 0 {
		updates[roleColumn] = info.GetRole()
	}
	updates[updatedAtColumn] = time.Now()

	query, args, err := squirrel.Update(tableName).
		SetMap(updates).
		Where(squirrel.Eq{idColumn: id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	q := db.Query{
		Name:     tableName + ".Update",
		QueryRaw: query,
	}
	_, err = r.db.DB().ExecContext(ctx, q, args...)

	return err
}

func (r *repoLayer) Delete(ctx context.Context, id int64) error {
	query, args, err := squirrel.Delete(tableName).
		Where(squirrel.Eq{idColumn: id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	q := db.Query{
		Name:     tableName + ".Delete",
		QueryRaw: query,
	}
	_, err = r.db.DB().ExecContext(ctx, q, args...)

	return err
}
