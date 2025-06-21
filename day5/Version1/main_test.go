package main

import (
	"math"
	"testing"
)

func TestCost(t *testing.T) {
	tests := []struct {
		name     string
		shape    Shape
		expected float64
	}{
		{
			name:     "Rectangle 5x4",
			shape:    &Rectangle{length: 5, breadth: 4},
			expected: 0.2 * 20, // area = 20
		},
		{
			name:     "Circle with radius 3",
			shape:    &Circle{radius: 3},
			expected: 0.5 * (math.Pi * 9),
		},
		{
			name: "Unknown shape (anonymous)",
			shape: &struct{ Shape }{
				Shape: fakeShape{area: 10},
			},
			expected: 1 * 10,
		},
	}

	for _, tt := range tests {
		got := cost(tt.shape)
		if math.Abs(got-tt.expected) > 1e-6 {
			t.Errorf("%s: cost() = %f, expected %f", tt.name, got, tt.expected)
		}
	}
}

type fakeShape struct {
	area float64
}

func (f fakeShape) Area() float64 {
	return f.area
}
