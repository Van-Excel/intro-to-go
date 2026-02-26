package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("http")

	http.ListenAndServe(":8000", nil)

}
