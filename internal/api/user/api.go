package user

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vbulash/auth/internal/model"
	"github.com/vbulash/auth/internal/service"
	desc "github.com/vbulash/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UsersAPI Слой API
type UsersAPI struct {
	desc.UnimplementedAuthV1Server
	serviceLayer service.UserService
}

// NewAPI Создание API
func NewAPI(serviceLayer service.UserService) *UsersAPI {
	return &UsersAPI{serviceLayer: serviceLayer}
}

// Create Создание пользователя
func (apiLayer *UsersAPI) Create(ctx context.Context, request *desc.CreateRequest) (*desc.CreateResponse, error) {
	fmt.Println("Сервер: создание пользователя")

	id, err := apiLayer.serviceLayer.Create(ctx, &model.UserInfo{
		Name:     request.GetName(),
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
		Role:     int32(request.GetRole()),
	})
	if err != nil {
		return nil, err
	}
	return &desc.CreateResponse{
		Id: id,
	}, nil
}

// Get Получение пользователя
func (apiLayer *UsersAPI) Get(ctx context.Context, request *desc.GetRequest) (*desc.GetResponse, error) {
	fmt.Println("Сервер: получение пользователя")

	userObj, err := apiLayer.serviceLayer.Get(ctx, request.GetId())
	if err != nil {
		return nil, err
	}

	var createdAt, updatedAt *timestamppb.Timestamp
	if userObj.UpdatedAt.Valid {
		updatedAt = timestamppb.New(userObj.UpdatedAt.Time)
	}
	createdAt = timestamppb.New(userObj.CreatedAt)

	return &desc.GetResponse{
		Id:        userObj.ID,
		Name:      userObj.Info.Name,
		Email:     userObj.Info.Email,
		Role:      desc.Role(userObj.Info.Role),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

// Update Изменение пользователя
func (apiLayer *UsersAPI) Update(ctx context.Context, request *desc.UpdateRequest) (*empty.Empty, error) {
	fmt.Println("Сервер: обновление пользователя")

	err := apiLayer.serviceLayer.Update(ctx, request.Id, &model.UserInfo{
		Name:  request.Name.GetValue(),
		Email: request.Email.GetValue(),
		Role:  int32(request.Role),
	})
	return &empty.Empty{}, err
}

// Delete Удаление пользователя
func (apiLayer *UsersAPI) Delete(ctx context.Context, request *desc.DeleteRequest) (*empty.Empty, error) {
	fmt.Println("Сервер: удаление пользователя")
	err := apiLayer.serviceLayer.Delete(ctx, request.Id)
	return &empty.Empty{}, err
}
