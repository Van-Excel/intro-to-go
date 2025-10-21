package main

import (
	"fmt"
)

func main() {

	// TYPES
	// numbers - floats and integers (signed and unsigned)
	//strings
	// boolean operators- &&-and || or  ! not

	fmt.Println("1 + 1 =", 1.0+1.0)
	fmt.Printf("1.0 + 2.0 = %f \n", 1.0+2.0)

	//Strings
	//fmt.Fprintln(os.Stdout, []any{len("Hello World,")}...)
	a, err := fmt.Println(len("hello world."))
	if err != nil {
		fmt.Println("Error detected")
	}
	fmt.Println("value of a:", a)
	c := len("hello world")
	fmt.Println("The length of hello world is:", c)
}
