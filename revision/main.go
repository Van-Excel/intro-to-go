package main

import (
	"fmt"
)

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

}
