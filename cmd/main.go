package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/DillonEnge/jolt/internal/api"
	server "github.com/DillonEnge/jolt/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	wait := make(chan os.Signal, 1)
	signal.Notify(wait, os.Interrupt)

	err := run(context.Background(), os.Args[1:], wait)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		log.Panic(err)
	}
}

func run(ctx context.Context, _ []string, wait chan os.Signal) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	config := api.NewConfig()

	dbPool, err := pgxpool.New(ctx, config.DBUrl)
	if err != nil {
		return err
	}
	defer dbPool.Close()

	// Run the service logic and wait for an interrupt.
	stopService, err := server.Service(ctx, dbPool, config)
	defer stopService()
	if err != nil {
		return err
	}
	<-wait

	slog.Info("Service has gracefully terminated.")
	return nil
}
