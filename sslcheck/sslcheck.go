package main

/*TODO
In case the failure is tls handshake failure, provide reasons
Reduce connection timeout
*/

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/empijei/sslcheck/checks"
)

var host = flag.String("host", "", "The domain to check")
var threads = flag.Int("threads", 10, "The domains to check in parallel")

func main() {
	log.SetOutput(os.Stderr)
	flag.Parse()
	if *host != "" {
		fmt.Println(checks.CanConnect(trimmer(*host)))
		return
	}

	//Run on stdin
	input := make(chan string)
	output := make(chan checks.Result)

	//Scan the input
	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		for scanner.Scan() {
			input <- trimmer(scanner.Text())
		}
		close(input)
	}()

	//Create *thread workers
	wg := sync.WaitGroup{}
	wg.Add(*threads)
	for i := 0; i < *threads; i++ {
		go func() {
			for host := range input {
				output <- checks.CanConnect(host)
			}
			wg.Done()
		}()
	}

	//Print output
	done := make(chan struct{})
	go func() {
		for i := range output {
			fmt.Println(i)
		}
		done <- struct{}{}
	}()

	//Wait for scanners
	wg.Wait()
	close(output)

	//Wait for printer
	<-done
}

//Trimmer takes a line in input and tries to extract the host and port information from it
func trimmer(line string) (hostPort string) {
	host := strings.Trim(line, " \t")
	host = strings.Replace(host, "https://", "", 1)
	host = strings.Replace(host, "http://", "", 1)
	host = strings.Split(host, "/")[0]
	port := "443"
	if columnSplit := strings.Split(host, ":"); len(columnSplit) > 1 {
		host = columnSplit[0]
		port = columnSplit[1]
	}
	return host + ":" + port
}
