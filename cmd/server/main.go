package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/ag-directory-manager/internal/api"
	"github.com/example/ag-directory-manager/internal/config"
	"github.com/example/ag-directory-manager/internal/database"
)

func main() {
	cfg := config.Load()

	db, err := database.OpenAndMigrate(cfg)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	router, err := api.NewRouter(cfg, db)
	if err != nil {
		log.Fatalf("failed to initialize api router: %v", err)
	}

	srv := &http.Server{
		Addr:              cfg.ServerAddress(),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("%s listening on %s", cfg.AppName, cfg.ServerAddress())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = srv.Shutdown(ctx)
	log.Println("server stopped gracefully")
}
