package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

var jobs = flag.Int("jobs", runtime.GOMAXPROCS(0), "The amount of jobs to run in parallel")
var counter = flag.Int("count", 1, "The amount of times to re-execute the command. \nThis option is ignored if the input parameters are passed via pipe.")
var command = flag.String("cmd", "echo", "The command to execute")

func main() {
	flag.Parse()
	var (
		wg sync.WaitGroup
		ch = make(chan string, *jobs)
		p  = isPipe()
	)
	wg.Add(*jobs)
	for i := 0; i < *jobs; i++ {
		go func() {
			for tok := range ch {
				strcmd := *command
				if p {
					strcmd = strings.Replace(strcmd, "{}", tok, -1)
				}
				c := exec.Command("bash", "-c", strcmd)
				out, err := c.CombinedOutput()
				if err != nil {
					log.Println(err)
					os.Exit(1)
				}
				fmt.Println(string(out))
			}
			wg.Done()
		}()
	}
	go func() {
		if p {
			scan := bufio.NewScanner(os.Stdin)
			for scan.Scan() {
				ch <- scan.Text()
			}
		} else {
			for i := 0; i < *counter; i++ {
				ch <- ""
			}
		}
		close(ch)
	}()
	wg.Wait()
}

func isPipe() bool {
	stdininfo, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while connecting to stdin: %s\n", err.Error())
		return false
	}
	if stdininfo.Mode()&os.ModeCharDevice == 0 {
		//The input is a pipe, so I assume it is what I'm going to use as a token source
		return true
	}
	return false
}
