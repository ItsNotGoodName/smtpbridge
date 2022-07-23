package memdb

import (
	"sort"
	"sync"

	"github.com/ItsNotGoodName/smtpbridge/core"
	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
)

type Endpoint struct {
	endpointsMu sync.Mutex
	endpoints   map[string]endpoint.Endpoint
}

func NewEndpoint() *Endpoint {
	return &Endpoint{
		endpoints: make(map[string]endpoint.Endpoint),
	}
}

func (e *Endpoint) CreateEndpoint(end endpoint.Endpoint) error {
	e.endpointsMu.Lock()
	_, ok := e.endpoints[end.Name]
	if ok {
		e.endpointsMu.Unlock()
		return core.ErrEndpointNameExists
	}
	e.endpoints[end.Name] = end
	e.endpointsMu.Unlock()

	return nil
}

func (e *Endpoint) GetEndpoint(name string) (endpoint.Endpoint, error) {
	e.endpointsMu.Lock()
	end, ok := e.endpoints[name]
	if !ok {
		e.endpointsMu.Unlock()
		return endpoint.Endpoint{}, core.ErrEndpointNotFound
	}
	e.endpointsMu.Unlock()

	return end, nil
}

func (e *Endpoint) ListEndpoint() []endpoint.Endpoint {
	ends := []endpoint.Endpoint{}

	e.endpointsMu.Lock()
	for _, end := range e.endpoints {
		ends = append(ends, end)
	}
	e.endpointsMu.Unlock()

	sort.Slice(ends, func(i, j int) bool {
		return ends[i].Name < ends[j].Name
	})

	return ends
}
