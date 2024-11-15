package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/vbulash/auth/internal/converter"
	"github.com/vbulash/auth/internal/repository"

	"github.com/vbulash/auth/internal/model"
	"github.com/vbulash/auth/internal/service/user"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	repositoryMocks "github.com/vbulash/auth/internal/repository/mocks"
	desc "github.com/vbulash/auth/pkg/user_v1"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx     context.Context
		request *model.UserInfo
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, true, true, true, true, 32)
		role     = desc.Role_USER

		serviceErr = fmt.Errorf("ошибка при тестировании")

		request = &desc.UserInfo{
			Name:     name,
			Email:    email,
			Password: password,
			Role:     role,
		}

		info = converter.DescUserInfoToModelUserInfo(request)
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		userRepositoryMock userRepositoryMockFunc
	}{
		{
			name: "Успешный вариант",
			args: args{
				ctx:     ctx,
				request: info,
			},
			want: id,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, request).Return(id, nil)
				return mock
			},
		},
		{
			name: "Неуспешный вариант",
			args: args{
				ctx:     ctx,
				request: info,
			},
			want: 0,
			err:  serviceErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, request).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepositoryMock := tt.userRepositoryMock(mc)
			// Упрощенный вариант инициализации сервиса - без менеджера транзакций
			service := user.NewUserService(userRepositoryMock, nil)

			resHandler, err := service.Create(tt.args.ctx, tt.args.request)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resHandler)
		})
	}
}
