// @title       Subscription API
// @version     1.0.0
// @description CRUDL service for user subscriptions
// @host        localhost:8080

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Neroframe/sub_crudl/config"
	"github.com/Neroframe/sub_crudl/internal/app"
	"github.com/Neroframe/sub_crudl/internal/infra/postgres"
	httpapi "github.com/Neroframe/sub_crudl/internal/interfaces/http"
	"github.com/Neroframe/sub_crudl/pkg/logger"

	_ "github.com/Neroframe/sub_crudl/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load("config/dev.yaml")
	if err != nil {
		log.Fatal(err)
	}

	log := logger.New(logger.Config(cfg.Log))
	log.Info("config loaded", "version", cfg.Version)

	// Connect to Postgres with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := sqlx.ConnectContext(ctx, "postgres", postgres.BuildDSN(cfg.Postgres))
	if err != nil {
		log.Fatal("db connect failed", "err", err)
	}

	// Ping to verify conn
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("db ping failed", "err", err)
	}

	// Apply pool settings
	db.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Postgres.ConnMaxLifetime)
	defer db.Close()

	// Wire layers
	// repo := postgres.NewSubscriptionRepo(db, log)
	repo := postgres.NewSubscriptionRepo(db.DB)
	service := app.NewSubscriptionService(repo, log)
	h := httpapi.NewHandler(service, log)

	// Gin setup
	router := gin.Default()
	httpapi.RegisterRoutes(router, h)
	// Init swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	//  Start server
	go func() {
		log.Info("starting HTTP server", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error", "err", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("shutdown error", "err", err)
	} else {
		log.Info("server stopped cleanly")
	}
}
