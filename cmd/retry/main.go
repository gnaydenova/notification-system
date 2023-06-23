package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gnaydenova/notification-system/app"
	"github.com/gnaydenova/notification-system/app/notifications"
)

func startConsumer(ctx context.Context, cfg *app.Config) {
	r := app.NewRegistryFromConfig(cfg)
	l := log.New(os.Stdout, "retry distributor: ", 0)
	distributor := notifications.NewDistributor(cfg.Distributor, r, l)

	cLog := log.New(os.Stdout, "retry consumer: ", 0)
	consumer := notifications.NewConsumer(cfg.RetryQueue, distributor, cLog)

	consumer.Consume(ctx)
}

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		fmt.Printf("\nreceived signal: %v\n", <-done)
		cancel()
	}()

    var configPath string
    flag.StringVar(&configPath, "config", "./config.yaml", "path to config file")
    flag.Parse()

	cfg, err := app.NewConfig(configPath)
	if err != nil {
		panic(err)
	}

	startConsumer(ctx, cfg)
}
