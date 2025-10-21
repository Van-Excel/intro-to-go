package main

import (
	"fmt"
)

func main() {

	var fName string = "Vanexcel"
	lName := "Acheampong"
	fmt.Println("first name is:", fName)
	fmt.Printf("The fullname is %s %s\n", fName, lName)
	fmt.Println("The fullname concatenated is:", fName+" "+lName)

	// variables are mappings of human readable names to memory addresses
	// you are essentially telling the compiler to allocate some space
	// and map the address of that memory block to a human readable name
}
