package cron

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
)

// WritePid writes a pid file
func WritePid(filename string) (err error) {
	// check if we already have a valid pidfile at that location
	if _, err := ReadPid(filename); err == nil {
		return errors.New("Valid pidfile already exists. ")
	}

	// if not, make one
	data := []byte(strconv.FormatInt(int64(os.Getpid()), 10))
	return ioutil.WriteFile(filename, data, 0644)
}

// ClearPid removes a pid file
func ClearPid(filename string) (err error) {
	return os.Remove(filename)
}

// ReadPid reads the pid file
func ReadPid(filename string) (process *os.Process, err error) {
	// read the file
	var b []byte
	if b, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	// read the pid
	var pid int64
	if pid, err = strconv.ParseInt(string(b), 10, 32); err != nil {
		return
	}

	// and turn it into a process
	return os.FindProcess(int(pid))
}
