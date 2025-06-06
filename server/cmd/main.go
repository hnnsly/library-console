package main

import (
	"context"
	"flag"
	"fmt"
	golog "log"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hnnsly/library-console/internal/config"
	"github.com/hnnsly/library-console/internal/handler"
	"github.com/hnnsly/library-console/internal/logger"
	"github.com/hnnsly/library-console/internal/repository"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	"github.com/hnnsly/library-console/internal/repository/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func main() {
	cfgPath := flag.String("c", "config.yml", "path to config")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		golog.Fatalf("config loading error: %v", err)
	}
	log.Logger = *logger.Setup(cfg.Log)

	// Context for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize dependencies
	pgPool := mustOpenPg(ctx, cfg.Db.URL())
	defer pgPool.Close()

	rd := mustOpenRedis(ctx, *cfg.Rd)
	defer rd.Close()

	repo := repository.New(postgres.New(pgPool), rd)

	// Create API handler and Fiber app
	h := handler.NewHandler(repo, cfg.Library)
	app := h.Router()

	// Start server
	go startServer(app, cfg.Library.Port)

	<-ctx.Done()
	log.Info().Msg("Shutdown initiated")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Library server shutdown error")
	} else {
		log.Info().Msg("Library server gracefully stopped")
	}
}

func mustOpenPg(ctx context.Context, dsn string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't connect to PostgreSQL")
	}
	log.Info().Msg("Connected to PostgreSQL")
	return pool
}

func mustOpenRedis(ctx context.Context, rc config.Redis) *redis.Redis {
	r, err := redis.New(ctx, rc)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't connect to Redis")
	}
	log.Info().Msg("Connected to Redis")
	return r
}

func startServer(app *fiber.App, port int) {
	addr := fmt.Sprintf(":%d", port)
	log.Info().Msgf("Identity starting on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatal().Err(err).Msg("HTTP server crashed")
	}
}
