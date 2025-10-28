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

type Rectangle struct {
	length float64
	width  float64
}

func circleDimension(c Circle) float64 {
	return math.Sqrt(c.x * c.y)

}

// calculate the area of a circle
func (circle *Circle) area() float64 {
	return math.Pi * circle.radius * circle.radius
}

// calculate the area of a rectangle
func (rectangle *Rectangle) area() float64 {
	return rectangle.length * rectangle.width

}

// both circle and rectangle have area methods.
// this is where we can abstract them and use interfaces
// an interface is a collection of methods
type Shape interface {
	area() float64
}

// we can now create a method which can accept both rectangles and circles

func areaOfShapes(s Shape) float64 {
	return s.area()
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
	fmt.Println("area of circle:", circle.area())

	fmt.Println("modify:", modifyCircle(&circle))

	r := Rectangle{2.0, 3.0}
	fmt.Println("area of rectangle:", r.area())
	fmt.Println("using interfaces:", areaOfShapes(&circle))

	//methods and receivers
	// implement some interfaces and learn
}
