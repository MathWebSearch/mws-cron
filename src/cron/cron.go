package cron

import (
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"

	cron "gopkg.in/robfig/cron.v2"
)

// Task represents a task run by cron
type Task func(Reason, func())

// Reason represents the reason a task was started
type Reason int

const (
	// Scheduled implies the CronJob has been scheduled
	Scheduled Reason = iota
	// Initial run of job during
	Initial
	// Signalled manually
	Signalled
)

// RunCron runs cron with a given interval
func RunCron(cronLine string, code Task, allowSignals bool, runInitial bool) bool {

	// mutex to lock the job
	mutex := &sync.Mutex{}

	// if we allow signals, listen to SIGUSR1 for manually starting the task
	if allowSignals {
		onSignal(func(_ os.Signal) {
			fmt.Printf("Cron: Received SIGUSR1. \n")
			quenueCronTask(mutex, code, Signalled)
		}, 0, syscall.SIGUSR1)
	}

	// run the initial go task
	if runInitial {
		fmt.Printf("Cron: Running initial job. \n")
		go quenueCronTask(mutex, code, Initial)
	}

	// create a cron instance
	c := cron.New()

	var entryID cron.EntryID
	var err error

	entryID, err = c.AddFunc(cronLine, func() {
		fmt.Printf("Cron: Running regular job, next scheduled at %s\n", c.Entry(entryID).Next.Format(time.RFC3339))
		quenueCronTask(mutex, code, Scheduled)
	})

	if err != nil {
		return false
	}

	// start cron
	c.Start()
	defer c.Stop()

	// log and wait for interrupt
	fmt.Printf("Cron: Starting, next scheduled at %s\n", c.Entry(entryID).Next.Format(time.RFC3339))

	end := make(chan int)

	// listen to a single interrupt
	onSignal(func(_ os.Signal) {
		fmt.Printf("Cron: Received INTERRUPT, exiting safely. \n")
		mutex.Lock()
		end <- 0
	}, 1, os.Interrupt)

	// wait forever; then return true
	<-end
	return true
}

func quenueCronTask(mutex *sync.Mutex, code Task, reason Reason) {
	var localTask func()
	localTask = func() {
		// lock the mutex until we are done
		mutex.Lock()
		defer mutex.Unlock()

		// and run the code
		code(reason, localTask)
	}

	go localTask()
}
