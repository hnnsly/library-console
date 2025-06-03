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
	"golang.org/x/crypto/bcrypt"

	"github.com/hnnsly/library-console/internal/config"
	"github.com/hnnsly/library-console/internal/handler"
	"github.com/hnnsly/library-console/internal/logger"
	libraryrepository "github.com/hnnsly/library-console/internal/repository"

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

	if cfg.LibraryService == nil { // Изменено с PharmacyService на LibraryService
		log.Fatal().Msg("Library service configuration section is missing in config file")
	}
	if cfg.Db == nil {
		log.Fatal().Msg("Database configuration section is missing in config file")
	}
	if cfg.Rd == nil {
		log.Fatal().Msg("Redis configuration section is missing in config file")
	}
	log.Info().Msg("config is normal")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pgPool := mustOpenPg(ctx, cfg.Db.URL())
	defer pgPool.Close()

	// Инициализация sqlc Queries
	pgQueries := postgres.New(pgPool) // postgres.New ожидает DBTX, pgxpool.Pool реализует его

	ensureFirstAdmin(ctx, pgPool)

	redisClient := mustOpenRedis(ctx, *cfg.Rd)
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing Redis connection for Library service")
		}
	}()

	// Создаем репозиторий библиотеки
	repo := libraryrepository.New(pgQueries, redisClient)

	// Создаем хендлер с авторизацией
	h := handler.NewHandler(repo, *cfg.LibraryService, pgQueries, redisClient)
	app := h.Router()

	go startServer(app, cfg.LibraryService.Port, "Library service")

	<-ctx.Done()
	log.Info().Msg("Shutdown initiated for Library service")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Library service server shutdown error")
	} else {
		log.Info().Msg("Library service server gracefully stopped")
	}
}

func mustOpenPg(ctx context.Context, dsn string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't connect to PostgreSQL for Library service")
	}
	// Проверка соединения
	if err := pool.Ping(ctx); err != nil {
		pool.Close() // Закрыть пул, если пинг не удался
		log.Fatal().Err(err).Msg("Failed to ping PostgreSQL for Library service")
	}
	log.Info().Msg("Connected to PostgreSQL for Library service")
	return pool
}

func mustOpenRedis(ctx context.Context, rc config.RedisConfig) *redis.Redis {
	r, err := redis.New(ctx, rc)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't connect to Redis for Library service")
	}
	log.Info().Msg("Connected to Redis for Library service")
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

func ensureFirstAdmin(ctx context.Context, pool *pgxpool.Pool) {
	// Проверяем, есть ли первый админ
	var exists bool
	err := pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE is_first_admin = true AND is_active = true)").Scan(&exists)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check first admin")
		return
	}

	if exists {
		log.Info().Msg("First admin already exists")
		return
	}

	// Создаем первого админа
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin-password"), bcrypt.DefaultCost)

	_, err = pool.Exec(ctx, `
		INSERT INTO users (username, email, password_hash, role, full_name, is_active, is_first_admin, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, true, true, NOW(), NOW())`,
		"root", "root@library.local", string(hashedPassword), "super_admin", "System Administrator")

	if err != nil {
		log.Error().Err(err).Msg("Failed to create first admin")
		return
	}

	log.Warn().
		Str("username", "root").
		Str("password", "admin-password").
		Msg("🚀 First admin created!")
}
