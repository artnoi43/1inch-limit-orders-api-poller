package utils

import (
	"time"

	"github.com/artnoi43/1inch-limit-orders-api-poller/lib/types"
)

type WorkerChan struct {
	Order chan []*types.LimitOrder
	Time  chan time.Duration
	Error chan error
}

func MakeChans(i int) *WorkerChan {
	orderChan := make(chan []*types.LimitOrder, i)
	timeChan := make(chan time.Duration)
	errorChan := make(chan error)
	return &WorkerChan{
		Order: orderChan,
		Time:  timeChan,
		Error: errorChan,
	}
}
