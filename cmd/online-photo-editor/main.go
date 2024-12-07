package main

import (
	"context"
	"log/slog"
	"net/http"
	"online-photo-editor/internal/config"
	"online-photo-editor/internal/http-server/handlers/image/blur"
	"online-photo-editor/internal/http-server/handlers/image/brightness"
	"online-photo-editor/internal/http-server/handlers/image/contrast"
	"online-photo-editor/internal/http-server/handlers/image/convert"
	"online-photo-editor/internal/http-server/handlers/image/crop"
	"online-photo-editor/internal/http-server/handlers/image/gamma"
	"online-photo-editor/internal/http-server/handlers/image/processor"
	"online-photo-editor/internal/http-server/handlers/image/resize"
	"online-photo-editor/internal/http-server/handlers/image/saturation"
	"online-photo-editor/internal/http-server/handlers/image/sharpen"
	"online-photo-editor/internal/http-server/handlers/image/upload"
	mwLogger "online-photo-editor/internal/http-server/middleware/logger"
	"online-photo-editor/internal/lib/logger/handlers/slogpretty"
	"online-photo-editor/internal/lib/logger/sl"
	imgStorage "online-photo-editor/internal/storage/filesystem"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	router := setupRouter(log, imageStorage, cfg.StorageImagePath)

	log.Info("starting server", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server", sl.Err(err))
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	log.Info("server stopped")

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

func setupRouter(log *slog.Logger, imageStorage *imgStorage.ImageStorage, storagePath string) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.RealIP, mwLogger.New(log), middleware.Recoverer, middleware.URLFormat)

	router.Post("/image", upload.New(log, imageStorage))

	router.Post("/image/crop", crop.New(log, imageStorage))

	router.Post("/image/resize", resize.New(log, imageStorage))

	router.Post("/image/convert", convert.New(log, imageStorage))

	router.Post("/image/blur", blur.New(log, imageStorage))

	router.Post("/image/brightness", brightness.New(log, imageStorage))

	router.Post("/image/contrast", contrast.New(log, imageStorage))

	router.Post("/image/gamma", gamma.New(log, imageStorage))

	router.Post("/image/saturation", saturation.New(log, imageStorage))

	router.Post("/image/sharpen", sharpen.New(log, imageStorage))

	router.Post("/image/process", processor.New(log, imageStorage))

	fileServer := http.FileServer(http.Dir(storagePath))
	router.Handle("/images/*", http.StripPrefix("/images", fileServer))

	return router
}
