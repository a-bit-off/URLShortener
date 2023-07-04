package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"golang.org/x/exp/slog"

	"URLShortener/internal/config"
	"URLShortener/internal/http-server/handlers/delete"
	"URLShortener/internal/http-server/handlers/redirect"
	"URLShortener/internal/http-server/handlers/url/save"
	mwLogger "URLShortener/internal/http-server/middleware/logger"
	"URLShortener/internal/lib/logger/hadlers/slogpretty"
	"URLShortener/internal/lib/logger/sl"
	"URLShortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// init config: cleanenv
	cfg := initConfig()

	// init logger: slog
	log := initLogger(cfg)

	// init storage: sqlite
	storage := initStorage(cfg, log)

	// init router: chi, chi render
	router := chi.NewRouter()

	// init middleware: chi Mux, middleware
	initMiddleware(router, log)

	// init handlers
	initHandlers(router, log, storage)

	// run server
	runServer(cfg, router, log)
}

func initConfig() *config.Config {
	configPath := flag.String("CONFIG_PATH", "", "path to config")
	flag.Parse()
	return config.MustLoad(*configPath)
}

func initLogger(cfg *config.Config) *slog.Logger {
	log := setupLogger(cfg.Env)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")
	return log
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
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

func initStorage(cfg *config.Config, log *slog.Logger) *sqlite.Storage {
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	return storage
}

func initMiddleware(router *chi.Mux, log *slog.Logger) {
	router.Use(middleware.RequestID) //
	router.Use(middleware.Logger)    //
	router.Use(mwLogger.New(log))    //
	router.Use(middleware.Recoverer) // обработка panic в handler
	router.Use(middleware.URLFormat) // обработка url
}

func initHandlers(router *chi.Mux, log *slog.Logger, storage *sqlite.Storage) {
	router.Post("/url", save.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))
	router.Delete("/url/{alias}", delete.New(log, storage))
}

func runServer(cfg *config.Config, router *chi.Mux, log *slog.Logger) {
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
