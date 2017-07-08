package main

import (
	"flag"
	"os"
	"os/signal"
)

var rules []listener

var nolog = flag.Bool("silent", false, "do not log the traffic")

func main() {
	flag.Parse()
	//TODO parse CLI parameters and create rules

	//TODO allow TLS -> Plain and Plain -> TLS

	for _, l := range rules {
		l := l
		go l.Listen()
	}

	//Wait for Ctrl+C
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
