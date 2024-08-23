package main

import "unsafe"

func main() {
	runes := []rune("👻abc")
	println(unsafe.Sizeof(runes))

}
