package main

import (
	"fmt"
)

func average(s []int) int {

	total := int(0)

	for _, number := range s {
		total += number

	}
	fmt.Println(total)
	return total
}

//

func main() {
	fmt.Println("ready!")

	//write a loop using range

	type User struct { //define a user struct
		name string
		age  uint
	}

	//declare variable
	var vanexcel User
	//initialize it
	vanexcel = User{
		name: "Vanexcel",
		age:  23,
	}
	fmt.Println(vanexcel)
	// shorthand for declaring and initializing variable
	adjoa := User{name: "adjoa", age: 12}
	fmt.Println(adjoa)
	listOfUsers := [2]User{vanexcel, adjoa}
	for index, element := range listOfUsers {
		fmt.Println("index:", index, "Username:", element.name, "user's age:", element.age)
	}

	//
	//slices
	a := []byte{1, 3, 4}
	fmt.Println(a)
	fmt.Println("length of underlying array for slice a:", len(a))
	fmt.Println("capacity of underlying array for slice a:", cap(a))

	//declaring and initializing an array of 10 elements that holds objects of type float32
	arrayOfFloats := [10]float64{2.3, 3.4, 5.6, 6.7, 6.5, 9.0, 0.9, 2.1, 34.5, 54.32}

	decimalHolder := arrayOfFloats[2:7]
	fmt.Println("decimalholder:", decimalHolder)
	fmt.Println("length of underlying array for slice decimalHolder:", len(decimalHolder))
	fmt.Println("capacity of underlying array for slice decimalHolder:", cap(decimalHolder))
	decimalHolder = append(decimalHolder, 23.8)
	fmt.Println("decimalholder:", decimalHolder)

	//task
	// implement a slice like object and append method
	// if append() is called if (cap- len) < cap and >= len(new elements to add) add to array
	// else create new array
	// if cap-len == 0 create new array
	// new array (create array with size 2* cap, copy existing elements into it, add new elements)

	//maps
	// unordered collection of key value pairs
	//basically a mapping of values
	// study memory structure

	var nameAge map[string]int
	nameAge = make(map[string]int)
	nameAge["van"] = 10
	nameAge["daph"] = 23
	fmt.Printf("age: %d\n", nameAge["daph"])

	schools := map[string]string{
		"A": "Alpha Beta",
		"B": "Barthew School",
	}
	fmt.Println(schools)
	fmt.Printf("length of map: %d \n", len(schools))

	// when you lookup a value and it doesnt exist it returns the zero value
	// for the key type
	hashvalue, ok := nameAge["daph"]
	if ok {
		fmt.Println("exists:", ok, ",", "value:", hashvalue)
	}

	students := map[string]map[string]string{
		"church1": map[string]string{
			"location": "Adenta",
			"country":  "Ghana",
		},
		"church2": map[string]string{
			"location": "Dansoman",
			"country":  "Ghana",
		},
	}
	fmt.Println(students)

	tt := []int{
		21, 45, 67, 89, 00, 43, 67,
	}

	fmt.Println("total:", average(tt))

}
