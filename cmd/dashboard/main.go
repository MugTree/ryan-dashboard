package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/MugTree/ryan_dashboard/dashboard"

	_ "embed"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	_ "modernc.org/sqlite"

	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {

	if err := run(); err != nil {
		slog.Error("app error", "error", err)
		os.Exit(1)
	}
}

func run() error {

	/* validate the environment
	---------------------------------------------------
	1. Locally-  Running locally we pull the env vars from a .env file
	2. Local docker - Running docker locally we can pull the env values from cmd line
	3. Cloud run -  When run on cloud run we pull the vals in from the environment no .env involved
	---------------------------------------------------
	*/

	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			return fmt.Errorf(".env file present but can't load it %v", err)
		}
	}

	env := dashboard.EnvVars{
		IsProd:      dashboard.MustEnvGetBool("IS_PRODUCTION"),
		LogLocation: dashboard.MustEnv("APP_LOG"),
	}

	host := dashboard.MustEnv("HOST")
	dbPath := dashboard.MustEnv("DB")

	rotator := &lumberjack.Logger{
		Filename:   env.LogLocation,
		MaxSize:    50, // MB
		MaxBackups: 5,
		MaxAge:     30, // days
		Compress:   true,
	}

	handler := slog.NewTextHandler(rotator, &slog.HandlerOptions{Level: slog.LevelInfo})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	db, err := sqlx.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	defer db.Close()

	_, err = db.Exec("PRAGMA foreign_keys = ON; PRAGMA journal_mode = WAL;")
	if err != nil {
		return fmt.Errorf("failed to enable foreign keys or wal mode: %v ", err)
	}

	s := dashboard.NewServer(db, host, &env)

	// Use an errgroup to wait for separate goroutines which can error
	eg, ctx := errgroup.WithContext(ctx)

	// Start the server within the errgroup.
	eg.Go(func() error {
		return s.Start()
	})

	// Wait for the context to be done, which happens when a signal is caught
	<-ctx.Done()
	slog.Info("Stopping the app")

	// Stop the server gracefully
	eg.Go(func() error {
		return s.Stop()
	})

	// Wait for the server to stop
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("wait group error: %v", err)
	}

	slog.Info("Stopped the app")

	return nil
}
