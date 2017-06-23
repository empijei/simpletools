package main

import (
	"flag"

	colorful "github.com/lucasb-eyer/go-colorful"
)

var color = flag.String("color", "#FF9933", "The color in #HEX notation")

func main() {
	flag.Parse()
	c, _ := colorful.Hex(*color)
	h, s, l := c.Hsl()
	h = float64(int(180+h) % 360)
	cc := colorful.Hsl(h, s, l)
	println(cc.Hex())
}
