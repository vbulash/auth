package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vbulash/auth/internal/repository/user"

	"github.com/vbulash/auth/internal/repository"

	"github.com/vbulash/auth/config"

	"github.com/golang/protobuf/ptypes/empty"
	desc "github.com/vbulash/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	desc.UnimplementedAuthV1Server
	userRepository repository.UserRepository
}

func (s *server) Create(ctx context.Context, request *desc.CreateRequest) (*desc.CreateResponse, error) {
	fmt.Println("Сервер: создание пользователя")

	id, err := s.userRepository.Create(ctx, &desc.UserInfo{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		return nil, err
	}
	return &desc.CreateResponse{
		Id: id,
	}, nil
}

func (s *server) Get(ctx context.Context, request *desc.GetRequest) (*desc.GetResponse, error) {
	fmt.Println("Сервер: получение пользователя")

	userObj, err := s.userRepository.Get(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &desc.GetResponse{
		Id:    userObj.Id,
		Name:  userObj.Info.Name,
		Email: userObj.Info.Email,
		Role:  userObj.Info.Role,
	}, nil
}

func (s *server) Update(ctx context.Context, request *desc.UpdateRequest) (*empty.Empty, error) {
	fmt.Println("Сервер: обновление пользователя")

	err := s.userRepository.Update(ctx, request.Id, &desc.UserInfo{
		Name:  request.Name.Value,
		Email: request.Email.Value,
		Role:  request.Role,
	})
	return &empty.Empty{}, err
}

func (s *server) Delete(ctx context.Context, request *desc.DeleteRequest) (*empty.Empty, error) {
	fmt.Println("Сервер: удаление пользователя")
	err := s.userRepository.Delete(ctx, request.Id)
	return &empty.Empty{}, err
}

func main() {
	ctx := context.Background()

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env: %v", err)
	}
	config.Config = conf

	poolConfig, err := pgxpool.ParseConfig(os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalf("Ошибка конфигурации pgxpool: %v", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Ошибка коннекта к БД: %v", err)
	}

	userRepo := user.NewUserRepository(pool)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Config.ServerPort))
	if err != nil {
		log.Fatalf("Фатальная ошибка запуска / прослушивания: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{
		userRepository: userRepo,
	})

	log.Printf("Сервер прослушивает: %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Фатальная ошибка запуска: %v", err)
	}
}
