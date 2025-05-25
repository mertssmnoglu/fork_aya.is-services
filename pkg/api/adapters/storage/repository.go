package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/eser/ajan/datafx"
	"github.com/eser/aya.is-services/pkg/lib/caching"
)

const (
	DEFAULT_CACHE_TTL = 1 * time.Hour
)

var ErrDatasourceNotFound = errors.New("datasource not found")

type Repository struct {
	queries  *Queries
	cache    *caching.Cache
	cacheTTL time.Duration
}

func NewRepositoryFromDefault(dataRegistry *datafx.Registry) (*Repository, error) {
	datasource := dataRegistry.GetDefault()

	if datasource == nil {
		return nil, fmt.Errorf("%w: default", ErrDatasourceNotFound)
	}

	return NewRepositoryFromDataSource(datasource), nil
}

func NewRepositoryFromNamed(dataRegistry *datafx.Registry, name string) (*Repository, error) {
	datasource := dataRegistry.GetNamed(name)

	if datasource == nil {
		return nil, fmt.Errorf("%w: %s", ErrDatasourceNotFound, name)
	}

	return NewRepositoryFromDataSource(datasource), nil
}

func NewRepositoryFromDataSource(datasource datafx.Datasource) *Repository {
	db := datasource.GetConnection()

	repository := &Repository{ //nolint:exhaustruct
		queries:  &Queries{db: db},
		cacheTTL: DEFAULT_CACHE_TTL,
	}

	repository.cache = caching.NewCache(
		func(ctx context.Context, key string, target any) (bool, error) {
			cachedMessage, err := repository.CacheGetSince(
				ctx,
				key,
				time.Now().Add(-1*repository.cacheTTL),
			)
			if err != nil {
				return false, err
			}

			if cachedMessage == nil {
				return false, nil
			}

			unmarshallErr := json.Unmarshal(*cachedMessage, target)
			if unmarshallErr != nil {
				return false, unmarshallErr //nolint:wrapcheck
			}

			return true, nil
		},
		func(ctx context.Context, key string, value any) error {
			message, marshallErr := json.Marshal(value)
			if marshallErr != nil {
				return marshallErr //nolint:wrapcheck
			}

			err := repository.CacheSet(ctx, key, message)
			if err != nil {
				return err
			}

			return nil
		},
	)

	return repository
}
