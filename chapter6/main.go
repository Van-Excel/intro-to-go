package main

import (
	"fmt"
)

//functions
// a body of statements that must be executed together
// used to model an action or task
// a functions signature is the parameters and return type

// function to calculate the average
func average(container []float64) float64 {
	total := 0.0
	for _, element := range container {
		total += element
	}

	return total

}

func splitNumbers(item ...int) int { // argument gathers the values passed into a slice
	total := 0
	for _, element := range item {
		total += element
	}
	return total
}

func bookPtr(element *[]int) *[]int {
	for index := range *element {
		(*element)[index] += 10

	}
	return element
}

func main() {
	fmt.Println("ready!")
	example := []float64{1.0, 3.4, 5.6, 6.7}
	answer := average(example)
	fmt.Println("average:", answer)
	fmt.Println("variadic function value:", splitNumbers(2, 3, 5, 6))

	// pointers

	bookIndex := []int{1, 3, 5, 7}
	bookPointer := &bookIndex
	fmt.Println(bookPtr(bookPointer))
	fmt.Println("current value of bookIndex:", bookIndex)

}
