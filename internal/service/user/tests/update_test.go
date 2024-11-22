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

func TestUpdate(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx     context.Context
		id      int64
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

		zeroRequest = &desc.UserInfo{
			Name:     "",
			Email:    "",
			Password: "",
			Role:     desc.Role_UNKNOWN,
		}

		info     = converter.DescUserInfoToModelUserInfo(request)
		zeroInfo = converter.DescUserInfoToModelUserInfo(zeroRequest)
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		err                error
		userRepositoryMock userRepositoryMockFunc
	}{
		{
			name: "Успешный вариант",
			args: args{
				ctx:     ctx,
				id:      id,
				request: info,
			},
			err: nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.UpdateMock.Expect(ctx, id, request).Return(nil)
				return mock
			},
		},
		{
			name: "Неуспешный вариант",
			args: args{
				ctx:     ctx,
				id:      id,
				request: info,
			},
			err: serviceErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.UpdateMock.Expect(ctx, id, request).Return(serviceErr)
				return mock
			},
		},
		{
			name: "Пустой вариант",
			args: args{
				ctx:     ctx,
				id:      id,
				request: zeroInfo,
			},
			err: nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				//mock.UpdateMock.Expect(ctx, id, zeroRequest).Return(nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepositoryMock := tt.userRepositoryMock(mc)
			service := user.NewUserService(userRepositoryMock)

			err := service.Update(tt.args.ctx, tt.args.id, tt.args.request)
			require.Equal(t, tt.err, err)
		})
	}
}
