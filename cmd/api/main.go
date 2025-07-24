package api

import (
	"context"
	"fmt"
	"log"

	"github.com/Neroframe/sub_crudl/pkg/logger"
	"honnef.co/go/tools/config"
)

func main() {
	cfg, err := config.Load("config/dev.yaml")
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.New(logger.Config(cfg.Log))
	logger.Info("config loaded", "version", cfg.Version)
	fmt.Printf("%#v\n", logger.Logger)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Mongo.ConnectTimeout)
	defer cancel()

	// app, err := bootstrap.New(ctx, cfg, logger)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// logger.Info("Starting HTTP server on ", ":", app.Server.Addr)
	// if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 	log.Fatalf("server error: %v", err)
	// }

}
