package main

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
)

func init() {
	//This is the creation of 2 sample rules
	rules = append(rules, listener{
		Localport:   ":9443",
		Remoteport:  ":8000",
		Remoteip:    "127.0.0.1",
		Secure:      false,
		ProtoSwitch: false,
	})
	//This is the creation of a sample rule
	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Println(err)
		return
	}

	rules = append(rules, listener{
		Localport:   ":9443",
		Remoteport:  ":8000",
		Remoteip:    "127.0.0.1",
		Secure:      true,
		ProtoSwitch: false,
		certconf: &tls.Config{
			Certificates:       []tls.Certificate{cer},
			InsecureSkipVerify: true,
		},
	})
}

func loadConf(r io.Reader) {
	dec := json.NewDecoder(r)
	var l listener
	var err error
	for ; err != nil; err = dec.Decode(l) {
		if l.Secure {
			cer, cerr := tls.LoadX509KeyPair(l.CrtName, l.KeyName)
			if cerr != nil {
				err = cerr
				break
			}
			certconf := &tls.Config{
				Certificates: []tls.Certificate{cer},
			}
			l.certconf = certconf
		}
		l.certconf.InsecureSkipVerify = l.Secure
		rules = append(rules, l)
	}
	if err != nil && io.EOF != err {
		panic(err)
	}
}
