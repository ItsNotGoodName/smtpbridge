package server

import (
	"context"
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
	"github.com/ItsNotGoodName/smtpbridge/right/boltdb"
	"github.com/ItsNotGoodName/smtpbridge/right/filedb"
	"github.com/ItsNotGoodName/smtpbridge/right/memdb"
)

func smtpAuthService(config *config.Config) auth.Service {
	if config.SMTP.Auth() {
		return auth.NewAccountService(config.SMTP.Username, config.SMTP.Password)
	} else {
		return auth.NewAnonymousService()
	}
}

func Start(ctx context.Context, config *config.Config) <-chan struct{} {
	// Background daemons
	backgrounds := []background.Background{}

	// Create stores
	var dataStore envelope.DataStore
	if config.Storage.IsMemory() {
		dataStore = memdb.NewData(config.Storage.Memory.Size)
	} else if config.Storage.IsFile() {
		dataStore = filedb.NewData(config.Storage.File.Path)
	} else {
		log.Fatalln("server.Start: storage invalid:", config.Storage.Type)
	}
	var envelopeStore envelope.Store
	if config.Database.IsMemory() {
		envelopeStore = memdb.NewEnvelope(config.Database.Memory.Limit)
	} else if config.Database.IsBolt() {
		boltdb := boltdb.NewDatabase(config.Database.Bolt.File)
		backgrounds = append(backgrounds, boltdb)
		envelopeStore = boltdb
	} else {
		log.Fatalln("server.Start: database invalid:", config.Database.Type)
	}

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
			Name:               end.Name,
			Type:               end.Type,
			Config:             end.Config,
			TextDisable:        end.TextDisable,
			TextTemplate:       end.TextTemplate,
			AttachmentsDisable: end.AttachmentsDisable,
		}); err != nil {
			log.Fatalf("server.Start: endpoint '%s': %s", end.Name, err)
		}

		log.Printf("server.Start: created endpoint '%s'", end.Name)
	}

	// Create bridges from config
	for i, brid := range config.Bridges {
		// Filters request
		var filtersReq []bridge.CreateFilterRequest
		filtersLength := len(brid.Filters)
		if filtersLength > 0 {
			filtersReq = make([]bridge.CreateFilterRequest, 0, filtersLength)
			for _, filter := range brid.Filters {
				filtersReq = append(filtersReq, bridge.CreateFilterRequest{
					From:          filter.From,
					To:            filter.To,
					FromRegex:     filter.FromRegex,
					ToRegex:       filter.ToRegex,
					MatchTemplate: filter.MatchTemplate,
				})
			}
		} else {
			filtersReq = []bridge.CreateFilterRequest{{
				From:          brid.From,
				To:            brid.To,
				FromRegex:     brid.FromRegex,
				ToRegex:       brid.ToRegex,
				MatchTemplate: brid.MatchTemplate,
			}}
		}

		// Create bridge
		if err := bridgeService.CreateBridge(&bridge.CreateBridgeRequest{
			Endpoints: brid.Endpoints,
			Filters:   filtersReq,
		}); err != nil {
			log.Fatalf("server.Start: bridge '%d': %s", i, err)
		}

		log.Printf("server.Start: bridge %d created with %d filters and %d endpoints", i, len(filtersReq), len(brid.Endpoints))
	}

	// Create HTTP server
	if !config.HTTP.Disable {
		backgrounds = append(backgrounds, http.New(
			config.HTTP.Addr(),
			dataStore.(envelope.LocalDataStore),
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
	return background.Run(ctx, backgrounds)
}
