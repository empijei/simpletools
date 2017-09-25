package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
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

var annotation = regexp.MustCompile("^\\s*\\[[A-Z].*\\]\\s*$")

func scanfile(path string) {
	file, err := os.Open(path)
	_panic(err)
	blacklistcheck := func(line string) bool {
		//TODO make this configurable
		blacklist := []string{"Obsolete", "assembly", "TestMethod"}
		for _, word := range blacklist {
			if strings.Contains(line, word) {
				return true
			}
		}
		return false
	}
	whitelistcheck := func(line string) bool {
		whitelist := []string{"Route"}
		if len(whitelist) == 0 {
			return true
		}
		for _, word := range whitelist {
			if strings.Contains(line, word) {
				return true
			}
		}
		return false
	}
	scanner := bufio.NewScanner(file)
	var isOutput, fileprinted, printed bool
	for scanner.Scan() {
		line := scanner.Text()
		if annotation.MatchString(line) && !blacklistcheck(line) {
			isOutput = true
		}
		if isOutput {
			if !annotation.MatchString(line) {
				isOutput = false
				if printed {
					fmt.Println(line)
					fmt.Println()
				}
				printed = false
				continue
			}
			if whitelistcheck(line) {
				if !fileprinted {
					fmt.Println("File: " + path)
					fileprinted = true
				}
				printed = true
				fmt.Println(line)
			}
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
