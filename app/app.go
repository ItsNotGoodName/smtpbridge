package app

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/ItsNotGoodName/smtpbridge/core"
)

type App struct {
	background     []core.BackgroundPort
	database       core.DatabasePort
	attachmentREPO core.AttachmentRepositoryPort
	messageREPO    core.MessageRepositoryPort
	endpointREPO   core.EndpointRepositoryPort
	authSVC        core.AuthServicePort
	bridgeSVC      core.BridgeServicePort
	endpointSVC    core.EndpointServicePort
	messageSVC     core.MessageServicePort
}

func New(
	background []core.BackgroundPort,
	databasePort core.DatabasePort,
	attachmentREPO core.AttachmentRepositoryPort,
	messageREPO core.MessageRepositoryPort,
	endpointREPO core.EndpointRepositoryPort,
	authSVC core.AuthServicePort,
	bridgeSVC core.BridgeServicePort,
	endpointSVC core.EndpointServicePort,
	messageSVC core.MessageServicePort,
) *App {
	return &App{
		background,
		databasePort,
		attachmentREPO,
		messageREPO,
		endpointREPO,
		authSVC,
		bridgeSVC,
		endpointSVC,
		messageSVC,
	}
}

func (a *App) Run() {
	done := make(chan struct{})
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	running := 0

	// Start background message service
	for _, background := range a.background {
		go background.Run(ctx, done)
		running++
	}

	// Wait for interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Stop background services
	cancel()

	// Close database
	if err := a.database.Close(); err != nil {
		log.Println("app.App.Run: could not close database:", err)
	}

	// Wait for background serivces
	for i := 0; i < running; i++ {
		<-done
	}

	log.Println("app.App.Run: stopped")
}
