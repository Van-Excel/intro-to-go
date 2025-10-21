package main

import "fmt"

//this is a comment

func main() {

	c, err := fmt.Println("this is my first program in go")
	if err != nil {
		fmt.Println("error detected")
	}

	fmt.Printf("The number of bytes returned: %d bytes", c)
}

// allows you to either just declare a variable or
// to declare and initialize
// go has a shorthand syntax for declaring and initializing variables.  f := "Van"
// when a variable is only declared it is auto initialized
// with the zero value of its data type
var d int
var e int = 6
var f = 7
