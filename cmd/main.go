package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/artnoi43/1inch-limit-orders-api-poller/config"
	"github.com/artnoi43/1inch-limit-orders-api-poller/lib/get"
	"github.com/artnoi43/1inch-limit-orders-api-poller/lib/utils"
)

var (
	conf *config.Config
)

func init() {
	var err error
	conf, err = config.Load()
	if err != nil {
		panic(err)
	}
	log.Printf(
		"Starting poller\nChain=%s Interval=%d, Limit=%d",
		conf.Chain,
		conf.Get.Interval,
		conf.Get.Limit,
	)
}

func main() {
	jobChan := utils.MakeChans(conf.MaxGuard)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(
		sigChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
		syscall.SIGTERM, // kill -SIGTERM XXXX
	)

	var getTimes []time.Duration
	ticker := time.NewTicker(
		time.Duration(conf.Get.Interval) * time.Second,
	)
	defer ticker.Stop()

	// This Goroutine periodically sets get.count and get.countSet
	go func() {
		for {
			get.PollCount(
				conf.Chain,
				get.GetCount,
				jobChan.Error,
			)
			time.Sleep(
				time.Duration(conf.Get.Interval) * time.Second,
			)
		}
	}()

	factory := get.NewFactory(
		conf.Chain,
		conf.Get.Limit,
		jobChan.Order,
		jobChan.Time,
		jobChan.Error,
	)
	loop := func() {
		poller := factory()
		poller.Poll()
	}

	go loop()
	for {
		select {
		// Launch new worker every interval second
		case <-ticker.C:
			go loop()
		//  Process orders here
		case orders := <-jobChan.Order:
			for _, order := range orders {
				order.ParseTime()
			}
		// Add duration to getTimes
		case getTime := <-jobChan.Time:
			getTimes = append(getTimes, getTime)
		// Handle 429
		case err := <-jobChan.Error:
			utils.HandleError(sigChan, err)
		// Graceful shutdown, here we just print average GET time
		case <-sigChan:
			log.Println("Shutting down poller")
			utils.CalcAvgTime(getTimes)
			os.Exit(0)
		}
	}
}
