package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/MathWebSearch/mws-cron/src/cron"
	"github.com/MathWebSearch/mws-cron/src/docker"
)

func main() {
	if shouldTrigger {
		mainTrigger()
	} else {
		mainCron()
	}
}

func mainCron() {
	fmt.Printf("Cron: Using pidfile %s\n", pidFile)
	// write the pidfile
	if cron.WritePid(pidFile) != nil {
		panic("Can't write pidfile")
	}

	// clear the pid on exit
	defer cron.ClearPid(pidFile)

	// and run cron
	minute := 1 * time.Minute
	cron.RunCron(cronLine, func(reason cron.Reason, retry func()) {
		docker.UpdateMWS(mwsLabel, retry, &minute)
	}, true, shouldRunInitial)
}

func mainTrigger() {
	// load the pidfile
	var proc *os.Process
	var err error
	if proc, err = cron.ReadPid(pidFile); err != nil {
		panic("Can't load pidfile")
	}

	// and send the signal
	cron.SignalCron(proc)
}

var pidFile string
var mwsLabel string
var cronLine string
var shouldTrigger bool
var shouldRunInitial bool

func init() {
	flag.BoolVar(&shouldTrigger, "trigger", false, "Trigger manually running a cron job in running instance")
	flag.BoolVar(&shouldRunInitial, "initial", true, "Run re-indexing on startup")
	flag.StringVar(&pidFile, "pidfile", "", "Pidfile to use")
	flag.StringVar(&cronLine, "schedule", "@midnight", "Cronline representing time to run job on, use '@never' to disable cronjobs. ")
	flag.StringVar(&mwsLabel, "label", "org.mathweb.mwsd", "Label for MathWebSearch daemon")
	flag.Parse()

	// if pidfile is empty, use a file int he current directory
	if pidFile == "" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(err)
		}
		pidFile = filepath.Join(dir, "mws-cron.pid")
	}
}
