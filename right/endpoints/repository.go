package endpoints

import (
	"fmt"
	"sync"

	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
)

type Repository struct {
	endpointMu  sync.RWMutex
	endpointMap map[string]endpoint.Endpoint
}

func NewRepository() *Repository {
	return &Repository{
		endpointMu:  sync.RWMutex{},
		endpointMap: make(map[string]endpoint.Endpoint),
	}
}

func (r *Repository) Get(name string) (*endpoint.Facade, error) {
	r.endpointMu.RLock()
	defer r.endpointMu.RUnlock()

	end, ok := r.endpointMap[name]
	if !ok {
		return nil, fmt.Errorf("%v: %w", name, endpoint.ErrNotFound)
	}

	return endpoint.NewFacade(name, end), nil
}

func (r *Repository) Create(name, endpointType string, config map[string]string) error {
	r.endpointMu.Lock()
	defer r.endpointMu.Unlock()

	if _, ok := r.endpointMap[name]; ok {
		return fmt.Errorf("%v: %w", name, endpoint.ErrNameConflict)
	}

	end, err := factory(endpointType, config)
	if err != nil {
		return fmt.Errorf("%v: %w", name, err)
	}

	r.endpointMap[name] = end

	return nil
}
