package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

const encr = "encrypted"
const unen = "unencrypted"

type listener struct {
	Localport, Remoteport, Remoteip string
	//Tells the listener to setup a tls socket
	Secure bool `"json:,omitempty"`
	//Makes the tool change protocol, this can be used to downgrade a secure
	//connection or to upgrade an insecure one
	ProtoSwitch bool `"json:,omitempty"`

	// Used only during deserialization, name of private key of listener
	KeyName string `"json:,omitempty"`
	// Used only during deserialization, name of public key of listener
	CrtName string `"json:,omitempty"`

	//If false does not validate certificates upon connection to remote hosts
	Verify bool `"json:,omitempty"`

	// Unexported field, will be created on load
	certconf *tls.Config
}

func (l *listener) String() string {
	var localtext, remotetext string
	if l.Secure {
		localtext = encr
		if l.ProtoSwitch {
			remotetext = unen
		} else {
			remotetext = encr
		}
	} else {
		localtext = unen
		if l.ProtoSwitch {
			remotetext = encr
		} else {
			remotetext = unen
		}
	}
	return localtext + " " + l.Localport + " -> " + l.Remoteip + l.Remoteport + " " + remotetext
}

func (l *listener) Listen() {
	//TODO check if the listener was properly constructed

	var ll net.Listener
	var err error
	if l.Secure {
		ll, err = tls.Listen("tcp", l.Localport, l.certconf)
	} else {
		ll, err = net.Listen("tcp", l.Localport)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Running rule " + l.String())
	for {
		conn, err := ll.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(lc net.Conn) {
			filename := fmt.Sprintf("%d", time.Now().UnixNano())
			log.Printf("Got incoming connection from ip %s on port %s", conn.RemoteAddr().String(), l.Localport)
			defer func() { _ = lc.Close() }()
			//Dump client traffic to stdout
			var lt io.Reader
			if *nolog {
				lt = lc
			} else {
				lt = io.TeeReader(lc, os.Stdout)

				//Dump client traffic to file
				outclientfile, filerr := os.OpenFile(
					fmt.Sprintf("%s_client.log", filename),
					os.O_APPEND|os.O_CREATE|os.O_WRONLY,
					0666)
				if filerr != nil {
					log.Printf("Could not open logfile: %s\n, proceeding without it", filerr.Error())
				} else {
					lt = io.TeeReader(lt, outclientfile)
				}
				_, _ = outclientfile.Write([]byte("Remote: " + conn.RemoteAddr().String() + " " + l.String() + "\n"))
			}

			var rc net.Conn
			log.Printf("Contacting remote server @%s", l.Remoteip+l.Remoteport)
			//This means (l.secure && !l.protoSwitch) || (!l.secure && l.protoSwitch)
			if l.Secure != l.ProtoSwitch {
				rc, err = tls.Dial("tcp", l.Remoteip+l.Remoteport,
					&tls.Config{
						//TODO make this configurable
						InsecureSkipVerify: true,
					})
			} else {
				rc, err = net.Dial("tcp", l.Remoteip+l.Remoteport)
			}
			if err != nil {
				//TODO handle this
				log.Println(err)
				return
			}
			log.Println("Connected to remote server")

			defer func() { _ = rc.Close() }()
			var rt io.Reader
			if *nolog {
				rt = rc
			} else {
				//Dump server traffic to stdout
				rt := io.TeeReader(rc, os.Stdout)
				//Dump server traffic to file
				outserverfile, filerr := os.OpenFile(
					fmt.Sprintf("%s_server.log", filename),
					os.O_APPEND|os.O_CREATE|os.O_WRONLY,
					0666)
				if filerr != nil {
					log.Printf("Could not open logfile: %s\n, proceeding without it", filerr.Error())
				} else {
					rt = io.TeeReader(rt, outserverfile)
				}
			}

			go func() {
				_, err := io.Copy(rc, lt)
				if err != nil {
					log.Printf("Error while sending to remote server: %s", err.Error())
				}
				_ = conn.Close()
				_ = rc.Close()
			}()
			_, err = io.Copy(lc, rt)
			_ = conn.Close()
			_ = rc.Close()
			if err != nil {
				log.Printf("Error while receiving from remote server: %s", err.Error())
			}
			log.Printf("Connection with %s closed", conn.RemoteAddr().String())
		}(conn)
	}

}
