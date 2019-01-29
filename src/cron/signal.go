package cron

import (
	"os"
	"os/signal"
	"syscall"
)

// onSignal registers a handler for a given signal
func onSignal(code func(os.Signal), maxNum int, sig ...os.Signal) {

	// create a channel for the signal
	var sigChan chan (os.Signal)

	if maxNum > 0 {
		sigChan = make(chan os.Signal, maxNum)
	} else {
		sigChan = make(chan os.Signal)
	}

	// notify the channel on the signals
	signal.Notify(sigChan, sig...)

	// and run code as long as there are some handlers left
	go (func() {
		for {
			select {
			case signal := <-sigChan:
				code(signal)
			}
		}
	})()
}

//SignalCron signals a cron instance
func SignalCron(proc *os.Process) (err error) {
	// and send SIGUSR1
	proc.Signal(syscall.SIGUSR1)
	return
}
