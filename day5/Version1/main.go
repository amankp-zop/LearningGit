package main

import (
	"fmt"
	"math"
)

// Shape is an interface that defines a method to calculate area.
type Shape interface {
	Area() float64
}

// Rectangle represents a rectangle with length and breadth.
type Rectangle struct {
	length  float64
	breadth float64
}

// Circle represents a circle with a given radius.
type Circle struct {
	radius float64
}

// Area calculates and returns the area of the circle.
func (cir *Circle) Area() float64 {
	return cir.radius * cir.radius * math.Pi
}

// Area calculates and returns the area of the rectangle.
func (rec *Rectangle) Area() float64 {
	return rec.breadth * rec.length
}

// findArea takes a Shape and returns its area.
func findArea(s Shape) float64 {
	return s.Area()
}

func main() {
	// Create a Rectangle instance.
	rec := Rectangle{
		length:  10,
		breadth: 10,
	}

	// Create a Circle instance.
	cir := Circle{
		radius: 20,
	}

	// Print the area of the circle.
	fmt.Println("The Area of the circle is:", findArea(&cir))
	// Print the area of the rectangle.
	fmt.Println("The Area of the rectangle is:", findArea(&rec))
}
