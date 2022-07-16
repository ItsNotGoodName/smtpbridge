package server

import (
	"github.com/ItsNotGoodName/smtpbridge/config"
	"github.com/ItsNotGoodName/smtpbridge/core/auth"
	"github.com/ItsNotGoodName/smtpbridge/core/background"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
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
	backgrounds := []background.Background{}

	// Create stores
	dataStore := memdb.NewData()
	envelopeStore := memdb.NewEnvelope(dataStore)

	// Create services
	envelopeService := envelope.NewEnvelopeService(envelopeStore, dataStore)
	smtpAuthService := smtpAuthService(config)

	// Create SMTP server
	backgrounds = append(backgrounds, smtp.New(smtp.NewBackend(envelopeService, smtpAuthService), config.SMTP.Addr, config.SMTP.Size))

	// Create HTTP server
	if config.HTTP.Enable {
		backgrounds = append(backgrounds, router.New(config.HTTP.Addr, router.NewHandler(envelopeService)))
	}

	// Start server
	background.Run(interrupt.Context(), backgrounds)
}
