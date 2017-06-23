package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/miekg/dns"
)

var ipv6 = flag.Bool("6", false, "Print also IPv6 IPs")

func main() {
	flag.Parse()
	ifaces, err := net.Interfaces()
	// handle err
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if !ip.IsLoopback() && (ip.To4() != nil || *ipv6) {
				fmt.Printf("%v: %v\n", i.Name, ip)
			}
			// process IP address
		}
	}
	fmt.Println("External: " + external())
}

//dig +short myip.opendns.com @resolver1.opendns.com
//TODO add ipv6 support?
func external() (ext string) {
	ext = "Not found"
	target := "myip.opendns.com"
	server := "resolver1.opendns.com"

	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(target+".", dns.TypeA)
	r, _, err := c.Exchange(&m, server+":53")
	if err != nil {
		//log.Println(err)
		return
	}
	for _, ans := range r.Answer {
		Arecord := ans.(*dns.A)
		ext = Arecord.A.String()
	}
	return
}
