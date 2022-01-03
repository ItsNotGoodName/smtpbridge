package endpoint

import (
	"fmt"
	"log"
	"sync"

	"github.com/ItsNotGoodName/smtpbridge/config"
	"github.com/ItsNotGoodName/smtpbridge/core"
)

type Repository struct {
	endpointMu  sync.RWMutex
	endpointMap map[string]core.EndpointPort
}

func NewRepository(cfg *config.Config) *Repository {
	r := Repository{
		endpointMap: make(map[string]core.EndpointPort),
		endpointMu:  sync.RWMutex{},
	}

	for _, c := range cfg.Endpoints {
		err := r.Create(c.Name, c.Type, c.Config)
		if err != nil {
			log.Fatalln("endpoint.NewRepository:", err)
		}
	}

	return &r
}

func (r *Repository) Get(name string) (core.EndpointPort, error) {
	r.endpointMu.RLock()
	defer r.endpointMu.RUnlock()

	endpoint, ok := r.endpointMap[name]
	if !ok {
		return nil, fmt.Errorf("%v: %s", core.ErrEndpointNotFound, name)
	}

	return endpoint, nil
}

func (r *Repository) Create(name, endpointType string, config map[string]string) error {
	r.endpointMu.Lock()
	defer r.endpointMu.Unlock()

	if _, ok := r.endpointMap[name]; ok {
		return fmt.Errorf("%v: %s", core.ErrEndpointNameConflict, name)
	}

	endpoint, err := factory(endpointType, config)
	if err != nil {
		return err
	}

	r.endpointMap[name] = endpoint

	return nil
}
