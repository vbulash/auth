package app

import (
	"context"
	"log"

	"github.com/vbulash/platform_common/pkg/config"

	"github.com/vbulash/platform_common/pkg/client/db/transaction"

	"github.com/vbulash/platform_common/pkg/client/db"
	"github.com/vbulash/platform_common/pkg/client/db/pg"
	"github.com/vbulash/platform_common/pkg/closer"

	api "github.com/vbulash/auth/internal/api/user"
	userAPI "github.com/vbulash/auth/internal/api/user"
	"github.com/vbulash/auth/internal/repository"
	userRepository "github.com/vbulash/auth/internal/repository/user"
	"github.com/vbulash/auth/internal/service"
	userService "github.com/vbulash/auth/internal/service/user"
	"github.com/vbulash/platform_common/pkg/config/env"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient     db.Client
	txManager    db.TxManager
	repoLayer    repository.UserRepository
	serviceLayer service.UserService
	apiLayer     *api.UsersAPI
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

// DBClient Клиент БД
func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		client, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("Ошибка создания db клиента: %v", err)
		}

		err = client.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("Ошибка пинга db: %v", err)
		}

		closer.Add(client.Close)

		s.dbClient = client
	}

	return s.dbClient
}

// TxManager Менеджер транзакций
func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

// RepoLayer Слой репозитория
func (s *serviceProvider) RepoLayer(ctx context.Context) repository.UserRepository {
	if s.repoLayer == nil {
		repoLayer := userRepository.NewUserRepository(s.DBClient(ctx))
		s.repoLayer = repoLayer
	}

	return s.repoLayer
}

// ServiceLayer Слой сервиса
func (s *serviceProvider) ServiceLayer(ctx context.Context) service.UserService {
	if s.serviceLayer == nil {
		serviceLayer := userService.NewUserService(
			s.RepoLayer(ctx),
			s.TxManager(ctx),
		)
		s.serviceLayer = serviceLayer
	}

	return s.serviceLayer
}

// APILayer Слой API
func (s *serviceProvider) APILayer(ctx context.Context) *api.UsersAPI {
	if s.apiLayer == nil {
		apiLayer := userAPI.NewAPI(s.ServiceLayer(ctx))
		s.apiLayer = apiLayer
	}

	return s.apiLayer
}
