package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pkg "github.com/ametow/getblock/internal"
	"github.com/ametow/getblock/internal/utils"
	"github.com/ametow/getblock/internal/web"
)

var configFile = flag.String("c", "etc/config.yaml", "config file")

func run() error {
	flag.Parse()

	conf, err := utils.ReadConfig(*configFile)
	if err != nil {
		return err
	}

	apiKey := os.Getenv("GETBLOCK_KEY")
	if len(apiKey) == 0 {
		log.Fatal("apiKey is empty")
	}
	conf.ApiKey = apiKey

	ctx, svcCancel := context.WithCancel(context.Background())
	defer svcCancel()

	service := pkg.NewService(conf)
	hc := web.NewHandlerContext(service, ctx)

	sm := http.NewServeMux()
	sm.HandleFunc("/api/get", hc.Get)

	server := &http.Server{
		Addr:    conf.ListenAddress,
		Handler: sm,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Starting server on:", conf.ListenAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", conf.ListenAddress, err)
		}
	}()

	<-quit
	svcCancel() // force cancel all running workers

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
	return err
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
