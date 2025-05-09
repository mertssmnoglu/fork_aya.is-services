package storage

import (
	"errors"
	"fmt"

	"github.com/eser/ajan/datafx"
)

var ErrDatasourceNotFound = errors.New("datasource not found")

type Repository struct {
	queries *Queries
}

func NewRepositoryFromDefault(dataRegistry *datafx.Registry) (*Repository, error) {
	datasource := dataRegistry.GetDefault()

	if datasource == nil {
		return nil, fmt.Errorf("%w - default", ErrDatasourceNotFound)
	}

	db := datasource.GetConnection()

	repository := &Repository{
		queries: &Queries{db: db},
	}

	return repository, nil
}

func NewRepositoryFromNamed(dataRegistry *datafx.Registry, name string) (*Repository, error) {
	datasource := dataRegistry.GetNamed(name)

	if datasource == nil {
		return nil, fmt.Errorf("%w - %s", ErrDatasourceNotFound, name)
	}

	db := datasource.GetConnection()

	repository := &Repository{
		queries: &Queries{db: db},
	}

	return repository, nil
}
