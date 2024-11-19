package user

import (
	"context"

	"github.com/vbulash/auth/internal/converter"
	"github.com/vbulash/auth/internal/model"
	"github.com/vbulash/auth/internal/repository"
	"github.com/vbulash/auth/internal/service"
)

type serviceLayer struct {
	repoLayer repository.UserRepository
}

// NewUserService Создание сервисного слоя
func NewUserService(repo repository.UserRepository) service.UserService {
	return &serviceLayer{
		repoLayer: repo,
	}
}

func (s *serviceLayer) Create(ctx context.Context, info *model.UserInfo) (int64, error) {
	id, err := s.repoLayer.Create(ctx, converter.ModelUserInfoToDescUserInfo(info))
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *serviceLayer) Get(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.repoLayer.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *serviceLayer) Update(ctx context.Context, id int64, info *model.UserInfo) error {
	if len(info.Name) == 0 && len(info.Email) == 0 && info.Role == 0 {
		return nil
	}

	err := s.repoLayer.Update(ctx, id, converter.ModelUserInfoToDescUserInfo(info))
	if err != nil {
		return err
	}

	return nil
}

func (s *serviceLayer) Delete(ctx context.Context, id int64) error {
	err := s.repoLayer.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
