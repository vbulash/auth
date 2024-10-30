package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vbulash/auth/internal/api"
	userAPI "github.com/vbulash/auth/internal/api/user"
	"github.com/vbulash/auth/internal/config"
	"github.com/vbulash/auth/internal/repository"
	userRepository "github.com/vbulash/auth/internal/repository/user"
	"github.com/vbulash/auth/internal/service"
	userService "github.com/vbulash/auth/internal/service/user"
	"log"
)

type serviceProvider struct {
	env          *config.Env
	pool         *pgxpool.Pool
	repoLayer    *repository.UserRepository
	serviceLayer *service.UserService
	apiLayer     *api.UserAPI
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// Env Конфигурация контейнера
func (s *serviceProvider) Env() *config.Env {
	if s.env == nil {
		env, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("Ошибка загрузки .env: %v", err)
		}
		s.env = env
	}

	return s.env
}

// Pool Пул соединений Postgres
func (s *serviceProvider) Pool(ctx context.Context) *pgxpool.Pool {
	if s.pool == nil {
		poolConfig, err := pgxpool.ParseConfig(s.Env().DSN)
		if err != nil {
			log.Fatalf("Ошибка конфигурации pgxpool: %v", err)
		}
		pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
		if err != nil {
			log.Fatalf("Ошибка коннекта к БД: %v", err)
		}

		s.pool = pool
	}
	return s.pool
}

// RepoLayer Слой репозитория
func (s *serviceProvider) RepoLayer(ctx context.Context) *repository.UserRepository {
	if s.repoLayer == nil {
		repoLayer := userRepository.NewUserRepository(s.Pool(ctx))
		s.repoLayer = &repoLayer
	}
	return s.repoLayer
}

// ServiceLayer Слой сервиса
func (s *serviceProvider) ServiceLayer(ctx context.Context) *service.UserService {
	if s.serviceLayer == nil {
		serviceLayer := userService.NewUserService(*s.RepoLayer(ctx))
		s.serviceLayer = &serviceLayer
	}
	return s.serviceLayer
}

// APILayer Слой API
func (s *serviceProvider) APILayer(ctx context.Context) *api.UserAPI {
	if s.apiLayer == nil {
		apiLayer := userAPI.NewAPI(*s.ServiceLayer(ctx))
		s.apiLayer = &apiLayer
	}
	return s.apiLayer
}
