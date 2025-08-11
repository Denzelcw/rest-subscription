package rest

import (
	"fmt"
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
	subscriptionHandler := handler.NewUserSubscriptionHandler(subscriptionService, log, cfg.HTTPServer.Timeout)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/subscriptions", subscriptionHandler.AddUserSubscriptionHandler)
	router.Get("/subscriptions/{id}", subscriptionHandler.GetUserSubscriptionHandler)
	router.Get("/subscriptions", subscriptionHandler.GetListUserSubscriptionHandler)
	router.Delete("/subscriptions/{id}", subscriptionHandler.DeleteUserSubscriptionHandler)
	router.Put("/subscriptions", subscriptionHandler.UpdateSubscriptionHandler)
	router.Get("/subscriptions/total_cost", subscriptionHandler.GetTotalCostHandler)

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout * time.Second,
		WriteTimeout: cfg.Timeout * time.Second,
		IdleTimeout:  cfg.IdleTimeout * time.Second,
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
	url := fmt.Sprintf("http://%s/swagger/index.html", a.cfg.Address)
	a.log.Info("starting server", slog.String("url", url))

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
