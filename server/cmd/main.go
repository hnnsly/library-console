package main

import (
	"context"
	"flag"
	"fmt"
	golog "log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	// Используйте правильные пути к вашим пакетам
	"github.com/hnnsly/library-console/internal/config"
	"github.com/hnnsly/library-console/internal/handler"
	"github.com/hnnsly/library-console/internal/logger"
	pharmacyrepository "github.com/hnnsly/library-console/internal/repository"

	// Предполагается, что у вас есть пакет, сгенерированный sqlc
	"github.com/hnnsly/library-console/internal/repository/postgres"
	"github.com/hnnsly/library-console/internal/repository/redis"
)

func main() {
	cfgPath := flag.String("c", "config.yml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		golog.Fatalf("Error loading config: %v", err)
	}

	if cfg.Log == nil {
		golog.Fatal("Logger configuration is missing in config file")
	}
	log.Logger = *logger.Setup(cfg.Log)

	if cfg.PharmacyService == nil {
		log.Fatal().Msg("Pharmacy service configuration section is missing in config file")
	}
	if cfg.Db == nil {
		log.Fatal().Msg("Database configuration section is missing in config file")
	}
	if cfg.Rd == nil {
		log.Fatal().Msg("Redis configuration section is missing in config file")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pgPool := mustOpenPg(ctx, cfg.Db.URL())
	defer pgPool.Close()

	// Инициализация sqlc Queries
	pgQueries := postgres.New(pgPool) // pharmacyPostgres.New ожидает DBTX, pgxpool.Pool реализует его

	redisClient := mustOpenRedis(ctx, *cfg.Rd)
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing Redis connection for Pharmacy service")
		}
	}()

	repo := pharmacyrepository.New(pgQueries, redisClient)
	h := handler.NewHandler(repo, *cfg.PharmacyService)
	app := h.Router()

	go startServer(app, cfg.PharmacyService.Port, "Pharmacy service")

	<-ctx.Done()
	log.Info().Msg("Shutdown initiated for Pharmacy service")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Pharmacy service server shutdown error")
	} else {
		log.Info().Msg("Pharmacy service server gracefully stopped")
	}
}

func mustOpenPg(ctx context.Context, dsn string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't connect to PostgreSQL for Pharmacy service")
	}
	// Проверка соединения
	if err := pool.Ping(ctx); err != nil {
		pool.Close() // Закрыть пул, если пинг не удался
		log.Fatal().Err(err).Msg("Failed to ping PostgreSQL for Pharmacy service")
	}
	log.Info().Msg("Connected to PostgreSQL for Pharmacy service")
	return pool
}

func mustOpenRedis(ctx context.Context, rc config.RedisConfig) *redis.Redis {
	r, err := redis.New(ctx, rc)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't connect to Redis for Pharmacy service")
	}
	log.Info().Msg("Connected to Redis for Pharmacy service")
	return r
}

func startServer(app *fiber.App, port int, serviceName string) {
	if port == 0 {
		log.Fatal().Msgf("%s port is not configured or is zero", serviceName)
	}
	addr := fmt.Sprintf(":%d", port)
	log.Info().Msgf("%s starting on %s", serviceName, addr)
	if err := app.Listen(addr); err != nil {
		// Проверка на ошибку syscall.EINVAL, которая может возникнуть при завершении работы
		// и попытке снова прослушать тот же порт перед полным освобождением.
		// В Fiber это обычно обрабатывается через ShutdownWithContext.
		if !os.IsTimeout(err) && !strings.Contains(err.Error(), "server closed") && !strings.Contains(err.Error(), "Listener closed") {
			log.Fatal().Err(err).Msgf("%s HTTP server crashed", serviceName)
		} else {
			log.Info().Err(err).Msgf("%s HTTP server stopped", serviceName)
		}
	}
}
