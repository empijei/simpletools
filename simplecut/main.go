package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var d = flag.String("d", " ", "The delimiter string")
var f = flag.String("f", "0", "The comma-separated fields to output. 0 represents the whole line")

func main() {
	flag.Parse()
	fieldslist := parsef()
	fieldsplitter(os.Stdout, os.Stdin, *d, fieldslist)
}

func fieldsplitter(in io.Reader, out io.Writer, d string, fieldslist []int) {
	outslice := make([]string, 0, len(fieldslist))
	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		t := scan.Text()
		splitted := strings.Split(t, d)
		for _, f := range fieldslist {
			if f == 0 {
				outslice = append(outslice, t)
				continue
			}
			f--
			if f < len(splitted) {
				outslice = append(outslice, splitted[f])
			} else {
				break
			}
		}
		fmt.Println(strings.Join(outslice, d))
		outslice = outslice[:0] //reset outslice
	}
}

func parsef() []int {
	strs := strings.Split(*f, ",")
	indexes := make([]int, 0, len(strs))
	for _, i := range strs {
		ii, err := strconv.Atoi(i)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		indexes = append(indexes, ii)
	}
	return indexes
}
