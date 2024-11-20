package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/vbulash/auth/internal/api/user"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/vbulash/auth/internal/model"
	"github.com/vbulash/auth/internal/service"
	serviceMocks "github.com/vbulash/auth/internal/service/mocks"
	desc "github.com/vbulash/auth/pkg/user_v1"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx     context.Context
		request *desc.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		name                 = gofakeit.Name()
		email                = gofakeit.Email()
		password             = gofakeit.Password(true, true, true, true, true, 32)
		passwordConfirmation = password
		role                 = desc.Role_USER

		serviceErr = fmt.Errorf("ошибка при тестировании")

		request = &desc.CreateRequest{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: passwordConfirmation,
			Role:            role,
		}

		info = &model.UserInfo{
			Name:     name,
			Email:    email,
			Password: password,
			Role:     int32(role),
		}

		response = &desc.CreateResponse{
			Id: id,
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "Успешный вариант",
			args: args{
				ctx:     ctx,
				request: request,
			},
			want: response,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, info).Return(id, nil)
				return mock
			},
		},
		{
			name: "Неуспешный вариант",
			args: args{
				ctx:     ctx,
				request: request,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, info).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := user.NewAPI(userServiceMock)

			resHandler, err := api.Create(tt.args.ctx, tt.args.request)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resHandler)
		})
	}
}
