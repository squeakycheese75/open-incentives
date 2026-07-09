package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	engine "github.com/squeakycheese75/open-incentives-engine"
	"github.com/squeakycheese75/open-incentives/configs"
	"github.com/squeakycheese75/open-incentives/pkg/services"

	"github.com/squeakycheese75/open-incentives/internal/admin"
	"github.com/squeakycheese75/open-incentives/internal/admin/auth"
	usecase_auth "github.com/squeakycheese75/open-incentives/internal/admin/auth/usecase"
	usecase_admin "github.com/squeakycheese75/open-incentives/internal/admin/usecase"
	"github.com/squeakycheese75/open-incentives/internal/db/sqlitedb"
	"github.com/squeakycheese75/open-incentives/internal/eval"
	"github.com/squeakycheese75/open-incentives/internal/httpserver"
	"github.com/squeakycheese75/open-incentives/internal/store"
)

func run(cfg *configs.APIConfig) error {
	rootCtx := context.Background()

	engine := engine.New()

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

	authContainer := usecase_auth.NewUserUsecaseContainer(store.Users(), store.Orgs(), tokenSvc, passwordSvc)
	adminContainer := usecase_admin.NewAdminUsecaseContainer(store.Projects(), store.Campaigns())

	adminHandler := admin.NewHandler(adminContainer)
	authHandler := auth.NewHandler(authContainer)
	evalHandler := eval.NewHandler(engine)

	// authCache := cache.NewAuthContextCache(5 * time.Minute)

	mux := httpserver.New(adminHandler, authHandler, evalHandler, tokenSvc, store.Orgs())

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

func main() {
	cfg, err := configs.LoadConfig[configs.APIConfig]("")
	if err != nil {
		panic(fmt.Errorf("failed to load config: %v", err.Error()))
	}

	if err := run(cfg); err != nil {
		panic(err)
	}
}
