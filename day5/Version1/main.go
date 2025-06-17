package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
}

type Rectangle struct {
	length  float64
	breadth float64
}

type Circle struct {
	radius float64
}

func (cir *Circle) Area() float64 {
	return cir.radius * cir.radius * math.Pi
}

func (rec *Rectangle) Area() float64 {
	return rec.breadth * rec.length
}

func findArea(s Shape) float64 {
	return s.Area()
}

func main() {

	rec := Rectangle{
		length:  10,
		breadth: 10,
	}

	cir := Circle{
		radius: 20,
	}

	fmt.Println("The Area of the circle is:", findArea(&cir))
	fmt.Println("The Area of the rectangle is:", findArea(&rec))

}
