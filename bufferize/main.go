package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	_, err = os.Stdout.Write(buf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
