package main

import (
	"fmt"
)

// Blocks of code are fundamental to programming —
// they let us model real-world actions, tasks, and logic.
// When we need to repeat a task multiple times, we use loops (for, while, or do-while),
// depending on whether we know the number of iterations in advance or not.
// To make decisions about whether or not to perform a task, we use conditional statements like if.
/*
We can combine if statements and loops to perform tasks conditionally and repeatedly,
depending on specific criteria or the number of iterations
*/

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
	// create variable to store state
	// write your conditional expression
	// create block of code you want to execute when
	// conditional expression evaluates to true
	// write state increment/decrement logic - prevents perpetual loop

	// write loop to check if number is even or odd
	for s := 1; s <= 10; s++ {
		if s%2 == 0 {
			fmt.Printf("%d: even\n", s)

		} else {
			fmt.Printf("%d: odd\n", s)

		}

	}

	//switch statement
	// used when there are many conditions

	condition := 2

	if condition == 1 {
		fmt.Println("one")
	} else if condition == 2 {
		fmt.Println("two")
	} else if condition == 3 {
		fmt.Println("three")
	} else if condition == 4 {
		fmt.Println("four")
	} else if condition == 5 {
		fmt.Println("five")
	}

	switch 25 / 5 {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	case 3:
		fmt.Println("three")
	case 4:
		fmt.Println("four")
	case 5:
		fmt.Println("five")
	}

	/*
			TASK
			Write a program that prints the numbers from 1 to 100, but for multiples of
		    three, print “Fizz” instead of the number, and for the multiples of five, print
		    “Buzz.” For numbers that are multiples of both three and five, print “FizzBuz */

	for fiz := 1; fiz <= 100; fiz++ {
		if fiz%3 == 0 {
			fmt.Println("Fizz")
		} else if fiz%5 == 0 {
			fmt.Println("buzz")
		} else if fiz%5 == 0 && fiz%3 == 0 {
			fmt.Println("fizzBuzz")
		} else {
			fmt.Println(fiz)
		}

	}

}
