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

	"github.com/joho/godotenv"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/harshitrajsinha/obj-rest/config"
	v1 "github.com/harshitrajsinha/obj-rest/internal/api/v1"
	"github.com/harshitrajsinha/obj-rest/internal/middleware"
	"github.com/harshitrajsinha/obj-rest/internal/store"
)

func init() {

	_ = godotenv.Load()

	log.SetFlags(log.LstdFlags | log.LUTC | log.Lshortfile)
	log.SetOutput(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    5,
		MaxBackups: 3,
		Compress:   true,
		MaxAge:     10,
	})

}

func main() {

	cfg := config.Load()

	storeClient := store.NewStore(cfg.BaseAPIURL)

	mux := http.NewServeMux()

	// register routes
	v1.RegisterV1Routes(mux, storeClient, cfg.AuthSecretKey)

	muxWithLogs := middleware.LoggingMiddleware(mux)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      muxWithLogs,
		ReadTimeout:  8 * time.Second,
		WriteTimeout: 8 * time.Second,
		IdleTimeout:  8 * time.Second,
	}

	go func() {
		log.Println("starting server at port: ", cfg.Port)
		fmt.Println("starting server at port: ", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error starting server, %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	log.Println("attempting to shutdown server gracefully")

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxWithTimeout); err != nil {
		log.Fatalf("error shutting down server gracefully, %v", err)
	}
}
