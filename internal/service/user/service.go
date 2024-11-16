package user

import (
	"context"

	"github.com/vbulash/platform_common/pkg/client/db"

	"github.com/vbulash/auth/internal/converter"
	"github.com/vbulash/auth/internal/model"
	"github.com/vbulash/auth/internal/repository"
	"github.com/vbulash/auth/internal/service"
)

type serviceLayer struct {
	repoLayer repository.UserRepository
	txManager db.TxManager
}

// NewUserService Создание сервисного слоя
func NewUserService(repo repository.UserRepository, txManager db.TxManager) service.UserService {
	return &serviceLayer{
		repoLayer: repo,
		txManager: txManager,
	}
}

func (s *serviceLayer) Create(ctx context.Context, info *model.UserInfo) (int64, error) {
	var id int64
	var err error
	if nil == s.txManager { // Ветка для упрощенного юнит-теста
		id, err = s.repoLayer.Create(ctx, converter.ModelUserInfoToDescUserInfo(info))
		if err != nil {
			return 0, err
		}
	} else {
		err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
			var err error
			id, err = s.repoLayer.Create(ctx, converter.ModelUserInfoToDescUserInfo(info))
			if err != nil {
				return err
			}
			// ..
			return nil
		})
	}

	// Транслируем внутреннюю ошибку во внешнюю без преобразования - хотя надо бы
	return id, err
}

func (s *serviceLayer) Get(ctx context.Context, id int64) (*model.User, error) {
	var user *model.User
	var err error
	if nil == s.txManager { // Ветка для упрощенного юнит-теста
		user, err = s.repoLayer.Get(ctx, id)
	} else {
		err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
			var err error
			user, err = s.repoLayer.Get(ctx, id)
			if err != nil {
				return err
			}
			// ...
			return nil
		})
	}

	// Транслируем внутреннюю ошибку во внешнюю без преобразования - хотя надо бы
	return user, err
}

func (s *serviceLayer) Update(ctx context.Context, id int64, info *model.UserInfo) error {
	if len(info.Name) == 0 && len(info.Email) == 0 && info.Role == 0 {
		return nil
	}

	var err error
	if nil == s.txManager { // Ветка для упрощенного юнит-теста
		err = s.repoLayer.Update(ctx, id, converter.ModelUserInfoToDescUserInfo(info))
	} else {
		err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
			err := s.repoLayer.Update(ctx, id, converter.ModelUserInfoToDescUserInfo(info))
			if err != nil {
				return err
			}
			// ...
			return nil
		})
	}

	// Транслируем внутреннюю ошибку во внешнюю без преобразования - хотя надо бы
	return err
}

func (s *serviceLayer) Delete(ctx context.Context, id int64) error {
	var err error

	if nil == s.txManager { // Ветка для упрощенного юнит-теста
		err = s.repoLayer.Delete(ctx, id)
	} else {
		err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
			err := s.repoLayer.Delete(ctx, id)
			if err != nil {
				return err
			}
			// ...
			return nil
		})
	}

	// Транслируем внутреннюю ошибку во внешнюю без преобразования - хотя надо бы
	return err
}
