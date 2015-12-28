package main

import (
	"os/signal"
	"syscall"

	lazlo "github.com/djosephsen/lazlo/lib"
)

func main() {

	lazlo.Logger.Debug(`creating broker`)
	//make a broker
	broker, err := lazlo.NewBroker()
	if err != nil {
		lazlo.Logger.Error(err)
		return
	}
	defer broker.Brain.Close()

	lazlo.Logger.Debug(`starting modules`)
	// register the modules
	if err := initModules(broker); err != nil {
		lazlo.Logger.Error(err)
		return
	}
	//start the Modules
	broker.StartModules()

	lazlo.Logger.Debug(`starting broker`)
	//start the broker
	go broker.Start()
	// Loop
	signal.Notify(broker.SigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	stop := false
	for !stop {
		select {
		case sig := <-broker.SigChan:
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				stop = true
			}
		}
	}
	// Stop listening for new signals
	signal.Stop(broker.SigChan)
	broker.Stop()

	//wait for the write thread to stop (so the shutdown hooks have a chance to run)
	<-broker.SyncChan
}
