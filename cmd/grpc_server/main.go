package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/brianvoe/gofakeit"
	"github.com/golang/protobuf/ptypes/empty"
	user "github.com/vbulash/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = ":50051"

type server struct {
	user.UnimplementedAuthV1Server
}

func (s *server) Create(_ context.Context, _ *user.CreateRequest) (*user.CreateResponse, error) {
	fmt.Println("Сервер: создание пользователя")
	return &user.CreateResponse{
		Id: gofakeit.Int64(),
	}, nil
}

func (s *server) Get(_ context.Context, request *user.GetRequest) (*user.GetResponse, error) {
	fmt.Println("Сервер: получение пользователя")
	id := request.Id
	record := &user.GetResponse{
		Id:        id,
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      user.Role_USER,
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()), // Может так статься, что эта дата будет раньше CreatedAt - пока пофиг
	}
	return record, nil
}

func (s *server) Update(_ context.Context, _ *user.UpdateRequest) (*empty.Empty, error) {
	fmt.Println("Сервер: получение пользователя")
	return &empty.Empty{}, nil
}

func (s *server) Delete(_ context.Context, _ *user.DeleteRequest) (*empty.Empty, error) {
	fmt.Println("Сервер: удаление пользователя")
	return &empty.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s", grpcPort))
	if err != nil {
		log.Fatalf("Фатальная ошибка запуска / прослушивания: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	user.RegisterAuthV1Server(s, &server{})

	log.Printf("Сервер прослушивает: %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Фатальная ошибка запуска: %v", err)
	}
}
