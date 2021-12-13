package app

// NewEndpoints creates a list of Endpoints from config file and factory.
func NewEndpoints(config []ConfigEndpoint, factory func(senderType string, config map[string]string) (EndpointPort, error)) (map[string]EndpointPort, error) {
	endpoints := make(map[string]EndpointPort)
	for _, c := range config {
		sender, err := factory(c.Type, c.Config)
		if err != nil {
			return nil, err
		}
		endpoints[c.Name] = sender
	}

	return endpoints, nil
}
