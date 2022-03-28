package server

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/ItsNotGoodName/smtpbridge/config"
	"github.com/ItsNotGoodName/smtpbridge/core/app"
	"github.com/ItsNotGoodName/smtpbridge/core/attachment"
	"github.com/ItsNotGoodName/smtpbridge/core/auth"
	"github.com/ItsNotGoodName/smtpbridge/core/background"
	"github.com/ItsNotGoodName/smtpbridge/core/bridge"
	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/core/event"
	"github.com/ItsNotGoodName/smtpbridge/core/janitor"
	"github.com/ItsNotGoodName/smtpbridge/core/message"
	"github.com/ItsNotGoodName/smtpbridge/left/json"
	"github.com/ItsNotGoodName/smtpbridge/left/router"
	"github.com/ItsNotGoodName/smtpbridge/left/smtp"
	"github.com/ItsNotGoodName/smtpbridge/right/boltdb"
	"github.com/ItsNotGoodName/smtpbridge/right/endpoints"
	"github.com/ItsNotGoodName/smtpbridge/right/filedb"
	"github.com/ItsNotGoodName/smtpbridge/right/mockdb"
)

type Server struct {
	background []background.Background
}

func New(config *config.Config) Server {
	background := []background.Background{}

	// Init repositories
	endpointRepository := endpoints.NewRepository()
	var (
		attachmentDataRepository attachment.DataRepository
		attachmentRepository     attachment.Repository
		messageRepository        message.Repository
		eventRepository          event.Repository
	)
	if config.Database.IsMock() {
		attachmentDataRepository = mockdb.NewData()
		attachmentRepository = mockdb.NewAttachment()
		messageRepository = mockdb.NewMessage()
		eventRepository = mockdb.NewEvent()
	} else if config.Database.IsBolt() {
		attachmentDataRepository = filedb.NewData(config.Storage.AttachmentsPath)
		db := boltdb.NewDatabase(config.Database.BoltPath)
		att := boltdb.NewAttachment(&db, attachmentDataRepository)
		msg := boltdb.NewMessage(&db, attachmentDataRepository)

		background = append(background, db)
		attachmentRepository = att
		messageRepository = msg
		eventRepository = boltdb.NewEvent(&db)
	} else {
		log.Fatalln("server.New: unknown database type:", config.Database.Type)
	}

	// Create endpoints
	for _, e := range config.Endpoints {
		err := endpointRepository.Create(e.Name, e.Type, e.Config)
		if err != nil {
			log.Fatalln("server.New: could not create endpoint:", err)
		}
	}

	// Create bridges
	var bridges []bridge.Bridge
	for _, b := range config.Bridges {
		var endpoints []bridge.Endpoint
		for _, e := range b.Endpoints {
			facade, err := endpointRepository.Get(e.Name)
			if err != nil {
				log.Fatalln("server.New: could not find endpoint:", err)
			}

			endpoints = append(endpoints, bridge.NewEndpoint(facade, e.NoText, e.NoAttachments))
		}

		var filters []bridge.Filter
		for _, f := range b.Filters {
			filter, err := bridge.NewFilter(f.To, f.From, f.ToRegex, f.FromRegex)
			if err != nil {
				log.Fatalln("server.New: could not create filter:", err)
			}

			filters = append(filters, filter)
		}

		bridges = append(bridges, *bridge.New(b.Name, endpoints, filters, b.MinAttachments))
	}

	// Init services
	var smtpAuthService auth.Service
	if config.SMTP.Auth {
		smtpAuthService = auth.NewAccountService(config.SMTP.Username, config.SMTP.Password)
	} else {
		smtpAuthService = auth.NewAnonymousService()
	}
	eventService := event.NewEventService(eventRepository)
	endpointService := event.NewEndpointService(eventService, endpoint.NewEndpointService())
	messageService := event.NewMessageService(eventService, message.NewMessageService(messageRepository))
	attachmentService := attachment.NewAttachmentService(attachmentRepository)
	bridgeService := bridge.NewBridgeService(bridges, messageService, endpointService)
	if !config.Database.IsMock() {
		background = append(background, janitor.NewJanitorService(attachmentRepository, messageRepository, attachmentDataRepository, config.Storage.Size))
	}

	// Init app
	app := app.New(
		attachmentRepository,
		attachmentService,
		bridgeService,
		attachmentDataRepository,
		endpointService,
		eventRepository,
		messageRepository,
		messageService,
		smtpAuthService,
	)

	// Init smtp server
	smtpServer := smtp.New(smtp.NewBackend(app), config.SMTP.Addr, config.SMTP.Size)
	background = append(background, smtpServer)

	// Init http server
	if config.HTTP.Enable {
		router := router.New(app, json.New(), config.HTTP.Addr)
		background = append(background, router)
	}

	return Server{
		background: background,
	}
}

func (s Server) Run() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	done := make(chan struct{}, len(s.background))
	running := 0

	// Start background message service
	for _, background := range s.background {
		go background.Run(ctx, done)
		running++
	}

	// Wait for interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Stop background services
	cancel()

	// Wait for background serivces
	for i := 0; i < running; i++ {
		<-done
	}

	log.Println("server.Server.Run: stopped")
}
