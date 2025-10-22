package main

import (
	"fmt"
)

// this variable can be accessed by other functions
var phone = "233531889081"

func main() {

	var fName string = "Vanexcel"
	lName := "Acheampong"
	fmt.Println("first name is:", fName)
	fmt.Printf("The fullname is %s %s\n", fName, lName)
	fmt.Println("The fullname concatenated is:", fName+" "+lName)

	//declare variable
	var school string

	// initialize variable
	school = "alpha beta"

	fmt.Println("name of school:", school)
	school = "martin de porres"
	fmt.Println("name of new school is:", school)
	fmt.Println("is your school alpha beta?", school == "alpha beta")
	outOfScope()

	// fmt.Println("Enter a word:")
	// var scannedValue string
	// fmt.Scanf("%s", &scannedValue)
	// fmt.Println("the scanned word is:", scannedValue)

	//constants
	const dob = "23-10-1997"
	const Pi float64 = 3.147

	// variables are mappings of human readable names to memory addresses
	// you are essentially telling the compiler to allocate some space
	// and map the address of that memory block to a human readable name

	//TASKS
	// write a program that converts from fahrenheit to celcius
	// calling function
	fmt.Println("Enter the temperature you want to convert to celcius")
	var userInput float64
	fmt.Scanf("%f", &userInput)
	fmt.Println("the user input is:", userInput)
	convertFahrenheitToCelcius(userInput)
}

func outOfScope() {
	fmt.Println("phone number:", phone)
}

//TASKS
// write a program that converts from fahrenheit to celcius

func convertFahrenheitToCelcius(temperature float64) float64 {
	var convertedInput float64 = (temperature - 32) * 5 / 9
	fmt.Printf("the value of %f fahrenheit converted to celcius is %f", temperature, convertedInput)

	return convertedInput

}
