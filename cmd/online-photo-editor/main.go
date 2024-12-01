package main

import (
	"log/slog"
	"net/http"
	"online-photo-editor/cmd/internal/config"
	"online-photo-editor/cmd/internal/http-server/handlers/img/save"
	mwLogger "online-photo-editor/cmd/internal/http-server/middleware/logger"
	"online-photo-editor/cmd/internal/lib/logger/handlers/slogpretty"
	"online-photo-editor/cmd/internal/lib/logger/sl"
	imgStorage "online-photo-editor/cmd/internal/storage/img"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting online-photo-editor", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	imageStorage, err := imgStorage.New(cfg.StorageImagePath)
	if err != nil {
		log.Error("failed to init image storage", sl.Err(err))
		os.Exit(1)
	}
	_ = imageStorage

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/img", save.New(log, imageStorage))

	fileServer := http.FileServer(http.Dir(cfg.StorageImagePath))
	router.Handle("/images/*", http.StripPrefix("/images", fileServer))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
