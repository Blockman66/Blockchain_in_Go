package main

import "fmt"

type obj struct {
	name string
	age  int
}

func main() {

	A := func(i int) {}
	B := A
	fmt.Printf("%p\n%v\n", &A, A)
	fmt.Printf("%p\n%v", &B, B)
}
