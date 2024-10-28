package converter

import (
	"database/sql"
	"time"

	"github.com/vbulash/auth/internal/model"
	desc "github.com/vbulash/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ModelUserInfoToDescUserInfo Преобразование из модели в GRPC
func ModelUserInfoToDescUserInfo(info *model.UserInfo) *desc.UserInfo {
	return &desc.UserInfo{
		Name:     info.Name,
		Email:    info.Email,
		Password: info.Password,
		Role:     desc.Role(info.Role),
	}
}

// ModelUserToDescUser Преобразование из модели в GRPC
func ModelUserToDescUser(info *model.User) *desc.User {
	var createdAt, updatedAt *timestamppb.Timestamp
	if info.UpdatedAt.Valid {
		updatedAt = timestamppb.New(info.UpdatedAt.Time)
	}
	createdAt = timestamppb.New(info.CreatedAt)

	return &desc.User{
		Id:        info.ID,
		Info:      ModelUserInfoToDescUserInfo(&info.Info),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// DescUserInfoToModelUserInfo Преобразование из GRPC в модель
func DescUserInfoToModelUserInfo(info *desc.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name:     info.Name,
		Email:    info.Email,
		Password: info.Password,
		Role:     int32(info.Role),
	}
}

// DescUserToModelUser Преобразование из GRPC в модель
func DescUserToModelUser(info *desc.User) *model.User {
	var createdAt time.Time
	var updateAt sql.NullTime
	createdAt = info.CreatedAt.AsTime()
	_ = updateAt.Scan(info.UpdatedAt.AsTime())

	translatedInfo := DescUserInfoToModelUserInfo(info.Info)

	return &model.User{
		ID:        info.Id,
		Info:      *translatedInfo,
		CreatedAt: createdAt,
		UpdatedAt: updateAt,
	}
}
