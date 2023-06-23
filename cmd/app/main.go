package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gnaydenova/notification-system/app"
	"github.com/gnaydenova/notification-system/app/handlers"
	"github.com/gnaydenova/notification-system/app/notifications"
)

func startServer(ctx context.Context, cfg *app.Config) {
	l := log.New(os.Stdout, "producer: ", 0)
	producer := notifications.NewProducer(cfg.Producer, l)

	http.Handle("/notifications", handlers.NewNotificationHandler(producer))

	go func() {
		<-ctx.Done()
		producer.Close()
	}()

	if err := http.ListenAndServe(":8090", nil); err != nil {
		panic(err)
	}
}

func startConsumer(ctx context.Context, cfg *app.Config) {
	r := app.NewRegistryFromConfig(cfg)
	l := log.New(os.Stdout, "distributor: ", 0)
	distributor := notifications.NewDistributor(cfg.Distributor, r, l)

	cLog := log.New(os.Stdout, "consumer: ", 0)
	consumer := notifications.NewConsumer(cfg.Consumer, distributor, cLog)

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

	go startServer(ctx, cfg)
	startConsumer(ctx, cfg)
}
