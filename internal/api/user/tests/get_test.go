package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/vbulash/auth/internal/api/user"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/vbulash/auth/internal/model"
	"github.com/vbulash/auth/internal/service"
	serviceMocks "github.com/vbulash/auth/internal/service/mocks"
	desc "github.com/vbulash/auth/pkg/user_v1"
)

func TestGet(t *testing.T) {
	//t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx     context.Context
		request *desc.GetRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		name      = gofakeit.Name()
		email     = gofakeit.Email()
		password  = gofakeit.Password(true, true, true, true, true, 32)
		role      = desc.Role_USER
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		serviceErr = fmt.Errorf("ошибка при тестировании")

		request = &desc.GetRequest{
			Id: id,
		}

		serviceResponse = &model.User{
			ID: id,
			Info: model.UserInfo{
				Name:     name,
				Email:    email,
				Password: password,
				Role:     int32(role),
			},
			CreatedAt: createdAt,
			UpdatedAt: sql.NullTime{
				Valid: true,
				Time:  updatedAt,
			},
		}

		response = &desc.GetResponse{
			Id:        id,
			Name:      name,
			Email:     email,
			Role:      role,
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *desc.GetResponse
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
				mock.GetMock.Expect(ctx, id).Return(serviceResponse, nil)
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
				mock.GetMock.Expect(ctx, id).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			//t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := user.NewAPI(userServiceMock)

			resHandler, err := api.Get(tt.args.ctx, tt.args.request)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resHandler)
		})
	}
}
