package main

import (
	"fmt"
	"math"
)

type Circle struct {
	x      float64
	y      float64
	radius float64
}

func circleDimension(c Circle) float64 {
	return math.Sqrt(c.x * c.y)

}

func circleArea(c Circle) float64 {
	return math.Pi * c.radius * c.radius
}

func modifyCircle(c *Circle) Circle {
	c.x = 20
	c.y = 20
	return *c
}

func main() {
	fmt.Println("ready!")

	// structs

	circle := Circle{x: 2.0, y: 2.3, radius: 9.0}
	result := circleDimension(circle)
	fmt.Println("circle:", result, "metres")
	fmt.Println(circleArea(circle))

	fmt.Println("modify:", modifyCircle(&circle))

	//methods and receivers
}
