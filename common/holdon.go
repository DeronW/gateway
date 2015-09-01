package common

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func HoldOn() {
	fmt.Println("Hold On...")
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Exit Signal", <-chSig)
}
