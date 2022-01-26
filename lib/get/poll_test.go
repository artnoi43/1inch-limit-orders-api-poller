package get

import (
	"sync"
	"testing"
	"time"

	"github.com/artnoi43/1inch-limit-orders-api-poller/enums"
)

// TODO: better test cases
func TestPoll(t *testing.T) {
	fail := func() {
		t.Fatal("Got error from channel")
	}
	factory := prepareFactory()
	poller := factory()
	var wg sync.WaitGroup
	// This Goroutine periodically sets get.Count
	go func() {
		for {
			PollCount(
				enums.ChainEthereum,
				GetCount,
				poller.ErrChan(),
			)
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		poller.Poll()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-poller.ErrChan():
			fail()
		case <-time.After(10 * time.Second):
			return
		}
	}()
	wg.Wait()
}
