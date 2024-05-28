package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ametow/getblock/pkg"
	"github.com/ametow/getblock/pkg/utils"
)

var configFile = flag.String("c", "etc/config.yaml", "config file")

func main() {
	flag.Parse()

	conf, err := utils.ReadConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	apiKey := os.Getenv("GETBLOCK_KEY")
	if len(apiKey) == 0 {
		log.Fatal("apiKey is empty")
	}
	conf.ApiKey = apiKey

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	svc := pkg.NewService(conf)
	res, err := svc.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Address: %s\nWei: %s\nEth: %.8f\nDollars: %s\n", res.Address, res.Wei.String(), res.Eth, res.Dollars)
}
