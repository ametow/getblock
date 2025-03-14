package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	pkg "github.com/ametow/getblock/internal"
	"github.com/ametow/getblock/internal/utils"
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
	fmt.Printf("Number of active goroutines: %d\n", runtime.NumGoroutine())
	// printGoroutines()
}

// func printGoroutines() {
// 	buf := make([]byte, 1<<20) // 1 MB buffer
// 	n := runtime.Stack(buf, true)
// 	fmt.Printf("%s", buf[:n])
// }
