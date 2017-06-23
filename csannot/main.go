package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	walk(".")
}

func walk(path string) {
	fileinfos, err := ioutil.ReadDir(path)
	_panic(err)
	for _, finfo := range fileinfos {
		if finfo.IsDir() {
			walk(path + string(os.PathSeparator) + finfo.Name())
			continue
		}
		if strings.HasSuffix(finfo.Name(), ".cs") {
			scanfile(path + string(os.PathSeparator) + finfo.Name())
		}
	}
}

func _panic(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

const BEGIN = "[Http"
const END = "public"

func scanfile(path string) {
	file, err := os.Open(path)
	_panic(err)
	var isOutput, fileprinted bool
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		trimmedline := strings.Trim(line, " \t")
		if strings.HasPrefix(trimmedline, BEGIN) {
			if !fileprinted {
				fmt.Println("File: " + path)
				fileprinted = true
			}
			isOutput = true
		}
		if isOutput {
			fmt.Println(line) // Println will add back the final '\n'
		}
		if strings.HasPrefix(trimmedline, END) && isOutput {
			isOutput = false
			fmt.Println()
		}
	}
	_panic(scanner.Err())
}

func catfile(in io.Reader) {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // Println will add back the final '\n'
	}
	_panic(scanner.Err())
}
