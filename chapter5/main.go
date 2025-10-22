package main

import (
	"fmt"
)

//initialize an array using loops
// array is like a box or container [] which will store items of the same type
// These items will be instances of a type which share similar properties
/*
	just realised an array is essentially modelled after a real life box or container which
	stores items of the same type. so instances of the same class which share similar properties

	when we use the size of the array conditionally in programming, you are asking if this box or
	container is full or has space to hold more items.
	size of array is also capacity of the box or container
	the items are also numbered or ordered sequentially

*/

func main() {

	var name [3]int
	name[1] = 23
	name[0] = 12
	name[2] = 24
	fmt.Println(name)

	// program to initialize an array

	var box [4]int
	for i := 0; i < 4; i++ { // check if box is empty
		box[i] = i // put something in box as long as its not full
	}
	fmt.Println("contents of container:", box)
	fmt.Println(box[0]) // give me the contents in the first box or first item in the box

	// range
	// similar to enumerate function in python
	// gives you access to the index and the value at the index

	var counter int
	// var index int
	for index, value := range box { // no need to declare iterator variable and value here
		fmt.Println(index, ":", value)
		counter += value

	}
	fmt.Println("current value of counter:", counter)

	// another way to declare array
	container := [4]int{2, 3, 5, 6}
	fmt.Println(container)

}
