package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"
)

var jobs = flag.Int("jobs", runtime.GOMAXPROCS(0), "The amount of jobs to run in parallel")
var counter = flag.Int("count", 1, "The amount of times to re-execute the command")
var command = flag.String("cmd", "echo", "The command to execute")

func main() {
	flag.Parse()
	var (
		wg sync.WaitGroup
		ch = make(chan struct{}, *jobs)
	)
	wg.Add(*jobs)
	for i := 0; i < *jobs; i++ {
		go func() {
			for _ = range ch {
				c := exec.Command("bash", "-c", *command)
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
		for i := 0; i < *counter; i++ {
			ch <- struct{}{}
		}
		close(ch)
	}()
	wg.Wait()
}
