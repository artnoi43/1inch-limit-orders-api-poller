package get

import (
	"testing"

	"github.com/artnoi43/1inch-limit-orders-api-poller/enums"
)

func mockCount() func(enums.Chain) (int, error) {
	c := 0
	return func(enums.Chain) (int, error) {
		c++
		return c, nil
	}
}

func TestPollCount(t *testing.T) {
	target := 4
	c := make(chan error)
	f := mockCount()
	for i := 0; i < target; i++ {
		PollCount(enums.ChainEthereum, f, c)
	}
	if !countSet {
		t.Fatal("Count not set")
	}
	if count != target {
		t.Fatal("Count invalid")
	}
}
