package main

import (
	"crypto/tls"
	"log"
)

func init() {
	//This is the creation of 2 sample rules
	rules = append(rules, listener{
		localport:   ":9443",
		remoteport:  ":8000",
		remoteip:    "127.0.0.1",
		secure:      false,
		protoSwitch: false,
	})
	//This is the creation of a sample rule
	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Println(err)
		return
	}

	rules = append(rules, listener{
		localport:   ":9443",
		remoteport:  ":8000",
		remoteip:    "127.0.0.1",
		secure:      true,
		protoSwitch: false,
		certconf: &tls.Config{
			Certificates:       []tls.Certificate{cer},
			InsecureSkipVerify: true,
		},
	})
}
