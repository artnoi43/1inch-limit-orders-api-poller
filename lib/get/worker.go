package get

import (
	"sync"
	"time"

	"github.com/artnoi43/1inch-limit-orders-api-poller/enums"
	"github.com/artnoi43/1inch-limit-orders-api-poller/lib/types"
)

type Poller interface {
	Poll()
	ID() int
	Chain() enums.Chain
	ErrChan() chan error
}

type factoryFunc func() Poller

type worker struct {
	workerID  int
	chain     enums.Chain
	limit     int
	orderChan chan<- []*types.LimitOrder
	timeChan  chan<- time.Duration
	errorChan chan error
}

func (w *worker) ID() int             { return w.workerID }
func (w *worker) Chain() enums.Chain  { return w.chain }
func (w *worker) ErrChan() chan error { return w.errorChan }

// newPoller is wrapped by NewFactory
func newPoller(
	workerID int,
	chain enums.Chain,
	limit int,
	orderChan chan<- []*types.LimitOrder,
	timeChan chan<- time.Duration,
	errorChan chan error,
) Poller {
	return &worker{
		workerID:  workerID,
		chain:     chain,
		limit:     limit,
		orderChan: orderChan,
		timeChan:  timeChan,
		errorChan: errorChan,
	}
}

// NewFactory wraps NewWorker with a closure
func NewFactory(
	chain enums.Chain,
	limit int,
	orderChan chan<- []*types.LimitOrder,
	timeChan chan<- time.Duration,
	errorChan chan error,
) factoryFunc {
	workerID := 0
	var mut sync.RWMutex
	return func() Poller {
		mut.Lock()
		workerID++
		mut.Unlock()
		return newPoller(
			workerID,
			chain,
			limit,
			orderChan,
			timeChan,
			errorChan,
		)
	}
}
