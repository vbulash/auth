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
	return &serviceLayer{repoLayer: repo}
}

func (s *serviceLayer) Create(ctx context.Context, info *model.UserInfo) (int64, error) {
	return s.repoLayer.Create(ctx, converter.ModelUserInfoToDescUserInfo(info))
}
func (s *serviceLayer) Get(ctx context.Context, id int64) (*model.User, error) {
	nonConverted, err := s.repoLayer.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return converter.DescUserToModelUser(nonConverted), nil
}

func (s *serviceLayer) Update(ctx context.Context, id int64, info *model.UserInfo) error {
	return s.repoLayer.Update(ctx, id, converter.ModelUserInfoToDescUserInfo(info))
}

func (s *serviceLayer) Delete(ctx context.Context, id int64) error {
	return s.repoLayer.Delete(ctx, id)
}
