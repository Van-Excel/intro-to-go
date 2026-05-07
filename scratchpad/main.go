package main

import (
	"bytes"
	"fmt"
	"time"
)

// create a struct
type testStruct struct {
	name string
	age  int
}

func NewtestStruct() testStruct {
	return testStruct{
		name: "van", age: 12,
	}
}

func ping(c chan testStruct, item testStruct) chan testStruct {
	fmt.Println("about to write data")

	c <- item // send message to channel c or write something to a queue
	fmt.Println("printing what is in a channel:", c)
	// msg := <-c // this is a deadlock since there is no writer waiting
	fmt.Println("writer unblocked and done writing to channel")
	return c

}

func main() {
	fmt.Print("scratch---")
	fmt.Println()

	// var arr [4]int
	// arr[0] = 9
	// fmt.Println(arr)

	// slice1 := arr[0:2] // what is the len and cap of slice?
	// fmt.Println(slice1)
	// fmt.Println("slice1:", "cap:", cap(slice1), "len:", len(slice1))
	// slice1 = append(slice1, 21) // does it overwrite arr or different array
	// fmt.Println("slice1 after append:", slice1)
	// fmt.Println("arr after append:", arr)
	// fmt.Println("slice1:", "cap:", cap(slice1), "len:", len(slice1), "slice1:", slice1, "arr:", arr)
	// slice1 = append(slice1, 90, 23, 45, 56)
	// fmt.Println("slice1 after append2:", slice1, "cap after append2:", cap(slice1))
	// fmt.Println("arr after append2:", arr)

	// arr := [5]float64{1, 2, 3, 4, 5}
	// x := arr[0:3] // len = 3, cap = 5
	// fmt.Println("len and cap of x:", len(x), cap(x))
	// x = append(x, 99, 100) // in this scenario append() over writes the last 2 elements
	// fmt.Println(x)

	// // new backing array is created because append() exceeded cap()
	// x = append(x, 99, 100, 101) // 6 elements total, cap 5 exceeded

	// fmt.Println("array for slice x:", x) // [1 2 3 99 100 101] — new backing array!
	// fmt.Println("array for arr:", arr)   // [1 2 3 4 5]          — original UNCHANGED

	// arr := [5]float64{1, 2, 3, 4, 5}
	// slice2 := arr[0:3] // len = 3, cap = 5
	// fmt.Println("first slicing:", slice2, "cap of slice2:", cap(slice2), "len of slice2:", len(slice2))
	// slice2 = arr[:5] // reslicing
	// fmt.Println("second slicing:", slice2, "cap of slice2:", cap(slice2), "len of slice2:", len(slice2))

	// buffer package
	// var b string
	// b = "hello"
	// buffer := bytes.NewBuffer([]byte(b))
	// fmt.Println(buffer.String())
	// buffer.Write([]byte("Love World"))
	// fmt.Println("after appending:", buffer)

	// // when you are getting confused about slice just default to array and add
	// // features of slice to your thinking. Build from the ground up.
	// // always try to trace state and data structure. makes it easier to understand
	// arr2 := [3]byte{1, 2, 4}
	// bufferSlice := arr2[:2]
	// bBuffer := bytes.NewBuffer(bufferSlice)
	// fmt.Println("before appending:", bBuffer.String())
	// bBuffer.Write([]byte("Love World"))
	// fmt.Println("after appending:", bBuffer.Bytes())

	// var c string
	// c = "Vanexcel"
	// bufSlice := bytes.NewBufferString(c)
	// fmt.Println(bufSlice.String())

	// revise buffer
	// create a string
	name := "Aloto Manola"
	// convert to byte slice so you can pass it to buffer in bytes struct
	// also means string will be copied to new array which slice will point to
	// since strings are immutable so in read only memory

	revBuffer := bytes.NewBuffer([]byte(name)) //[]byte() is type casting
	fmt.Println("logging items in buffer as a string:", revBuffer.String())
	fmt.Println("logging items in buffer as bytes:", revBuffer.Bytes())

	// secret to understanding go code is understanding the internals of slices
	// which means you have to understand memory and how data looks in it
	// functions, variadic functions, unpacking collections as individual items
	// collecting individual items into collections, generics and types

	// play with channels
	col := NewtestStruct()
	newChannel := make(chan testStruct)
	// fmt.Println(ping(newChannel, *col)) // why is this wrong thinking
	go ping(newChannel, col) // why doesn't it print the fmt line in the function
	// how do you use c ?
	time.Sleep(2 * time.Millisecond)
	fmt.Println(newChannel) // why is it a pointer or address in memory? channels are reference types
	msg := <-newChannel
	fmt.Println(msg)

	time.Sleep(time.Millisecond * 5)

}
