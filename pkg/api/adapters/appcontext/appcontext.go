package appcontext

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/eser/ajan/configfx"
	"github.com/eser/ajan/datafx"
	"github.com/eser/ajan/logfx"
	"github.com/eser/ajan/metricsfx"
	"github.com/eser/ajan/queuefx"
	"github.com/eser/aya.is-services/pkg/api/adapters/arcade"
	_ "github.com/lib/pq"
)

var ErrInitFailed = errors.New("failed to initialize app context")

type AppContext struct {
	Config  *AppConfig
	Logger  *logfx.Logger
	Metrics *metricsfx.MetricsProvider

	Data  *datafx.Registry
	Queue *queuefx.Registry

	Arcade *arcade.Arcade
}

func New() *AppContext {
	return &AppContext{} //nolint:exhaustruct
}

func (a *AppContext) Init(ctx context.Context) error {
	// ----------------------------------------------------
	// Config
	// ----------------------------------------------------
	cl := configfx.NewConfigManager()

	a.Config = &AppConfig{} //nolint:exhaustruct

	err := cl.LoadDefaults(a.Config)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInitFailed, err)
	}

	// ----------------------------------------------------
	// Logger
	// ----------------------------------------------------
	a.Logger, err = logfx.NewLoggerAsDefault(os.Stdout, &a.Config.Log)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInitFailed, err)
	}

	a.Logger.InfoContext(
		ctx,
		"[AppContext] Initialization in progress",
		slog.String("module", "appcontext"),
		slog.String("name", a.Config.AppName),
		slog.String("environment", a.Config.AppEnv),
		slog.Any("features", a.Config.Features),
	)

	// ----------------------------------------------------
	// Metrics
	// ----------------------------------------------------
	a.Metrics = metricsfx.NewMetricsProvider()

	err = a.Metrics.RegisterNativeCollectors()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInitFailed, err)
	}

	// ----------------------------------------------------
	// Data
	// ----------------------------------------------------
	a.Data = datafx.NewRegistry(a.Logger)

	err = a.Data.LoadFromConfig(ctx, &a.Config.Data)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInitFailed, err)
	}

	// ----------------------------------------------------
	// Queue
	// ----------------------------------------------------
	a.Queue = queuefx.NewRegistry(a.Logger)

	err = a.Queue.LoadFromConfig(ctx, &a.Config.Queue)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInitFailed, err)
	}

	// ----------------------------------------------------
	// Arcade
	// ----------------------------------------------------
	a.Arcade = arcade.New(a.Config.Externals.Arcade)

	return nil
}
