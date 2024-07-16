package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/crnvl96/distributed-tracing-and-spam/service_b/api/handler"
	"github.com/crnvl96/distributed-tracing-and-spam/service_b/api/middleware"
	"github.com/crnvl96/distributed-tracing-and-spam/service_b/infra/opentelemetry"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelShutdown, err := opentelemetry.Init(ctx)
	if err != nil {
		return
	}

	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/zipcode", middleware.CallTrace(http.HandlerFunc(handler.CalculateTemperature)))

	srv := &http.Server{
		Addr:         ":8081",
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}

	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	select {
	case err = <-srvErr:
		return
	case <-ctx.Done():
		stop()
	}

	err = srv.Shutdown(context.Background())
}
