package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/akemoon/crowdfunding-app-auth/api"
	userClient "github.com/akemoon/crowdfunding-app-auth/cluster/user/resty"
	"github.com/akemoon/crowdfunding-app-auth/config"
	_ "github.com/akemoon/crowdfunding-app-auth/docs"
	infraRedis "github.com/akemoon/crowdfunding-app-auth/infra/redis"
	"github.com/akemoon/crowdfunding-app-auth/repo/creds/postgres"
	redisRepo "github.com/akemoon/crowdfunding-app-auth/repo/token/redis"
	authService "github.com/akemoon/crowdfunding-app-auth/service/auth"
	"github.com/akemoon/crowdfunding-app-auth/service/creds"
	"github.com/akemoon/crowdfunding-app-auth/service/token"
	"github.com/akemoon/crowdfunding-app-auth/tool/hasher/bcrypt"
	"github.com/redis/go-redis/v9"
)

const (
	pgMigrationsDir = "/app/migrations/postgres"
	envPgDSN        = "POSTGRES_DSN"

	envRedisAddr = "REDIS_ADDR"
	envRedisDB   = "REDIS_DB"
	envRedisPass = "REDIS_PASSWORD"

	envJWTSecret      = "JWT_SECRET"
	envUserServiceURL = "USER_SERVICE_URL"
)

// @title Auth Service API
// @version 1.0
func main() {
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pg, err := initPostgres(mainCtx)
	if err != nil {
		log.Fatalf("init db err: %s", err)
	}
	defer func() {
		if err := pg.Close(); err != nil {
			log.Printf("close db err: %s", err)
		}
	}()

	redisClient, err := initRedis(mainCtx)
	if err != nil {
		log.Fatalf("init redis err: %s", err)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("close redis err: %s", err)
		}
	}()

	jwtSecret := strings.TrimSpace(os.Getenv(envJWTSecret))
	if jwtSecret == "" {
		log.Fatalf("env %s is empty", envJWTSecret)
	}

	userServiceURL := strings.TrimSpace(os.Getenv(envUserServiceURL))
	if userServiceURL == "" {
		log.Fatalf("env %s is empty", envUserServiceURL)
	}

	tokenRepo := redisRepo.NewRefreshTokenRepository(redisClient)
	tokenSvc := token.NewService(tokenRepo, jwtSecret)

	credsRepo := postgres.NewCredsRepo(pg)
	hasher := bcrypt.NewHasher(0)
	credsSvc := creds.NewService(credsRepo, hasher)

	userSvc := userClient.NewClient(userServiceURL)
	authSvc := authService.NewService(userSvc, credsSvc, tokenSvc)

	srv := api.NewServer()
	srv.AddAuthHandlers(authSvc)
	srv.AddTokenHandlers(tokenSvc)
	srv.AddSwaggerUI()
	srv.AddMetrics()

	srv.ListenAndServe(":80")
}

func initPostgres(ctx context.Context) (*sql.DB, error) {
	dsn := strings.TrimSpace(os.Getenv(envPgDSN))
	if dsn == "" {
		return nil, fmt.Errorf("missing required env var %s", envPgDSN)
	}

	db, err := postgres.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("db connection error: %s", err)
	}

	err = postgres.Migrate(ctx, db, pgMigrationsDir)
	if err != nil {
		return nil, fmt.Errorf("db migration error: %s", err)
	}

	return db, nil
}

func initRedis(ctx context.Context) (*redis.Client, error) {
	for _, key := range []string{
		envRedisAddr,
		envRedisDB,
	} {
		value := strings.TrimSpace(os.Getenv(key))
		if value == "" {
			return nil, fmt.Errorf("missing required env var: %s", key)
		}
	}

	db, err := strconv.Atoi(strings.TrimSpace(os.Getenv(envRedisDB)))
	if err != nil {
		return nil, fmt.Errorf("invalid %s: %w", envRedisDB, err)
	}

	cfg := config.Redis{
		Addr:     strings.TrimSpace(os.Getenv(envRedisAddr)),
		Password: os.Getenv(envRedisPass),
		DB:       db,
	}

	return infraRedis.NewRedisClient(ctx, cfg)
}
