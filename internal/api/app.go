package api

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	user3 "github.com/vbulash/auth/internal/api/user"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vbulash/auth/config"
	"github.com/vbulash/auth/internal/repository/user"
	user2 "github.com/vbulash/auth/internal/service/user"

	desc "github.com/vbulash/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// AppRun Инициализация и запуск приложения
func AppRun() error {
	ctx := context.Background()

	conf, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("Ошибка загрузки .env: %v", err)
	}
	config.Config = conf

	poolConfig, err := pgxpool.ParseConfig(os.Getenv("DB_DSN"))
	if err != nil {
		return fmt.Errorf("Ошибка конфигурации pgxpool: %v", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("Ошибка коннекта к БД: %v", err)
	}

	userRepo := user.NewUserRepository(pool)
	serviceLayer := user2.NewUserService(userRepo)
	apiLayer := user3.NewAPI(serviceLayer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Config.ServerPort))
	if err != nil {
		return fmt.Errorf("Фатальная ошибка запуска / прослушивания: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, apiLayer)

	log.Printf("Сервер ожидает вызовов: %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		return fmt.Errorf("Фатальная ошибка запуска сервера: %v", err)
	}

	return nil
}
