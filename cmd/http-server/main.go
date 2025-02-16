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

	"github.com/vit6556/avito-internship-assignment/internal/app"
	"github.com/vit6556/avito-internship-assignment/internal/config"
)

func main() {
	dbPool := app.InitDatabase()

	cfg := config.LoadServerConfig()
	echo := app.InitServer(cfg, dbPool)

	go func() {
		if err := echo.Start(fmt.Sprintf(":%d", cfg.HTTPServer.Port)); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %s", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := echo.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	log.Println("Closing database connection...")
	dbPool.Close()

	log.Println("Server stopped gracefully")
}
