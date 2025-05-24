package http

import (
	"context"

	"github.com/eser/ajan/datafx"
	"github.com/eser/ajan/httpfx"
	"github.com/eser/ajan/httpfx/middlewares"
	"github.com/eser/ajan/httpfx/modules/healthcheck"
	"github.com/eser/ajan/httpfx/modules/openapi"
	"github.com/eser/ajan/httpfx/modules/profiling"
	"github.com/eser/ajan/logfx"
	"github.com/eser/ajan/metricsfx"
)

func Run(
	ctx context.Context,
	config *httpfx.Config,
	metricsProvider *metricsfx.MetricsProvider,
	logger *logfx.Logger,
	dataRegistry *datafx.Registry,
) (func(), error) {
	routes := httpfx.NewRouter("/")
	httpService := httpfx.NewHttpService(config, routes, metricsProvider, logger)

	// http middlewares
	routes.Use(middlewares.ErrorHandlerMiddleware())
	routes.Use(middlewares.ResolveAddressMiddleware())
	routes.Use(middlewares.ResponseTimeMiddleware())
	routes.Use(middlewares.CorrelationIdMiddleware())
	routes.Use(middlewares.CorsMiddleware())
	routes.Use(middlewares.MetricsMiddleware(httpService.InnerMetrics))
	// routes.Use(AuthMiddleware(dataRegistry))

	// http modules
	healthcheck.RegisterHttpRoutes(routes, config)
	openapi.RegisterHttpRoutes(routes, config)
	profiling.RegisterHttpRoutes(routes, config)

	// --- OAuth Service wiring ---
	githubOAuthService := NewGitHubOAuthService(dataRegistry)

	// http routes
	RegisterHttpRoutesForUsers(routes, logger, dataRegistry, githubOAuthService) //nolint:contextcheck
	RegisterHttpRoutesForSite(routes, logger, dataRegistry)                      //nolint:contextcheck
	RegisterHttpRoutesForProfiles(routes, logger, dataRegistry)                  //nolint:contextcheck
	RegisterHttpRoutesForStories(routes, logger, dataRegistry)                   //nolint:contextcheck

	// run
	return httpService.Start(ctx) //nolint:wrapcheck
}
