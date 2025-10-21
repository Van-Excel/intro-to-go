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

	//declare variable
	var school string
	school = "alpha beta"
	fmt.Println("name of school:", school)
	school = "martin de porres"
	fmt.Println("name of new school is:", school)
	fmt.Println("is your school alpha beta?", school == "alpha beta")

	// variables are mappings of human readable names to memory addresses
	// you are essentially telling the compiler to allocate some space
	// and map the address of that memory block to a human readable name
}
