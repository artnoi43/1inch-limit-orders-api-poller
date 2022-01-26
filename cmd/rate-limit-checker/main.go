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
	f    flags // CLI flags
)

func init() {
	var err error
	conf, err = config.Load()
	if err != nil {
		panic(err)
	}
	f.parse(conf)
	log.Printf(
		"Starting poller\nChain=%s MaxGuard=%d Interval=%d Limit=%d",
		conf.Chain,
		f.maxGuard,
		conf.Get.Interval,
		conf.Get.Limit,
	)
}

func main() {
	jobChan := utils.MakeChans(f.maxGuard)
	sigChan := make(chan os.Signal, 1)
	guard := make(chan struct{}, f.maxGuard)
	signal.Notify(
		sigChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
		syscall.SIGTERM, // kill -SIGTERM XXXX
	)

	pollerFactory := get.NewFactory(
		conf.Chain,
		conf.Get.Limit,
		jobChan.Order,
		jobChan.Time,
		jobChan.Error,
	)
	loop := func() {
		poller := pollerFactory()
		poller.Poll()
		<-guard
	}

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

	// Needed for average GET time
	var getTimes []time.Duration
	for {
		guard <- struct{}{}
		// This Goroutine launches new worker
		go func() {
			loop()
		}()
		// This Goroutine listens on channels
		go func() {
			for {
				select {
				// Process orders here
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
		}()
	}
}
