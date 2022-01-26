package utils

import (
	"errors"
	"log"
	"os"
	"syscall"

	"github.com/artnoi43/1inch-limit-orders-api-poller/enums"
)

func HandleError(c chan os.Signal, e error) {
	if errors.Is(enums.ErrRateLimitExceeded, e) {
		c <- syscall.SIGABRT
	} else {
		log.Println("Found error:", e.Error())
	}
}
