/*******************************************************************************
 * Maxima Pool is a maxima server to deal with requests from moodle-qtype_stack
 *
 * @author     Lars Thoms <lars@thoms.io>
 * @date       2022-11-25
 ******************************************************************************/

package main

import (
	"Moodle_Maxima_Pool/services"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	ctxTimeout         = 3 * time.Second
	terminationTimeout = 3 * time.Second
)

var (
	logger     = NewLogger()
	waitGroup  sync.WaitGroup
	terminator = make(chan struct{})
)

func init() {
	go func() {

		// Listen to interrupt and termination signals
		termSignal := make(chan os.Signal)
		signal.Notify(termSignal, os.Interrupt, syscall.SIGTERM)
		<-termSignal
		close(terminator)

		// Guarantee termination after specified timeout
		waitChannel := make(chan struct{})
		go func() {
			defer close(waitChannel)
			waitGroup.Wait()
		}()

		select {
		case <-waitChannel:
		case <-time.After(terminationTimeout):
		}

		os.Exit(143)
	}()
}

func main() {
	if err := loadConfig(); err != nil {
		logger.Fatal(err)
	}
	logger.SetLevel(viper.GetInt("loglevel"))

	if *createSnapshots {
		err := services.MaximaSnapshotCreate()
		if err != nil {
			logger.Fatal(err)
		}
	} else if _, err := services.MaximaSnapshotGet(""); err != nil {
		logger.Fatal(err)
	} else {
		startHTTPServer()
		waitGroup.Wait()
	}
}
