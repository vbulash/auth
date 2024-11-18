package pg

import (
	"context"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	desc "github.com/vbulash/auth/pkg/user_v1"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/vbulash/auth/internal/model"
	"github.com/vbulash/auth/internal/repository"
	"github.com/vbulash/auth/internal/repository/user/redis/converter"
	modelRepo "github.com/vbulash/auth/internal/repository/user/redis/model"
	"github.com/vbulash/platform_common/pkg/client/cache"
)

type repoLayer struct {
	cl cache.RedisClient
}

// NewUserRepository ...
func NewUserRepository(cl cache.RedisClient) repository.UserRepository {
	return &repoLayer{cl: cl}
}

func (r *repoLayer) Create(ctx context.Context, info *desc.UserInfo) (int64, error) {
	id := gofakeit.Int64()

	user := &modelRepo.User{
		ID:        id,
		Name:      info.Name,
		Email:     info.Email,
		Password:  info.Password,
		Role:      int32(info.Role),
		CreatedAt: time.Now().UnixNano(),
	}

	idStr := strconv.FormatInt(id, 10)
	err := r.cl.HashSet(ctx, idStr, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repoLayer) Get(ctx context.Context, id int64) (*model.User, error) {
	idStr := strconv.FormatInt(id, 10)
	values, err := r.cl.HGetAll(ctx, idStr)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, model.ErrorUserNotFound
	}

	var user modelRepo.User
	err = redigo.ScanStruct(values, &user)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

func (r *repoLayer) Update(ctx context.Context, id int64, info *desc.UserInfo) error {
	idStr := strconv.FormatInt(id, 10)
	updatedAt := time.Now().UnixNano()

	user := &modelRepo.User{
		ID:        id,
		Name:      info.Name,
		Email:     info.Email,
		Password:  info.Password,
		Role:      int32(info.Role),
		UpdatedAt: &updatedAt,
	}

	err := r.cl.HashSet(ctx, idStr, user)
	if err != nil {
		return err
	}

	return nil
}

func (r *repoLayer) Delete(ctx context.Context, id int64) error {
	idStr := strconv.FormatInt(id, 10)
	err := r.cl.Expire(ctx, idStr, 1)
	if err != nil {
		return err
	}

	return nil
}
