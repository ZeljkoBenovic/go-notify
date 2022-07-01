package main

import (
	"fmt"
	"github.com/ZeljkoBenovic/go-notify/config"
	"github.com/ZeljkoBenovic/go-notify/monitor"
	"github.com/ZeljkoBenovic/go-notify/notify"
	"os"
)

//TODO: Install service flag

func main() {
	// get conf
	conf, err := config.NewConfig()
	if err != nil {
		fmt.Println("could not create conf error=", err.Error())
		os.Exit(1)
	}

	conf.Logger.Info("Config successfully initialized.")

	// setup notifier instance
	notifier, err := notify.NewNotifier(conf)
	if err != nil {
		conf.Logger.Error("Could not set up notifier service", "error", err.Error())
	}

	// set and run monitor
	newMon, monErr := monitor.NewMonitor(conf)
	if monErr != nil {
		conf.Logger.Error("Could not set up a new sender instance")
	}

	newMon.SetNotifier(notifier)
	newMon.Run()

}
