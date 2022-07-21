package server

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/config"
	"github.com/ItsNotGoodName/smtpbridge/core/auth"
	"github.com/ItsNotGoodName/smtpbridge/core/background"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/ItsNotGoodName/smtpbridge/core/event"
	"github.com/ItsNotGoodName/smtpbridge/left/controller"
	"github.com/ItsNotGoodName/smtpbridge/left/router"
	"github.com/ItsNotGoodName/smtpbridge/left/smtp"
	"github.com/ItsNotGoodName/smtpbridge/pkg/interrupt"
	"github.com/ItsNotGoodName/smtpbridge/right/memdb"
)

func smtpAuthService(config *config.Config) auth.Service {
	if config.SMTP.Auth {
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
		log.Fatalf("server.Start: invalid database or storage: '%s' '%s'", config.Database.Type, config.Storage.Type)
	}

	// Create stores
	dataStore := memdb.NewData(config.Storage.Memory.Limit, config.Storage.Memory.Size)
	envelopeStore := memdb.NewEnvelope(config.Database.Memory.Limit)

	// Create services
	pub := event.NewPub()
	envelopeService := event.NewEnvelopeService(envelope.NewEnvelopeService(envelopeStore, dataStore), pub)
	smtpAuthService := smtpAuthService(config)

	// Create HTTP server
	if config.HTTP.Enable {
		backgrounds = append(backgrounds, router.New(config.HTTP.Addr(), controller.New(envelopeService), dataStore.DataFS()))
	}

	// Create SMTP server
	if config.SMTP.Enable {
		backgrounds = append(backgrounds, smtp.New(smtp.NewBackend(envelopeService, smtpAuthService), config.SMTP.Addr(), config.SMTP.Size))
	}

	// Start background daemons
	background.Run(interrupt.Context(), backgrounds)

	log.Println("server.Start: stopped server")
}
