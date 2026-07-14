package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	engine "github.com/squeakycheese75/open-incentives-engine"
	"github.com/squeakycheese75/open-incentives/configs"
	"github.com/squeakycheese75/open-incentives/internal/adapters"
	"github.com/squeakycheese75/open-incentives/internal/evaluate"

	"github.com/squeakycheese75/open-incentives/internal/services"

	"github.com/squeakycheese75/open-incentives/internal/admin"
	"github.com/squeakycheese75/open-incentives/internal/admin/auth"
	usecase_auth "github.com/squeakycheese75/open-incentives/internal/admin/auth/usecase"
	usecase_admin "github.com/squeakycheese75/open-incentives/internal/admin/usecase"
	"github.com/squeakycheese75/open-incentives/internal/db/sqlitedb"
	usecase_eval "github.com/squeakycheese75/open-incentives/internal/evaluate/usecase"

	"github.com/squeakycheese75/open-incentives/internal/httpserver"
	"github.com/squeakycheese75/open-incentives/internal/httputil"
	"github.com/squeakycheese75/open-incentives/internal/store"
)

func run(cfg *configs.APIConfig) error {
	rootCtx := context.Background()

	runtimeAdapter := adapters.New(engine.New())
	actionApplier := services.NewActionApplier()

	db, err := sqlitedb.NewDB(rootCtx, sqlitedb.Config{
		Path:      cfg.DatabasePath,
		Bootstrap: cfg.Bootstrap,
	})
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer func() {
		_ = db.Close()
	}()

	store := store.New(db)

	passwordSvc := services.NewBcryptPasswordService()
	tokenSvc := services.NewJWTTokenService(cfg.ServerJWTSecret)
	publicIDGenerator := services.NanoIDGenerator{}
	cryptoSvc := services.NewCryptoService()

	authUsecaseFactory := usecase_auth.NewUserUsecaseFactory(store.Users(), store.Orgs(), tokenSvc, passwordSvc)
	adminUsecaseFactory := usecase_admin.NewAdminUsecaseFactory(store.Projects(), store.Campaigns(), store.APIKeys(), publicIDGenerator, runtimeAdapter, cryptoSvc, passwordSvc)
	evalUsecaseFactory := usecase_eval.NewAdminUsecaseFactory(store.Campaigns(), runtimeAdapter, actionApplier)

	adminHandler := admin.NewHandler(adminUsecaseFactory)
	authHandler := auth.NewHandler(authUsecaseFactory)

	evalHandler := evaluate.NewHandler(evalUsecaseFactory)

	// authCache := cache.NewAuthContextCache(5 * time.Minute)

	corsCfg := httputil.CORSConfig{AllowedOrigins: parseAllowedOrigins(cfg.CORSAllowedOrigins)}

	mux := httpserver.New(adminHandler, authHandler, evalHandler, tokenSvc, store.Orgs(), store.APIKeys(), passwordSvc, corsCfg)

	srv := &http.Server{
		Addr:    ":" + fmt.Sprint(cfg.ServerPort),
		Handler: mux,
	}

	go func() {
		log.Printf("Starting HTTP server on port %d", cfg.ServerPort)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 2 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Warn("Received request to shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server shutdown failed", "error", err)
	}

	<-ctx.Done()
	slog.Warn("Server exiting ...")

	return nil
}

func parseAllowedOrigins(raw string) []string {
	var origins []string
	for _, origin := range strings.Split(raw, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			origins = append(origins, origin)
		}
	}
	return origins
}

func main() {
	cfg, err := configs.LoadConfig[configs.APIConfig]("")
	if err != nil {
		panic(fmt.Errorf("failed to load config: %v", err.Error()))
	}

	if err := run(cfg); err != nil {
		panic(err)
	}
}
