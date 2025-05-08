package appcontext

import (
	"github.com/eser/ajan"
	"github.com/eser/aya.is-services/pkg/api/adapters/arcade"
)

type FeatureFlags struct {
	Dummy bool `conf:"DUMMY" default:"false"` // dummy feature flag
}

type Externals struct {
	Arcade arcade.Config `conf:"ARCADE"`
}

type AppConfig struct {
	ajan.BaseConfig

	Features FeatureFlags `conf:"FEATURES"`

	Externals Externals `conf:"EXTERNALS"`
}
