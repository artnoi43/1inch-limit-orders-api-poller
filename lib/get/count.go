package get

import (
	"log"
	"sync"

	"github.com/artnoi43/1inch-limit-orders-api-poller/enums"
	"github.com/pkg/errors"
)

// Dependency injection
type getCountFunc func(enums.Chain) (int, error)

var count int
var countSet bool

// PollCount uses given getFunc (usually GetCount)
// to poll number of orders,
// and set it to the package's global var count and countSet
func PollCount(
	chain enums.Chain,
	f getCountFunc,
	errChan chan<- error,
) {
	var mut sync.RWMutex
	if newCount, err := f(chain); err != nil {
		errChan <- errors.Wrap(
			err, "failed to get count",
		)
	} else if count != newCount {
		log.Println("Count set to", newCount)
		mut.Lock()
		count = newCount
		countSet = true
		mut.Unlock()
	}
}
