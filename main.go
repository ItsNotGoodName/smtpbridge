package main

import (
	"context"
	"fmt"
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
	cli := config.ReadAndParseCLI()

	if cli.Version {
		fmt.Println(version)
		return
	}

	ctx, shutdown := context.WithCancel(interrupt.Context())
	defer shutdown()

	<-run(ctx, shutdown, cli)
}

func run(ctx context.Context, shutdown context.CancelFunc, cli config.CLI) <-chan struct{} {
	raw, err := config.Read(cli)
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

	procs.MailmanBackground(ctx, app)
	procs.GardenerBackground(ctx, app, cfg.RetentionPolicy)
	procs.VacuumStart(ctx, app)

	var backgrounds []background.Background

	// SMTP
	if !cfg.SMTPDisable {
		smtp := smtp.New(app, shutdown, cfg.SMTPAddress, cfg.SMTPMaxMessageBytes)
		backgrounds = append(backgrounds, smtp)
	}

	// HTTP
	if !cfg.HTTPDisable {
		http := http.New(app, shutdown, cfg.HTTPAddress, cfg.HTTPBodyLimit, cfg.RetentionPolicy)
		backgrounds = append(backgrounds, http)
	}

	// Start
	return background.Run(ctx, backgrounds...)
}

var (
	builtBy    = "unknown"
	commit     = ""
	date       = ""
	version    = "dev"
	repoURL    = "https://github.com/ItsNotGoodName/smtpbridge"
	releaseURL = ""
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	build.Current = build.Build{
		BuiltBy:    builtBy,
		Commit:     commit,
		Date:       date,
		Version:    version,
		RepoURL:    repoURL,
		ReleaseURL: releaseURL,
	}
}
