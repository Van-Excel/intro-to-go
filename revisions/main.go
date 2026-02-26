package main

import (
	"fmt"
	// "net/http"
)

type drive interface {
	drive() string
}

type Bmw struct {
	color string
	year  int
	model string
}

func (b *Bmw) drive() string {
	return b.model + "is driving"
}

type Benz struct {
	model string
	color string
}

func (b *Benz) drive() string {
	return b.model + " " + "is driving"
}

func driveCar(i drive) string {
	return i.drive()
}




func main() { // typical design problem you need to think through
	fmt.Println("test")
	b := Benz{model: "C700", color: "red"}
	fmt.Println(driveCar(&b))

}

// dsa, system design (HLD, LLD), comp.arch, comp network, OS,
// econs(finance, accounting), structures(linear algebra, calculus, probability and statistics)
// for work - python, sql, drf, devOps, frontend (react +nextjs), cloud
