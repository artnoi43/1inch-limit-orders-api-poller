package main

import (
	"flag"

	"github.com/artnoi43/1inch-limit-orders-api-poller/config"
)

type flags struct {
	maxGuard int
}

func (f *flags) parse(conf *config.Config) {
	flag.IntVar(&f.maxGuard, "c", conf.MaxGuard, "MaxGuard")
	flag.Parse()
}
