package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.com/mildd/flow-gateway/internal/adapters/secret"
	secretCmd "gitlab.com/mildd/flow-gateway/internal/app/commands/secret"
	secret2 "gitlab.com/mildd/flow-gateway/internal/app/query/secret"
	"net"
	"net/http"

	"gitlab.com/mildd/flow-gateway/internal/app/graceful"
	"gitlab.com/mildd/flow-gateway/internal/app/handlers"
	"gitlab.com/mildd/flow-gateway/internal/configs"
	"gitlab.com/mildd/flow-gateway/pkg/logging"
	"gitlab.com/mildd/flow-gateway/pkg/metric"
	"golang.org/x/sync/errgroup"

	"github.com/gin-gonic/gin"
)

type App struct {
	cfg *configs.Config

	router            *gin.Engine
	httpServer        *http.Server
	Options           *Options
	Commands          Commands
	Queries           Queries
	ShutdownFunctions []func() error
}

type Commands struct {
	GetSecret secretCmd.GetHandler
}

type Queries struct {
	GetSecret secret2.SecretHandler
}

func NewApp(ctx context.Context, config *configs.Config, options *Options) (*App, error) {
	logging.Info(ctx, "router initializing")

	router := gin.Default()

	dbURL := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s database=%s sslmode=disable",
		config.PostgreSQL.Username, config.PostgreSQL.Password, config.PostgreSQL.Host, config.PostgreSQL.Port, config.PostgreSQL.Database,
	)
	logging.Info(ctx, "handlers initializing")

	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		return &App{}, errors.New("failed to open database")
	}
	defer func() {
		_ = db.Close()
	}()

	db.SetMaxIdleConns(30)

	app, err := NewApplication(ctx, &Options{
		DB: db,
	})
	if err != nil {
		return &App{}, errors.New("failed to initialize application")
	}

	handler := handlers.NewHandler(db, config.Secret)
	handler.RegisterRoutes()

	logging.Info(ctx, "heartbeat metric initializing")

	metricHandler := metric.Handler{}
	metricHandler.Register(router)

	return &App{
		cfg:    config,
		router: router,
	}, nil

}

func (a *App) NewApplication(_ context.Context, options *Options) (*App, error) {
	secretRepo := secret.NewPGRepository(options.DB)
	return &App{
		Options: options,
		Commands: Commands{
			// Secret
			GetSecret: secretCmd.NewGetHandler(secretRepo),
		},
		Queries: Queries{
			GetSecret: secret2.NewSecretHandler(secretRepo),
		},
		ShutdownFunctions: []func() error{},
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return a.startHTTP(ctx)
	})
	return grp.Wait()
}

func (a *App) startHTTP(ctx context.Context) error {
	logger := logging.WithFields(ctx, map[string]interface{}{
		"IP":   a.cfg.HTTP.IP,
		"Port": a.cfg.HTTP.Port,
	})
	logger.Info("HTTP Server initializing")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.HTTP.IP, a.cfg.HTTP.Port))
	if err != nil {
		logger.WithError(err).Fatal("failed to create listener")
	}

	handler := a.router

	a.httpServer = &http.Server{
		Handler:      handler,
		WriteTimeout: a.cfg.HTTP.WriteTimeout,
		ReadTimeout:  a.cfg.HTTP.ReadTimeout,
	}

	if err = a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Warning("server shutdown")
		default:
			logger.Fatal(err)
		}
	}

	httpErrChan := make(chan error, 1)
	httpShutdownChan := make(chan struct{})

	graceful.PerformGracefulShutdown(a.httpServer, httpErrChan, httpShutdownChan)

	return err
}
