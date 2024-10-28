package api

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	desc "github.com/vbulash/auth/pkg/user_v1"
)

// UserAPI Слой API
type UserAPI interface {
	Create(ctx context.Context, request *desc.CreateRequest) (*desc.CreateResponse, error)
	Get(ctx context.Context, request *desc.GetRequest) (*desc.GetResponse, error)
	Update(ctx context.Context, request *desc.UpdateRequest) (*empty.Empty, error)
	Delete(ctx context.Context, request *desc.DeleteRequest) (*empty.Empty, error)
}
