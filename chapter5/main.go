package main

import (
	"encoding/json"
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

	// slices
	//creating a slice
	// var car []string

	//using make function to create a slice
	v := make([]string, 5, 8)
	fmt.Println("length of slice v:", len(v))
	fmt.Println("capacity of arr:", cap(v))
	fmt.Println("slice:", v)
	v[0] = "Local"
	fmt.Println("slice:", v)
	//  []string{"van", "ship", "bike"}

	arr := []int{1, 2, 4, 6, 7, 8, 9} // declare and initialize an array
	s := arr[:4]                      // slice the first 4 items in the array
	fmt.Println("new slice:", s)

	//append

	arr2 := [8]int{1, 2, 3, 4, 5}
	slice1 := arr2[0:2]
	// fmt.Println("first viz of slice:", slice1)
	fmt.Printf("flen: %d fcapacity:%d fslice: %v\n", len(slice1), cap(slice1), slice1)
	slice1 = append(slice1, 6, 7)
	fmt.Printf("len: %d capacity:%d slice: %v\n", len(slice1), cap(slice1), slice1)
	fmt.Println("array:", arr2)

	//copy
	slice3 := []int{1, 3, 5}
	slice4 := make([]int, 2)
	copy(slice3, slice4)
	fmt.Println(slice3, slice4)

	// map
	// unordered collection of key value pairs
	//also called a hash table or dictionary or associative arrays
	// noticed in the syntax, the keys are the items in the brackets
	// followed by a value type

	hashtable := make(map[string]int)
	hashtable["dog"] = 4
	hashtable["cat"] = 10
	fmt.Println((hashtable))
	result, status := hashtable["tiger"]
	fmt.Println("result:", result, "status:", status)
	if !status {
		fmt.Println("Key doesn't exist in table")
	}

	//map with int key and person object value with keys name and age
	workers := make(map[int]map[string]any) // map an integer to a map of str:str
	workers[1] = make(map[string]any)
	workers[1]["name"] = "Adjoa"
	workers[1]["age"] = 23
	fmt.Println(workers)
	//safe way with any is type assertion
	worker, ok := workers[1]["age"].(int) // checks if value is int before extraction
	if ok {
		fmt.Println(worker)

	}

	// initializing outer map and then inner map
	// Create the outer map
	employees := make(map[int]map[string]string)

	// Initialize inner maps one by one
	employees[1] = make(map[string]string)
	employees[1]["name"] = "Kofi"
	employees[1]["age"] = "34"

	employees[2] = make(map[string]string)
	employees[2]["name"] = "Ama"
	employees[2]["age"] = "28"

	//use case - serialization
	var d map[string]any
	jsonData := `{"name":"excel", "age":"23"}`
	json.Unmarshal([]byte(jsonData), &d)
	fmt.Println(d)

	// implement a slice on your own

}
