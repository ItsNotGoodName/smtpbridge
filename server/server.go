package server

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/config"
	"github.com/ItsNotGoodName/smtpbridge/core/auth"
	"github.com/ItsNotGoodName/smtpbridge/core/background"
	"github.com/ItsNotGoodName/smtpbridge/core/bridge"
	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/ItsNotGoodName/smtpbridge/core/event"
	"github.com/ItsNotGoodName/smtpbridge/left/http"
	"github.com/ItsNotGoodName/smtpbridge/left/smtp"
	"github.com/ItsNotGoodName/smtpbridge/pkg/interrupt"
	"github.com/ItsNotGoodName/smtpbridge/right/memdb"
)

func smtpAuthService(config *config.Config) auth.Service {
	if config.SMTP.Auth() {
		return auth.NewAccountService(config.SMTP.Username, config.SMTP.Password)
	} else {
		return auth.NewAnonymousService()
	}
}

func Start(config *config.Config) {
	log.Println("server.Start: starting server")

	// Background daemons
	backgrounds := []background.Background{}

	// Only memdb database and storage is supported
	if !(config.Database.IsMemDB() && config.Storage.IsMemDB()) {
		log.Fatalf("server.Start: invalid database or storage: '%s' or '%s'", config.Database.Type, config.Storage.Type)
	}

	// Create stores
	dataStore := memdb.NewData(config.Storage.Memory.Size)
	envelopeStore := memdb.NewEnvelope(config.Database.Memory.Limit)

	// Create services
	pub := event.NewPub()
	envelopeService := event.NewEnvelopeService(envelope.NewEnvelopeService(envelopeStore, dataStore), pub)
	smtpAuthService := smtpAuthService(config)
	endpointService := endpoint.NewEndpointService(memdb.NewEndpoint())
	bridgeService := bridge.NewBridgeService(pub, envelopeService, endpointService)
	backgrounds = append(backgrounds, bridgeService)

	// Create endpoints from config
	for _, end := range config.Endpoints {
		if err := endpointService.CreateEndpoint(endpoint.CreateEndpointRequest{
			Name:     end.Name,
			Type:     end.Type,
			Config:   end.Config,
			Template: end.TextTemplate,
		}); err != nil {
			log.Fatalf("server.Start: endpoint '%s': %s", end.Name, err)
		}
	}

	// Create bridges from config
	for i, brid := range config.Bridges {
		if err := bridgeService.CreateBridge(&bridge.CreateBridgeRequest{
			From:          brid.From,
			To:            brid.To,
			FromRegex:     brid.FromRegex,
			ToRegex:       brid.ToRegex,
			Endpoints:     brid.Endpoints,
			MatchTemplate: brid.MatchTemplate,
		}); err != nil {
			log.Fatalf("server.Start: bridge '%d': %s", i, err)
		}
	}

	// Create HTTP server
	if !config.HTTP.Disable {
		backgrounds = append(backgrounds, http.New(
			config.HTTP.Addr(),
			dataStore.DataFS(),
			envelopeService,
			endpointService,
		))
	}

	// Create SMTP server
	if !config.SMTP.Disable {
		backgrounds = append(backgrounds, smtp.New(
			smtp.NewBackend(envelopeService, smtpAuthService),
			config.SMTP.Addr(),
			config.SMTP.Size,
		))
	}

	// Start background daemons
	background.Run(interrupt.Context(), backgrounds)

	log.Println("server.Start: stopped server")
}
