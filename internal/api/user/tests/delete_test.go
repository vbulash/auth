package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/vbulash/auth/internal/api/user"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/vbulash/auth/internal/service"
	serviceMocks "github.com/vbulash/auth/internal/service/mocks"
	desc "github.com/vbulash/auth/pkg/user_v1"
)

func TestDelete(t *testing.T) {
	//t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx     context.Context
		request *desc.DeleteRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		serviceErr = fmt.Errorf("ошибка при тестировании")

		request = &desc.DeleteRequest{
			Id: id,
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
				mock.DeleteMock.Expect(ctx, id).Return(nil)
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
				mock.DeleteMock.Expect(ctx, id).Return(serviceErr)
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

			_, err := api.Delete(tt.args.ctx, tt.args.request)
			require.Equal(t, tt.err, err)
		})
	}
}
