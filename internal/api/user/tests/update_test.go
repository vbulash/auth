package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/vbulash/auth/internal/api/user"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/vbulash/auth/internal/model"
	"github.com/vbulash/auth/internal/service"
	serviceMocks "github.com/vbulash/auth/internal/service/mocks"
	desc "github.com/vbulash/auth/pkg/user_v1"
)

func TestUpdate(t *testing.T) {
	//t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx     context.Context
		request *desc.UpdateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		name  = wrappers.StringValue{Value: gofakeit.Name()}
		email = wrappers.StringValue{Value: gofakeit.Email()}
		role  = desc.Role_USER

		serviceErr = fmt.Errorf("ошибка при тестировании")

		request = &desc.UpdateRequest{
			Id:    id,
			Name:  &name,
			Email: &email,
			Role:  role,
		}

		info = &model.UserInfo{
			Name:  name.GetValue(),
			Email: email.GetValue(),
			Role:  int32(role),
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "Успешный вариант",
			args: args{
				ctx:     ctx,
				request: request,
			},
			err: nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, id, info).Return(nil)
				return mock
			},
		},
		{
			name: "Неуспешный вариант",
			args: args{
				ctx:     ctx,
				request: request,
			},
			err: serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, id, info).Return(serviceErr)
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

			_, err := api.Update(tt.args.ctx, tt.args.request)
			require.Equal(t, tt.err, err)
		})
	}
}
