package main

import (
	"context"
	"os"

	"github.com/ItsNotGoodName/smtpbridge/config"
	"github.com/ItsNotGoodName/smtpbridge/internal/build"
	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/db"
	"github.com/ItsNotGoodName/smtpbridge/internal/procs"
	"github.com/ItsNotGoodName/smtpbridge/pkg/background"
	"github.com/ItsNotGoodName/smtpbridge/pkg/interrupt"
	"github.com/ItsNotGoodName/smtpbridge/smtp"
	"github.com/ItsNotGoodName/smtpbridge/web/http"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, shutdown := context.WithCancel(interrupt.Context())
	defer shutdown()

	<-run(ctx, shutdown)
}

func run(ctx context.Context, shutdown context.CancelFunc) <-chan struct{} {
	raw, err := config.Read()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config")
	}

	cfg, err := config.Parse(raw)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse config")
	}

	// Database
	bunDB, err := db.New(cfg.DatabasePath)
	if err != nil {
		log.Fatal().Err(err).Str("path", cfg.DatabasePath).Msg("Failed to open database")
	}
	if err := db.Migrate(ctx, bunDB); err != nil {
		log.Fatal().Err(err).Str("path", cfg.DatabasePath).Msg("Failed to migrate database")
	}

	// File Store
	fileStore := core.NewFileStore(cfg.AttachmentsDirectory)

	// App
	app := core.NewApp(bunDB, fileStore)
	if err := procs.InternalSync(app.Context(ctx), cfg.Endpoints, cfg.Rules, cfg.RuleEndpoints); err != nil {
		log.Fatal().Err(err).Msg("Failed to sync app from config")
	}

	procs.MailmanStart(ctx, app)
	procs.GardenerStart(ctx, app, cfg.RetentionPolicy)
	procs.VacuumStart(ctx, app)

	// SMTP
	smtp := smtp.New(app, shutdown, cfg.SMTPAddress, cfg.SMTPMaxMessageBytes)

	// HTTP
	http := http.New(app, shutdown, cfg.HTTPAddress, cfg.HTTPBodyLimit, cfg.RetentionPolicy)

	// Start
	return background.Run(ctx,
		smtp,
		http,
	)
}

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	build.Current = build.Build{
		BuiltBy: builtBy,
		Commit:  commit,
		Date:    date,
		Version: version,
	}
}
