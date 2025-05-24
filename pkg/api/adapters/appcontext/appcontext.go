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
	"github.com/eser/aya.is-services/pkg/api/adapters/storage"
	"github.com/eser/aya.is-services/pkg/api/business/profiles"
	"github.com/eser/aya.is-services/pkg/api/business/stories"
	"github.com/eser/aya.is-services/pkg/api/business/users"
	_ "github.com/lib/pq"
)

var ErrInitFailed = errors.New("failed to initialize app context")

type AppContext struct {
	// Adapters
	Config  *AppConfig
	Logger  *logfx.Logger
	Metrics *metricsfx.MetricsProvider

	Data  *datafx.Registry
	Queue *queuefx.Registry

	Arcade *arcade.Arcade

	Repository *storage.Repository

	// Business
	ProfilesService   *profiles.Service
	UsersService      *users.Service
	UsersOAuthService *users.GitHubOAuthService
	StoriesService    *stories.Service
}

func New() *AppContext {
	return &AppContext{} //nolint:exhaustruct
}

func (a *AppContext) Init(ctx context.Context) error {
	// ----------------------------------------------------
	// Adapter: Config
	// ----------------------------------------------------
	cl := configfx.NewConfigManager()

	a.Config = &AppConfig{} //nolint:exhaustruct

	err := cl.LoadDefaults(a.Config)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInitFailed, err)
	}

	// ----------------------------------------------------
	// Adapter: Logger
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
	// Adapter: Metrics
	// ----------------------------------------------------
	a.Metrics = metricsfx.NewMetricsProvider()

	err = a.Metrics.RegisterNativeCollectors()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInitFailed, err)
	}

	// ----------------------------------------------------
	// Adapter: Data
	// ----------------------------------------------------
	a.Data = datafx.NewRegistry(a.Logger)

	err = a.Data.LoadFromConfig(ctx, &a.Config.Data)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInitFailed, err)
	}

	// ----------------------------------------------------
	// Adapter: Queue
	// ----------------------------------------------------
	a.Queue = queuefx.NewRegistry(a.Logger)

	err = a.Queue.LoadFromConfig(ctx, &a.Config.Queue)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInitFailed, err)
	}

	// ----------------------------------------------------
	// Adapter: Arcade
	// ----------------------------------------------------
	a.Arcade = arcade.New(a.Config.Externals.Arcade)

	// ----------------------------------------------------
	// Adapter: Repository
	// ----------------------------------------------------
	a.Repository, err = storage.NewRepositoryFromDefault(a.Data)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInitFailed, err)
	}

	// ----------------------------------------------------
	// Business Services
	// ----------------------------------------------------
	a.ProfilesService = profiles.NewService(a.Logger, a.Repository)
	a.UsersService = users.NewService(a.Logger, a.Repository)
	a.UsersOAuthService = users.NewGitHubOAuthService(a.Logger, a.Repository)
	a.StoriesService = stories.NewService(a.Logger, a.Repository)

	return nil
}
