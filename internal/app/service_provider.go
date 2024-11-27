package app

import (
	"context"
	"log"

	redigo "github.com/gomodule/redigo/redis"
	userRepositoryPg "github.com/vbulash/auth/internal/repository/user/pg"
	userRepositoryRedis "github.com/vbulash/auth/internal/repository/user/redis"
	"github.com/vbulash/platform_common/pkg/client/cache"
	"github.com/vbulash/platform_common/pkg/client/cache/redis"

	"github.com/vbulash/platform_common/pkg/config"

	"github.com/vbulash/platform_common/pkg/client/db"
	"github.com/vbulash/platform_common/pkg/client/db/pg"
	"github.com/vbulash/platform_common/pkg/closer"

	api "github.com/vbulash/auth/internal/api/user"
	userAPI "github.com/vbulash/auth/internal/api/user"
	"github.com/vbulash/auth/internal/repository"
	"github.com/vbulash/auth/internal/service"
	userService "github.com/vbulash/auth/internal/service/user"
	"github.com/vbulash/platform_common/pkg/config/env"
)

type serviceProvider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	httpConfig    config.HTTPConfig
	redisConfig   config.RedisConfig
	storageConfig config.StorageConfig
	swaggerConfig config.SwaggerConfig

	redisPool   *redigo.Pool
	redisClient cache.RedisClient

	dbClient     db.Client
	repoLayer    repository.UserRepository
	serviceLayer service.UserService
	apiLayer     *api.UsersAPI
}

const (
	redisMode = "redis"
	pgMode    = "pg"
)

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

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := env.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redis config: %s", err.Error())
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

func (s *serviceProvider) StorageConfig() config.StorageConfig {
	if s.storageConfig == nil {
		cfg, err := env.NewStorageConfig()
		if err != nil {
			log.Fatalf("failed to get storage config: %s", err.Error())
		}

		s.storageConfig = cfg
	}

	return s.storageConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := env.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
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

func (s *serviceProvider) RedisPool() *redigo.Pool {
	if s.redisPool == nil {
		s.redisPool = &redigo.Pool{
			MaxIdle:     s.RedisConfig().MaxIdle(),
			IdleTimeout: s.RedisConfig().IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", s.RedisConfig().Address())
			},
		}
	}

	return s.redisPool
}

func (s *serviceProvider) RedisClient() cache.RedisClient {
	if s.redisClient == nil {
		s.redisClient = redis.NewClient(s.RedisPool(), s.RedisConfig())
	}

	return s.redisClient
}

// RepoLayer Слой репозитория
func (s *serviceProvider) RepoLayer(ctx context.Context) repository.UserRepository {
	var repoLayer repository.UserRepository
	if s.repoLayer == nil {
		switch s.StorageConfig().Mode() {
		case redisMode:
			repoLayer = userRepositoryRedis.NewUserRepository(s.RedisClient())
			break
		case pgMode:
			repoLayer = userRepositoryPg.NewUserRepository(s.DBClient(ctx))
			break
		default:
			repoLayer = nil
		}
		s.repoLayer = repoLayer
	}

	return s.repoLayer
}

// ServiceLayer Слой сервиса
func (s *serviceProvider) ServiceLayer(ctx context.Context) service.UserService {
	if s.serviceLayer == nil {
		serviceLayer := userService.NewUserService(s.RepoLayer(ctx))
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
