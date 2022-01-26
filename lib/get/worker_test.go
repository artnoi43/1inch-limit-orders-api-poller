package get

import (
	"testing"

	"github.com/artnoi43/1inch-limit-orders-api-poller/enums"
	"github.com/artnoi43/1inch-limit-orders-api-poller/lib/utils"
)

func prepareFactory() factoryFunc {
	chans := utils.MakeChans(10)
	return NewFactory(
		enums.ChainEthereum,
		69,
		chans.Order,
		chans.Time,
		chans.Error,
	)
}

func TestFactory(t *testing.T) {
	c := 4
	var poller Poller
	factory := prepareFactory()
	for i := 0; i < c; i++ {
		poller = factory()
	}
	if poller.ID() != c {
		t.Fatal("workerID not matched")
	}
}
