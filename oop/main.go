package main

import (
	"fmt"
)

type Box struct {
	length float64
	width  float64
}

func (b *Box) area() float64 {
	return b.length * b.width

}

// embedded types: is-a relationship
// interface was created after implementing human and employee structs
// noticed patterns and abstracted it above to create humanbeing interface
// lesson is implement concrete types first and abstract as pattern emerges
type HumanBeing interface {
	speak() string
}

type Human struct {
	name        string
	age         int
	nationality string
}

func (h Human) speak() string {
	// fmt.Println("my name is", h.name)
	return "my name is" + " " + h.name + " " + "and I am a" + " " + h.nationality

}

type Employee struct {
	Human
	skills  []string
	company string
}

// variadic allows caller to pass either slice or individual item
// []Humanbeing limits caller to pass only slice
func allEmployees(e ...HumanBeing) string {
	for _, employee := range e {
		fmt.Println(employee.speak())
	}
	return "all employees spoke"
}

func main() {

	fmt.Println("hello world of OOP")

	b := Box{length: 23.9, width: 45.4}
	fmt.Println("area of box:", b.area())
	h := Human{name: "Vanexel", age: 23, nationality: "Ghanaian"}
	e := Employee{Human: h, skills: []string{"python", "go"}, company: "g4m"}
	fmt.Println("employee speaking:", e.speak())

	//creating humans for employees slice
	h1 := Human{name: "Kojo", age: 23, nationality: "Ghanaian"}
	h2 := Human{name: "Chibeke", age: 23, nationality: "Nigerian"}
	h3 := Human{name: "Ama", age: 20, nationality: "Ghanaian"}

	e1 := Employee{Human: h1, skills: []string{"accounting", "go"}, company: "g4m"}
	e2 := Employee{Human: h2, skills: []string{"python", "go"}, company: "google"}
	e3 := Employee{Human: h3, skills: []string{"devOps", "go"}, company: "airbnb"}

	humanbeings := []HumanBeing{e1, e2, e3}
	fmt.Println(allEmployees(humanbeings...))

}
