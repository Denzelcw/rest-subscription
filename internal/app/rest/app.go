package rest

import (
	"log/slog"
	"net/http"
	"os"
	"task_manager/internal/config"
	"task_manager/internal/http_server/handler"
	"task_manager/internal/http_server/middleware/logger"
	"task_manager/internal/lib/logger/sl"
	"task_manager/internal/storage/postgres"
	"task_manager/internal/usecases"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	log *slog.Logger
	cfg *config.Config
	srv *http.Server
}

func New(cfg *config.Config, log *slog.Logger) *App {
	storage, err := postgres.New(cfg.DbConfig)
	if err != nil {
		log.Error("failed to init storage: ", sl.Err(err))
		os.Exit(1)
	}

	subscriptionService := usecases.NewSubscriptionService(storage, log)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService, log)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/subscriptions", subscriptionHandler.AddSubscriptionHandler)
	router.Get("/subscriptions/{id}", subscriptionHandler.GetSubscriptionHandler)
	router.Get("/subscriptions", subscriptionHandler.GetListSubscriptionHandler)
	router.Delete("/subscriptions/{id}", subscriptionHandler.DeleteSubscriptionHandler)
	router.Put("/subscriptions/{id}", subscriptionHandler.UpdateSubscriptionHandler)
	router.Get("/subscriptions/total_cost", subscriptionHandler.GetTotalCostHandler)

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout * time.Second,
		WriteTimeout: cfg.HTTPServer.Timeout * time.Second,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout * time.Second,
	}

	return &App{
		log: log,
		cfg: cfg,
		srv: srv,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	a.log.Info("starting server", slog.String("address", a.cfg.Address))
	if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.log.Error("server error", slog.Any("err", err))
		return err
	}

	return nil
}

func (a *App) Stop() {
	a.log.Info("stopping server")
	if err := a.srv.Close(); err != nil {
		a.log.Error("failed to stop server", slog.Any("err", err))
	}
}
