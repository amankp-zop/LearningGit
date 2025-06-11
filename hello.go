package main

import (
	"fmt"
)

var x = 10

const (
	VALUE = 2001
	PI    = 3.14
	AGE   = 100
)

func add(x, y int) int {
	return x + y
}

func swap(x, y int) (int, int) {
	return y, x
}

func main() {

	fmt.Println(VALUE)
	fmt.Println(PI)
	fmt.Println(AGE)
	fmt.Println("HEY there")

	fmt.Printf("Value has value: %v and type: %T", VALUE, VALUE)
	var arr = [...]int{1, 2, 3}
	fmt.Println(arr)
	arr1 := [...]string{"Aman", "Kumar", "Pandey"}
	fmt.Println(arr1)
	fmt.Println(arr1[0])

	slice1 := []int{1, 2, 3}
	fmt.Println(slice1)

	slice2 := arr1[0:2]
	fmt.Println(slice2)

	fmt.Println(len(slice1))
	fmt.Println(cap(slice1))

	slice1 = append(slice1, 50)
	fmt.Println(slice1)
	fmt.Println(add(10, 20))

	fmt.Println(swap(40, 50))

	var name string

	fmt.Scanln(&name)
	fmt.Println(name)

}
