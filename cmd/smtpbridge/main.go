package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/config"
	"github.com/ItsNotGoodName/smtpbridge/cron"
	"github.com/ItsNotGoodName/smtpbridge/internal/app"
	"github.com/ItsNotGoodName/smtpbridge/internal/build"
	"github.com/ItsNotGoodName/smtpbridge/internal/bus"
	"github.com/ItsNotGoodName/smtpbridge/internal/database"
	"github.com/ItsNotGoodName/smtpbridge/internal/file"
	"github.com/ItsNotGoodName/smtpbridge/internal/mailman"
	"github.com/ItsNotGoodName/smtpbridge/internal/repo"
	"github.com/ItsNotGoodName/smtpbridge/migrations"
	"github.com/ItsNotGoodName/smtpbridge/pkg/secret"
	"github.com/ItsNotGoodName/smtpbridge/pkg/sutureext"
	"github.com/ItsNotGoodName/smtpbridge/smtp"
	"github.com/ItsNotGoodName/smtpbridge/web"
	"github.com/ItsNotGoodName/smtpbridge/web/helpers"
	"github.com/ItsNotGoodName/smtpbridge/web/http"
	"github.com/ItsNotGoodName/smtpbridge/web/session"
	"github.com/Rican7/lieut"
	"github.com/gorilla/sessions"
	"github.com/reugn/go-quartz/quartz"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

func main() {
	ctx := context.Background()

	flags := config.WithFlagSet(flag.NewFlagSet(os.Args[0], flag.ExitOnError))

	app := lieut.NewSingleCommandApp(
		lieut.AppInfo{
			Name:    "smtpbridge",
			Version: build.Current.Version,
			Summary: "Bridge email to other messaging services.",
		},
		run(flags),
		flags,
		os.Stdout,
		os.Stderr,
	)

	code := app.Run(ctx, os.Args[1:])

	os.Exit(code)
}

func run(flags *flag.FlagSet) lieut.Executor {
	return func(ctx context.Context, arguments []string) error {
		// Config
		parser, err := config.NewParser(flags)
		if err != nil {
			return err
		}
		raw, err := parser.Read()
		if err != nil {
			return err
		}
		cfg, err := parser.Parse(raw)
		if err != nil {
			return err
		}
		if cfg.Debug {
			fmt.Println(helpers.JSON(cfg))
		}

		// Database
		db, err := database.New(cfg.DatabasePath, cfg.Debug)
		if err != nil {
			return err
		}
		// Database migrations
		if err := migrations.Migrate(db); err != nil {
			return err
		}
		// Sync database with config
		err = repo.InternalSync(ctx, db, cfg.InternalEndpoints, cfg.InternalRules, cfg.InternalRuleToEndpoint)
		if err != nil {
			return err
		}

		// File store
		fileStore := file.NewStore(cfg.AttachmentsDirectory)

		// Bus
		bus, err := bus.New()
		if err != nil {
			return err
		}

		// App
		webTestFileStore := app.NewWebTestFileStore("apple-touch-icon.png", fmt.Sprintf("http://127.0.0.1:%d", cfg.HTTPPort))
		app, release := app.New(db, fileStore, bus, cfg.Config, cfg.EndpointFactory, webTestFileStore)
		defer release()

		// Supervisor
		super := suture.New("root", suture.Spec{
			EventHook: sutureext.EventHook(),
		})

		// Cron
		scheduler := cron.NewScheduler()
		super.Add(scheduler)

		// Cron schedule
		super.Add(cron.ScheduleJob(func(ctx context.Context) error {
			{
				job := cron.NewRetentionPolicy(app)
				trigger := quartz.NewSimpleTrigger(cfg.Config.RetentionPolicy.MinAge)
				if err := scheduler.ScheduleJob(ctx, job, trigger); err != nil {
					return err
				}

				job.Execute(ctx)
			}

			{
				job := cron.NewAttachmentOrphan(app)
				trigger := quartz.NewSimpleTrigger(30 * time.Minute)
				if err := scheduler.ScheduleJob(ctx, job, trigger); err != nil {
					return err
				}

				job.Execute(ctx)
			}

			if cfg.HealthcheckURL != "" {
				job := cron.NewHealthcheck(app, cfg.HealthcheckURL)
				trigger := quartz.NewSimpleTrigger(cfg.HealthcheckInterval)
				if err := scheduler.ScheduleJob(ctx, job, trigger); err != nil {
					return err
				}

				if cfg.HealthcheckStartup {
					job.Execute(ctx)
				}
			}

			{
				job := cron.NewDatabaseVacuum(app)
				trigger := quartz.NewSimpleTrigger(24 * time.Hour)
				if err := scheduler.ScheduleJob(ctx, job, trigger); err != nil {
					return err
				}

				job.Execute(ctx)
			}

			return nil
		}))

		for i := 1; i <= cfg.MailmanWorkers; i++ {
			// Mailman
			mailman := mailman.New(i, app, bus, fileStore, cfg.EndpointFactory)
			super.Add(mailman)
		}

		// SMTP
		if !cfg.SMTPDisable {
			backend := smtp.NewBackend(app)
			smtp := smtp.New(backend, cfg.SMTPAddress, cfg.SMTPMaxMessageBytes)
			super.Add(smtp)
		}

		// HTTP
		if !cfg.HTTPDisable {
			// HTTP CSRF
			csrfSecret, err := secret.GetOrCreate(cfg.CSRFSecretPath)
			if err != nil {
				return err
			}

			// HTTP session
			sessionSecret, err := secret.GetOrCreate(cfg.SessionSecretPath)
			if err != nil {
				return err
			}
			sessionStore := sessions.NewFilesystemStore(cfg.SessionsDirectory, sessionSecret)
			session.ConfigureOptions(sessionStore.Options)

			controller := http.NewController(app, cfg.TimeHourFormat)
			router := http.NewRouter(controller, app, fileStore, csrfSecret, sessionStore)
			server := http.NewServer(router, cfg.HTTPAddress)
			super.Add(server)
		}

		// // IMAP
		// if !cfg.IMAPDisable {
		// 	backend := imap.NewBackend(app)
		// 	server := imap.NewServer(backend, cfg.IMAPAddress)
		// 	super.Add(server)
		// }

		// Web
		if web.DevMode {
			web := web.NewRefresher()
			super.Add(web)
		}

		return super.Serve(ctx)
	}
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
