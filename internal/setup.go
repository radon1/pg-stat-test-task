package internal

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"

	"github.com/radon1/pg-stat-test-task/internal/config"
	"github.com/radon1/pg-stat-test-task/internal/controllers"
	"github.com/radon1/pg-stat-test-task/internal/repositories"
	"github.com/radon1/pg-stat-test-task/internal/system/connections/pg"
	"github.com/radon1/pg-stat-test-task/internal/system/logger"
)

type app struct {
	ctx      context.Context
	cfg      *config.Config
	logger   *zerolog.Logger
	pgPool   *pgxpool.Pool
	fiberApp *fiber.App
}

func NewApp(ctx context.Context) (*app, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	appLogger, err := logger.New(cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	pgPool, err := pg.NewPool(ctx, cfg.PGDSN)
	if err != nil {
		return nil, err
	}

	fiberApp := fiber.New()
	repoContainer := repositories.New(pgPool)
	controllers.Register(appLogger, repoContainer, fiberApp)

	return &app{
		ctx:      ctx,
		cfg:      cfg,
		logger:   appLogger,
		pgPool:   pgPool,
		fiberApp: fiberApp,
	}, nil
}

func (a *app) Start() {
	if err := a.fiberApp.Listen(fmt.Sprintf(":%s", a.cfg.Port)); err != nil {
		a.logger.Error().Err(err).Msg("failed to start http server")
	}
}

func (a *app) Shutdown() error {
	if err := a.fiberApp.Shutdown(); err != nil {
		return err
	}
	a.pgPool.Close()
	return nil
}
