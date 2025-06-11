package main

import (
	"fmt"
)

func main() {
	var (
		name   string
		age    int
		salary int
	)

	fmt.Println("Enter your name : ")
	fmt.Scanln(&name)
	fmt.Println("Enter your age : ")
	fmt.Scanln(&age)
	fmt.Println("Enter your salary : ")
	fmt.Scanln(&salary)

	fmt.Printf("Hello, %v. I am glad to know you are %v years old. Your salary will be %v rupees", name, age, salary)

}
