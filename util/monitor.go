package util

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// Monitor monitoring OS signals.
func Monitor(done chan bool) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
}
