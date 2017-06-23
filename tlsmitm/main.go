package main

import (
	"os"
	"os/signal"
)

var rules []listener

func main() {

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
