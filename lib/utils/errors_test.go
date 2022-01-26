package utils

import (
	"os"
	"syscall"
	"testing"

	"github.com/artnoi43/1inch-limit-orders-api-poller/enums"
)

func TestHandleError(t *testing.T) {
	c := make(chan os.Signal)
	go HandleError(c, enums.ErrRateLimitExceeded)
	if sig := <-c; sig != syscall.SIGABRT {
		t.Fatal("Unexpected behaviour")
	}
}
