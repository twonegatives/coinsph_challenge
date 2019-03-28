package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"database/sql"

	"github.com/go-kit/kit/log"
	_ "github.com/lib/pq"
	"github.com/twonegatives/coinsph_challenge/pkg/banking"
	"github.com/twonegatives/coinsph_challenge/pkg/config"
	"github.com/twonegatives/coinsph_challenge/pkg/pgstorage"
)

func main() {
	cfg := config.NewConfig()
	logger := initLogger()

	dbString := cfg.GetString("DB")
	db, err := sql.Open("postgres", dbString)
	defer func() {
		if err = db.Close(); err != nil {
			logger.Log("func", "main", "err", fmt.Sprintf("can't close DB connection to %s", dbString), err)
		}
	}()

	if err != nil {
		logger.Log("func", "main", "err", fmt.Sprintf("can't open DB connection to %s", dbString), err)
		os.Exit(1)
	}

	pgStorage := pgstorage.NewPgStorage(db)

	bankingService := banking.NewService(pgStorage)
	bankingHandler := banking.MakeHandler(bankingService, logger)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", bankingHandler))

	srv := &http.Server{
		Addr:    cfg.GetString("LISTEN"),
		Handler: mux,
	}

	errs := make(chan error)
	go func() {
		logger.Log("func", "srv.ListenAndServe", "msg", "Server is starting", "host", srv.Addr)
		errs <- srv.ListenAndServe()
	}()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case <-stop:
		logger.Log("func", "main", "msg", "interrupt signal received")
	case err := <-errs:
		logger.Log("func", "srv.ListenAndServe", "err", err)
	}

	ctxSD, cancel := context.WithTimeout(context.Background(), cfg.GetDuration("SHUTDOWN_TIMEOUT"))
	defer cancel()
	err = srv.Shutdown(ctxSD)
	logger.Log("msg", "Server was gracefully stopped", "err", err)
	os.Exit(1)
}

func initLogger() log.Logger {
	kitLogger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	kitLogger = log.With(kitLogger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	return kitLogger
}
