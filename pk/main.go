package main

import (
	"fmt"
	"strings"
)

// packages- bytes, string, strconv, time, random

func main() {
	fmt.Println("Go Packages")

	//strings- contains, count, HasPrefix,

	v := "Vanexcel"
	fmt.Println(strings.Contains(v, "van"))
	fmt.Println(strings.HasPrefix(v, "va"))
	fmt.Println("returns 1 + number of code points",
		strings.Count(v, "")) // revise strings, runes and unicode code point

	s := "van"
	for _, r := range s {
		fmt.Println(r) // 118, 97, 110  ← code points
		fmt.Printf("unicode code point:%d, character:%c, utf-8: %b, type: %T \n", r, r, []byte(string(r)), r)
	}
}
