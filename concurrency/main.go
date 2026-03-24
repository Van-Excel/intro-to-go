package main

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

func countLetters(url string, table []int, resourceLock *sync.Mutex, waiter *sync.WaitGroup) []int {
	defer waiter.Done()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Webpage is not accessible")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	for _, b := range body {
		convertedchar := strings.ToLower(string(b))
		c := strings.Index(allLetters, convertedchar)
		if c >= 0 {
			resourceLock.Lock()
			table[c] = table[c] + 1
			resourceLock.Unlock()

		}

	}
	fmt.Println("completed:", url)

	return table

}

func main() {
	fmt.Println("Concurrency ")
	fmt.Println("number of CPUs:", runtime.NumCPU())

	// showing how goroutines can share variables using the main goroutine
	count := 5
	go countdown(&count)
	for count > 0 {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("count:", count)
	}
	var resourceLock sync.Mutex
	var waiter sync.WaitGroup

	frequency := make([]int, 26)
	for i := 1000; i < 1200; i++ {
		url := fmt.Sprintf("https://rfceditor.org/rfc/rfc%d.txt", i)
		waiter.Add(1)
		go countLetters(url, frequency, &resourceLock, &waiter)

	}
	waiter.Wait()

	// table :=
	// fmt.Println("table", table)
	// time.Sleep(10 * time.Second)
	for index, character := range allLetters {
		fmt.Printf("character: %c- count:%d\n", character, frequency[index])
	}

}

func countdown(count *int) {
	for *count > 0 {
		time.Sleep(1 * time.Second)
		*count -= 1

	}
}

// type LetterCounter struct {
//     mu      sync.Mutex
//     waiter  sync.WaitGroup
//     table   []int
// }

// func (lc *LetterCounter) count(url string) {
//     defer lc.waiter.Done()

//     resp, err := http.Get(url)
//     if err != nil {
//         return
//     }
//     defer resp.Body.Close()

//     body, _ := io.ReadAll(resp.Body)
//     for _, b := range body {
//         convertedchar := strings.ToLower(string(b))
//         c := strings.Index(allLetters, convertedchar)
//         if c >= 0 {
//             lc.mu.Lock()
//             lc.table[c]++
//             lc.mu.Unlock()
//         }
//     }
// }

// // in main
// lc := &LetterCounter{
//     table: make([]int, 26),
// }

// for i := 1000; i < 1200; i++ {
//     url := fmt.Sprintf("https://rfceditor.org/rfc/rfc%d.txt", i)
//     lc.waiter.Add(1)
//     go lc.count(url)
// }

// lc.waiter.Wait()
