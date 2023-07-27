package main

import (
	"fmt"
	"slice-test/slice"
)

func main() {
	s := []int{1, 2, 3, 4, 5, 6, 7}

	println(cap(s))
	println(len(s))

	s = slice.DeleteElement(s, 2)

	println(cap(s))
	println(len(s))
	fmt.Printf("%#v\n", s)

}
