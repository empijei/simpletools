package main

import (
	"flag"
	"os"
	"os/signal"
)

var rules []listener

var nolog = flag.Bool("silent", false, "do not log the traffic")

//TODO print example
var conf = flag.String("conf", "", "a conf file")

func main() {
	flag.Parse()
	//TODO parse CLI parameters and create rules

	for _, l := range rules {
		l := l
		go l.Listen()
	}

	//Wait for Ctrl+C
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
