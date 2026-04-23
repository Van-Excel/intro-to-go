package main

import (
	"fmt"
	// "net/http"
)

func main() { // typical design problem you need to think through
	fmt.Println("test")

	fmt.Println("enter your name:")

	var name string
	fmt.Scanf("%s", &name)
	fmt.Printf("user entered:%s \n", name)

	var i int
	for i = 1; i < 11; i++ {
		if i%2 == 0 {
			fmt.Printf("%d: even\n", i)

		} else {
			fmt.Printf("%d: odd\n", i)

		}

	}

	//switch statement
	age := 2
	if age == 1 {
		fmt.Println(3, "lord")
	} else if age == 2 {
		fmt.Print(2, "dog")
	}
	fmt.Print("name")
	fmt.Print("\n")
	switch age + 2 {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	case 3:
		fmt.Println("three")
	case 4:
		fmt.Println("four")
	default:
		fmt.Println("not found")
	}

	//collections

	//declare an array
	var arrayOfInt [10]int

	// initialise it
	arrayOfInt[0] = 12
	arrayOfInt[1] = 13
	fmt.Println("array of integers:", arrayOfInt)

	//shorthand of declaring and initializing arrays
	vec := [4]int{1, 4, 5, 6}
	fmt.Println("shorthand arrays:", vec)

	//computing total of elements in vec
	total := 0
	for i := 0; i < len(vec); i++ {
		total += vec[i]

	}
	fmt.Println("total of vec:", total)

	//finding average of elements in vec
	fmt.Println("average of vec:", total/len(vec))

	// use range. we can access iterable variable with this form
	for index, value := range vec {
		fmt.Println(index, ":", value)
		fmt.Printf("%d:%d\n", index, value)
	}

	// slices, make, append, copy

	var slice1 []float64          // this has no backing array. need to use append
	slice1 = append(slice1, 1.43) // sets cap and len values
	fmt.Println("checking slice:", slice1)
	fmt.Println("cap and len of slice1", cap(slice1), len(slice1))

	// declare and initialise slice
	// used when you know elements ahead of time
	slice2 := []int{1, 3, 5, 6, 78, 69} //literal declaration. cap and len are known and set
	fmt.Println("cap of slice2 literal:", cap(slice2))
	fmt.Println("len of slice2 literal:", len(slice2))
	slice2[3] = 100 // modify within len of slice2
	fmt.Println("slice2 after mod:", slice2)

	// to add  a new element to slice2 we use append.
	// create new backing array, copy old elements into it and append new element
	slice2 = append(slice2, 900)

	// use make if we know either len or both len and future growth but not elements
	// assigns a zero value of int64 to backing array slots
	// this will allocate memory for an array that can hold 20 int64
	slice3 := make([]int64, 20)
	slice3[1] = 23
	fmt.Println("slice3:", slice3)

	// creating slices from existing array
	arr := [5]float64{1, 2, 3, 4, 5}
	x := arr[0:3] // len = 3, cap = 5
	fmt.Println("len and cap of x:", len(x), cap(x))
	x = append(x, 99, 100) // in this scenario append() over writes the last 2 elements
	fmt.Println(x)

	// new backing array is created because append() exceeded cap()
	x = append(x, 99, 100, 101) // 6 elements total, cap 5 exceeded

	fmt.Println("array for slice x:", x) // [1 2 3 99 100 101] — new backing array!
	fmt.Println("array for arr:", arr)   // [1 2 3 4 5]          — original UNCHANGED

	// x[:5] we can also reslice if we want to access last 2 elements
	// allows us to either read or modify or overwrite it
	//x[3] = 99      // modify existing arr[3] (was 4)
	//x[4] = 100     // modify existing arr[4] (was 5)
	// we can also use existing array to modify since we point to same array
	//arr[3] = 99
	//arr[4] = 100

	// we can also use a 3 value slicing index to manipulate cap
	//x := arr[0:3:3]  // len 3, cap 3

	//x = append(x, 99)  // exceeds cap → ALWAYS allocates new array
	// arr is safe, x is independent
	ya := [5]int{2, 3, 4, 5, 6}
	xa := ya[:3:3]
	fmt.Println("cap and len of xa:", cap(xa), len(xa))
	xa = append(xa, 89, 80)
	fmt.Println("new cap and len of xa:", cap(xa), len(xa))
	// cap is 6 here instead of 5 because when go allocates new backing array, it usually
	// uses 2 * cap()
	//new array -> [2 | 3 | 4 | 89 | 80 | _ ]

	/*🧠 Mental Model (your “pages” analogy)

		Think of:

		ya = a full book (5 pages)
		xa := ya[:3:3] = you photocopy ONLY first 3 pages AND forbid access to the rest

		So when you try to “write more pages”:

	👉 You must create a new book
	*/

}

// dsa, system design (HLD, LLD), comp.arch, comp network, OS,
// econs(finance, accounting), structures(linear algebra, calculus, probability and statistics)
// for work - python, sql, drf, devOps, frontend (react +nextjs), cloud
