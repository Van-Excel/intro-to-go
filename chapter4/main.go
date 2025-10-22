package main

import (
	"fmt"
)

func main() {
	fmt.Println("chapter 4")
	// control structures

	fmt.Println("1")
	fmt.Println("2")
	fmt.Println("3")
	fmt.Println("4")

	// write a loop that makes writing this code less tedious
	// loops are used to repeat blocks of code or tasks to be performed
	// blocks of code represent real world tasks to be performed

	var iteration int
	for iteration = 1; iteration <= 10; iteration++ {
		fmt.Println("iteration:", iteration)
	}

	for count := 1; count <= 10; count++ {
		fmt.Println("count:", count)

	}

}
