package endpoint

import (
	"fmt"
	"log"
	"sync"

	"github.com/ItsNotGoodName/smtpbridge/domain"
)

type Repository struct {
	endpointMu  sync.RWMutex
	endpointMap map[string]domain.EndpointPort
}

func NewRepository(configEndpoints []domain.ConfigEndpoint) *Repository {
	r := Repository{
		endpointMap: make(map[string]domain.EndpointPort),
		endpointMu:  sync.RWMutex{},
	}

	for _, c := range configEndpoints {
		err := r.Create(c.Name, c.Type, c.Config)
		if err != nil {
			log.Fatalln("endpoint.NewRepository:", err)
		}
	}

	return &r
}

func (r *Repository) Get(name string) (domain.EndpointPort, error) {
	r.endpointMu.RLock()
	defer r.endpointMu.RUnlock()

	endpoint, ok := r.endpointMap[name]
	if !ok {
		return nil, fmt.Errorf("%s: %v", name, domain.ErrEndpointNotFound)
	}

	return endpoint, nil
}

func (r *Repository) Create(name, endpointType string, config map[string]string) error {
	r.endpointMu.Lock()
	defer r.endpointMu.Unlock()

	if _, ok := r.endpointMap[name]; ok {
		return fmt.Errorf("%s: %v", name, domain.ErrEndpointNameConflict)
	}

	endpoint, err := factory(endpointType, config)
	if err != nil {
		return err
	}

	r.endpointMap[name] = endpoint

	return nil
}
